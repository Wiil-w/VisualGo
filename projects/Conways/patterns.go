package coneways

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Patterns struct {
	Screen        *ebiten.Image
	Template      map[string][][][]int
	PatternsTypes []string
	Visible       [][][]int
	Pause         bool
	Type          int
	Page          int
	xScale        float64
	yScale        float64
	xStart        int
	yStart        int
	// position [9][2]int
}

func (patterns *Patterns) Draw(bg *ebiten.Image) error {

	if patterns.Visible == nil {
		patterns.getVisibleTemplates()
		// patterns.Pause = true
	}
	patterns.Screen.Fill(color.Black)

	for i, pattern := range patterns.Visible {

		position := i % 9                                       // which of the 9 patterns it is on the screen
		y := patterns.yStart*(position/3+1) - len(pattern)/2    // middle of the pattern in the screen
		x := patterns.xStart*(position%3+1) - len(pattern[0])/2 // middle of the pattern in the screen

		for j, row := range pattern {
			for k := range row {
				if pattern[j][k] == 1 {
					ebitenutil.DrawRect(patterns.Screen, float64((x+k))*patterns.xScale, float64((y+j))*patterns.yScale, patterns.xScale, patterns.yScale, color.White)
				}
			}
		}

	}

	opts := &ebiten.DrawImageOptions{}
	bg.DrawImage(patterns.Screen, opts)
	return nil
}

func (patterns *Patterns) Update() error {
	// patterns.getVisibleTemplates()

	for p, pattern := range patterns.Visible {

		// duplicate the board
		copyPattern := make([][]int, len(pattern))
		for i := range pattern {
			copyPattern[i] = make([]int, len(pattern[i]))
			copy(copyPattern[i], pattern[i])
		}

		// Iterate over the cells
		var liveNeighbors int
		for y := 0; y < len(pattern); y++ {
			for x := 0; x < len(pattern[y]); x++ {

				liveNeighbors = patterns.checkNeighbors(x, y, pattern)

				switch pattern[y][x] {
				case 1:
					if !(liveNeighbors > 1 && liveNeighbors < 4) {
						copyPattern[y][x] = 0
					}
				case 0:
					if liveNeighbors == 3 {
						copyPattern[y][x] = 1
					}
				}
			}
		}
		copy(patterns.Visible[p], copyPattern)
	}
	return nil
}

func (patterns *Patterns) checkNeighbors(x, y int, pattern [][]int) int {

	// Iterate over the cell's neightbors
	var liveNeighbors, x1, y1 int
	for i := -1; i < 2; i++ {
		y1 = y + i
		switch {
		case y1 == -1:
			y1 = len(pattern) - 1
		case y1 == len(pattern):
			y1 = 0
		default:
			if y1 < -1 || y1 > len(pattern) {
				continue
			}
		}
		for k := -1; k < 2; k++ {
			x1 = x + k
			switch {
			case x1 == -1:
				x1 = len(pattern[0]) - 1
			case x1 == len(pattern[0]):
				x1 = 0
			default:
				if x1 < -1 || x1 > len(pattern[0]) || (i == k && i == 0) {
					continue
				}
			}
			if pattern[y1][x1] == 1 {
				liveNeighbors++
				if liveNeighbors > 3 {
					return liveNeighbors
				}
			}
		}
	}
	return liveNeighbors
}

func (patterns *Patterns) SelectPattern(_mouseX, _mouseY int) [][]int {

	mouseX := float64(_mouseX)
	mouseY := float64(_mouseY)
	pixelDist := 2 * 2 // distance of 2 for now

	// Click on a Pattern
	for i := 0; i < 9; i++ {
		position := i % 9                          // which of the 9 patterns it is on the screen
		patY := patterns.yStart * (position/3 + 1) // middle of the pattern in the screen
		patX := patterns.xStart * (position%3 + 1) // middle of the pattern in the screen

		if mouseX > patterns.xScale*float64(patX-pixelDist) && mouseX < patterns.xScale*float64(patX+pixelDist) && mouseY > patterns.yScale*float64(patY-pixelDist) && mouseY < patterns.xScale*float64(patY+pixelDist) {
			return patterns.Template[patterns.PatternsTypes[patterns.Type]][(patterns.Page*9)+i]
		}
	}

	// Click to change pages

	return nil
}

func (patterns *Patterns) getVisibleTemplates() error {

	patternType := patterns.PatternsTypes[patterns.Type]
	min := patterns.Page * 9
	max := (patterns.Page + 1) * 9
	if max > len(patterns.Template[patternType]) {
		max = len(patterns.Template[patternType])
	}
	if patterns.Visible == nil {
		patterns.Visible = make([][][]int, len(patterns.Template[patternType]))
	}
	for p, pattern := range patterns.Template[patternType][min:max] {
		// duplicate the visible patterns
		patterns.Visible[p] = make([][]int, len(pattern))
		for i := range pattern {
			patterns.Visible[p][i] = make([]int, len(pattern[i]))
			copy(patterns.Visible[p][i], pattern[i])
		}
	}
	// spew.Dump(patterns.Visible)

	return nil
}
