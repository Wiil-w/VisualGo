package noise

import (
	"image"
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Noise struct {
	Map    [][]float64
	Seed   int
	Size   int
	Type   string
	Screen *image.RGBA
	// Tick   int
}

func (noise *Noise) Generate() {
	height := 255.0
	precision := 0.01

	GeneratePermutation()

	// generate noise
	switch noise.Type {
	case "perlin":
		for y := range noise.Map {
			for x := range noise.Map {
				n := Perlin(float64(x)*precision, float64(y)*precision)

				n = (n + 1.0) / 2.0 // transfor range from [-1, 1] to [0, 1]

				noise.Map[y][x] = math.Round(height * n)
			}
		}
	}
}

func (noise *Noise) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || noise.Screen == nil {
		noise.Generate()
		noise.Screen = image.NewRGBA(image.Rect(0, 0, noise.Size, noise.Size))

		for y := 0; y < noise.Size; y++ {
			for x := 0; x < noise.Size; x++ {
				i := (y*noise.Size + x) * 4
				noise.Screen.Pix[i] = uint8(noise.Map[y][x])
				noise.Screen.Pix[i+1] = uint8(noise.Map[y][x])
				noise.Screen.Pix[i+2] = uint8(noise.Map[y][x])
				noise.Screen.Pix[i+3] = 0xff
			}
		}
		spew.Dump(noise.Screen.Pix[:96])

		// f, _ := os.Create("Modules/Noise/imageTest8.png")
		// png.Encode(f, noise.Screen)
	}

	// noise.Tick++
	// if noise.Tick == 60 {
	// 	noise.Tick = 0
	// }
	return nil
}

func (noise *Noise) Draw(bg *ebiten.Image) error {
	if noise.Screen != nil {
		// bg.Fill(color.White)
		bg.ReplacePixels(noise.Screen.Pix)
	}
	return nil
}

func Start(size int, noiseType string) *Noise {

	noise := Noise{Type: noiseType, Size: size}

	// create map
	switch noise.Type {
	case "perlin":
		grid := make([][]float64, noise.Size)
		for i := range grid {
			grid[i] = make([]float64, noise.Size)
		}
		noise.Map = grid
	}
	return &noise
}
