package coneways

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Menu struct {
	Screen *ebiten.Image
}

func (menu *Menu) Draw(bg *ebiten.Image) error {
	if menu.Screen == nil {
		menu.Screen.Fill(color.Black)

		opts := &ebiten.DrawImageOptions{}
		menu.Screen.DrawImage(menu.Screen, opts)
	}
	opts := &ebiten.DrawImageOptions{}
	bg.DrawImage(menu.Screen, opts)
	return nil
}

func (menu *Menu) Update() error {
	return nil
}
