package coneways

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Board struct {
	Screen [][]int
	// ScreenCopy [][]int
	Pattern [][]int
	X, Y    int
	Gen     int
}

func (board *Board) Update() error {

	// duplicate the board
	copyScreen := make([][]int, board.Y)
	for i := range board.Screen {
		copyScreen[i] = make([]int, board.X)
		copy(copyScreen[i], board.Screen[i])
	}

	// Iterate over the cells
	var liveNeighbors int
	for y := 0; y < board.Y; y++ {
		for x := 0; x < board.X; x++ {

			liveNeighbors = board.checkNeighbors(x, y)

			switch board.Screen[y][x] {
			case 1:
				if !(liveNeighbors > 1 && liveNeighbors < 4) {
					copyScreen[y][x] = 0
				}
			case 0:
				if liveNeighbors == 3 {
					copyScreen[y][x] = 1
				}
			}
		}
	}
	copy(board.Screen, copyScreen)
	return nil
}

func (board *Board) checkNeighbors(x, y int) int {

	// Iterate over the cell's neightbors
	var liveNeighbors, x1, y1 int
	for i := -1; i < 2; i++ {
		y1 = y + i
		switch {
		case y1 == -1:
			y1 = board.Y - 1
		case y1 == board.Y:
			y1 = 0
		default:
			if y1 < -1 || y1 > board.Y {
				continue
			}
		}
		for k := -1; k < 2; k++ {
			x1 = x + k
			switch {
			case x1 == -1:
				x1 = board.X - 1
			case x1 == board.X:
				x1 = 0
			default:
				if x1 < -1 || x1 > board.X || (i == k && i == 0) {
					continue
				}
			}
			if board.Screen[y1][x1] == 1 {
				liveNeighbors++
				if liveNeighbors > 3 {
					return liveNeighbors
				}
			}
		}
	}
	return liveNeighbors
}

func (board *Board) Draw(bg *ebiten.Image) error {
	for y := 0; y < board.Y; y++ {
		for x := 0; x < board.X; x++ {
			if board.Screen[y][x] == 1 {
				ebitenutil.DrawRect(bg, float64(x*xScale), float64(y*yScale), xScale, yScale, color.White)
			}
		}
	}
	return nil
}

func (board *Board) PlacePattern(xMouse, yMouse int) {
	// Verify if x and y is in screen
	xVals := []int{0, xMouse, screenWidth - 1}
	sort.Ints(xVals)
	xMouse = xVals[1] / xScale

	yVals := []int{0, yMouse, screenHeight - 1}
	sort.Ints(yVals)
	yMouse = yVals[1] / yScale

	// TODO: Place board.Pattern at mouse coords

	// Cerate a + of live cells
	board.Screen[yMouse][xMouse] = 1
	if xMouse+1 < board.X {
		board.Screen[yMouse][xMouse+1] = 1
	}
	if xMouse-1 >= 0 {
		board.Screen[yMouse][xMouse-1] = 1
	}
	if yMouse+1 < board.Y {
		board.Screen[yMouse+1][xMouse] = 1
	}
	if yMouse-1 >= 0 {
		board.Screen[yMouse-1][xMouse] = 1
	}

}

func (board *Board) PlacePixel(xMouse, yMouse int) {
	// Verify if x and y is on the window
	xVals := []int{0, xMouse, screenWidth - 1}
	sort.Ints(xVals)
	xMouse = xVals[1] / xScale

	yVals := []int{0, yMouse, screenHeight - 1}
	sort.Ints(yVals)
	yMouse = yVals[1] / yScale

	board.Screen[yMouse][xMouse] = 1
}
