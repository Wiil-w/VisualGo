package coneways

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Board    Board
	Config   Config
	Menu     Menu
	Patterns Patterns
	Screen   string
	Pause    bool
}

type Config struct {
	Tick         int
	TickDivision int
}

const (
	windowSizeScale = 1.5
	x               = 120 * windowSizeScale
	y               = 120 * windowSizeScale
	screenWidth     = 480 * windowSizeScale // 480, 640
	screenHeight    = 480 * windowSizeScale
	xScale          = screenHeight / x
	yScale          = screenWidth / y
)

func (game *Game) Draw(bg *ebiten.Image) error {

	// draw bg color
	bg.Fill(color.RGBA{0, 0, 0, 0xff})

	switch game.Screen {

	case "menu":
		game.Menu.Draw(bg)

	case "patterns":
		game.Patterns.Draw(bg)

	case "game":
		game.Board.Draw(bg)
	}
	return nil
}

func (game *Game) Update() error {

	switch game.Screen {
	case "menu":
		game.Menu.Update()

	case "patterns":
		if !(game.Patterns.Pause) {
			// board update speed
			if game.Config.Tick%game.Config.TickDivision != 0 {
				game.Config.Tick++
			} else {
				// update board
				game.Config.Tick = 1
				game.Patterns.Update()
			}
		}

		// Left Mouse Button to select a pattern
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			pattern := game.Patterns.SelectPattern(ebiten.CursorPosition())
			if pattern != nil {
				game.Screen = "game"
				game.Board.Pattern = pattern
			}
		}

		// Space Key to pause the pattern update
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			game.Patterns.Pause = !game.Patterns.Pause
		}

	case "game":
		if !(game.Pause) {
			// board update speed
			if game.Config.Tick%game.Config.TickDivision != 0 {
				game.Config.Tick++
			} else {
				// update board
				game.Config.Tick = 1
				game.Board.Update()
			}
		}

		// Left Mouse Button to place a pattern on the board
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			game.Board.PlacePattern(ebiten.CursorPosition())
		}

		// Right Mouse Button to place a pixel on the board
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			game.Board.PlacePixel(ebiten.CursorPosition())
		}

		// Space Key to pause the board update
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			game.Pause = !game.Pause
		}
	}

	// Esc Key to open menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if game.Screen == "menu" {
			game.Screen = "game"
		} else {
			game.Screen = "menu"
		}
	}

	// Q Key to open patterns
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if game.Screen == "patterns" {
			game.Screen = "game"
		} else {
			game.Screen = "patterns"
		}
	}

	// S Key to slow down cell update
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		game.Config.TickDivision = 6
	}
	// F Key to speed up cell update
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		game.Config.TickDivision = 2
	}

	return nil
}

// func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

func Start() *Game {

	// ebiten.SetWindowSize(screenWidth, screenHeight)
	rand.Seed(time.Now().Unix())
	boardScreen := make([][]int, y)

	for k := 0; k < y; k++ {
		boardScreen[k] = make([]int, x)

		for i := 0; i < x; i++ {
			if rand.Intn(10) == 1 {
				boardScreen[k][i] = 1
			}
		}
	}

	board := Board{
		Screen: boardScreen,
		Gen:    1,
		X:      x,
		Y:      y,
	}

	patterns := Patterns{
		Screen: ebiten.NewImage(screenWidth, screenHeight),
		Template: map[string][][][]int{
			"static": { // still life
				{{0, 0, 0, 0}, {0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
				{{0, 0, 0, 0, 0, 0}, {0, 0, 1, 1, 0, 0}, {0, 1, 0, 0, 1, 0}, {0, 0, 1, 1, 0, 0}, {0, 0, 0, 0, 0, 0}},
				{{0, 0, 0, 0, 0, 0}, {0, 0, 1, 1, 0, 0}, {0, 1, 0, 0, 1, 0}, {0, 0, 1, 0, 1, 0}, {0, 0, 0, 1, 0, 0}, {0, 0, 0, 0, 0, 0}},
				{{0, 0, 0, 0, 0}, {0, 1, 1, 0, 0}, {0, 1, 0, 1, 0}, {0, 0, 1, 0, 0}, {0, 0, 0, 0, 0}},
				{{0, 0, 0, 0, 0}, {0, 0, 1, 0, 0}, {0, 1, 0, 1, 0}, {0, 0, 1, 0, 0}, {0, 0, 0, 0, 0}},
			},
			"oscillators": {
				{{0, 0, 0, 0, 0}, {0, 0, 0, 0, 0}, {0, 1, 1, 1, 0}, {0, 0, 0, 0, 0}, {0, 0, 0, 0, 0}},
				{{0, 0, 0, 0, 0}, {0, 0, 0, 0, 0}, {0, 0, 1, 1, 1}, {0, 1, 1, 1, 0}, {0, 0, 0, 0, 0}},
				{{0, 0, 0, 0, 0}, {0, 1, 1, 0, 0}, {0, 1, 1, 0, 0}, {0, 0, 0, 1, 1}, {0, 0, 0, 1, 1}},
			},
			"spaceships": {
				{{0, 0, 0, 0, 0}, {0, 0, 0, 1, 0}, {0, 1, 0, 1, 0}, {0, 0, 1, 1, 0}, {0, 0, 0, 0, 0}},
			},
		},
		PatternsTypes: []string{"static", "oscillators", "spaceships"},
		xScale:        xScale * 2,
		yScale:        yScale * 2,
		xStart:        int(x) / 8,
		yStart:        int(y) / 8,
	}

	menu := Menu{
		Screen: ebiten.NewImage(screenWidth, screenHeight),
	}

	return &Game{
		Board:    board,
		Config:   Config{TickDivision: 2},
		Menu:     menu,
		Patterns: patterns,
		Screen:   "game",
	}
}
