package main

import (
	"image/color"
	cave "killtime/projects/CaveGen"
	conways "killtime/projects/Conways"
	noise "killtime/projects/Noise"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Menu struct {
	Modulos     Modulos
	Running     string
	ModuleNames []string
}

type Modulos interface {
	Draw(*ebiten.Image) error
	Update() error
}

const (
	windowSizeScale = 1.5
	screenWidth     = int(480 * windowSizeScale) // 480, 640
	screenHeight    = int(480 * windowSizeScale)
)

func (menu *Menu) Draw(bg *ebiten.Image) {

	// draw bg color
	bg.Fill(color.RGBA{0, 0, 0, 0xff})

	if menu.Modulos == nil {
		// draw menu
	} else {
		menu.Modulos.Draw(bg)
	}
}

func (menu *Menu) Update() error {

	if menu.Modulos == nil {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// Calcula em qual botao apertou
			menu.Running = "cave"
			switch menu.Running {
			case "conways":
				menu.Modulos = conways.Start()
			case "perlin":
				menu.Modulos = noise.Start(screenWidth, "perlin")
			case "cave":
				menu.Modulos = cave.Start(screenWidth, screenHeight, 50, 50)
			}
		}
	} else {
		menu.Modulos.Update()
	}

	return nil
}

func (menu *Menu) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func MainMenu() *Menu {
	var mod Modulos
	return &Menu{
		Running: "menu",
		Modulos: mod,
		ModuleNames: []string{
			"coneways",
			"perlin",
			"cave",
		},
	}
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Will's Corner")

	if err := ebiten.RunGame(MainMenu()); err != nil {
		log.Fatal(err)
	}
}
