package cave

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	base "killtime/base"
	noise "killtime/projects/Noise"
)

const (
	// number o iterations
	maxGen = 0
)

type Cave struct {
	Mesh   base.Mesh
	Grid   [][]Square
	Map    [][]bool
	Height int
	Width  int
	Gen    int
	Count  int
}

func (cave *Cave) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cave.new()
	}

	// Simple Sqaure
	// cave.Mesh.AddVertex(0,0)
	// cave.Mesh.AddVertex(0,1)
	// cave.Mesh.AddVertex(1,0)
	// cave.Mesh.AddVertex(1,1)
	// cave.Mesh.AddTriangle(0, 1, 2)
	// cave.Mesh.AddTriangle(1, 2, 3)

	return nil
}

func (cave *Cave) Draw(bg *ebiten.Image) error {
	cave.Mesh.Draw(bg)

	return nil
}

func Start(screenWidth, screenHeight int, args ...interface{}) *Cave {

	if 2 > len(args) {
		panic("Cave gen needs 2 args to start (height and width)")
	}
	h, ok := args[0].(int)
	if !ok {
		panic("Cave height needs to be int")
	}
	w, ok := args[1].(int)
	if !ok {
		panic("Cave width needs to be int")
	}

	cave := &Cave{
		Mesh:   base.NewMesh(screenWidth, screenHeight),
		Height: h,
		Width:  w,
	}

	cave.new()

	return cave
}

func (cave *Cave) new() {

	cave.Gen = 0
	cave.Count = 0

	// create/generate map
	cave.newMap()

	// update map epoch
	for cave.Gen < maxGen {
		cave.updateMap()
		cave.Gen++
	}

	// create/populate grid
	cave.newGrid()

	// populate vertices and triangles arrays
	cave.createMesh()
	fmt.Println(cave.Count)
}

func (cave *Cave) newMap() {
	fmt.Println("new map")
	precision := 0.05

	noise.GeneratePermutation()
	cave.Map = make([][]bool, cave.Height)
	for i := 0; i < cave.Height; i++ {

		cave.Map[i] = make([]bool, cave.Width)
		for k := 0; k < cave.Width; k++ {

			cave.Map[int(i)][int(k)] = noise.Perlin(float64(k)*precision, float64(i)*precision) > 0.0
		}
	}
}

func (cave *Cave) updateMap() {
	fmt.Println("update map")

	duplicate := make([][]bool, cave.Height)
	for i, row := range cave.Map {
		duplicate[i] = make([]bool, cave.Width)
		for k := range row {
			duplicate[i][k] = cave.countNeighbors(k, i)
		}
	}

	cave.Map = duplicate
}

func (cave *Cave) countNeighbors(x, y int) bool {

	var count int
	for i := y - 1; i <= y+1; i++ {
		for k := x - 1; k <= x+1; k++ {
			if i == y && k == x {
				continue
			}
			if i < 0 || k < 0 || i == cave.Height || k == cave.Width || cave.Map[i][k] {
				count++
				if count > 4 {
					return true
				}
			}
		}
	}
	return false
}

func (cave *Cave) newGrid() {
	fmt.Println("new grid")

	var nodeMap = make([][]MainNode, cave.Height)
	for i, row := range cave.Map {
		nodeMap[i] = make([]MainNode, cave.Width)
		for k, value := range row {
			position := base.Vector2{X: float64(k), Y: float64(i)}
			nodeMap[i][k] = MainNode{
				Node{position, -1},
				value,
				Node{position.Right(.5), -1},
				Node{position.Top(.5), -1},
			}
		}
	}

	cave.Grid = make([][]Square, cave.Height-1)
	for i := 0; i < cave.Height-1; i++ {
		cave.Grid[i] = make([]Square, cave.Width-1)
		for k := 0; k < cave.Width-1; k++ {

			cave.Grid[i][k] = Square{
				topLeft:      nodeMap[i][k],
				topRight:     nodeMap[i][k+1],
				bottomLeft:   nodeMap[i+1][k],
				bottomRight:  nodeMap[i+1][k+1],
				centerTop:    nodeMap[i][k].Right,
				centerLeft:   nodeMap[i][k].Top,
				centerRight:  nodeMap[i][k+1].Top,
				centerBottom: nodeMap[i+1][k].Right,
			}
			if cave.Grid[i][k].topLeft.Active {
				cave.Grid[i][k].value += 1
			}
			if cave.Grid[i][k].topRight.Active {
				cave.Grid[i][k].value += 2
			}
			if cave.Grid[i][k].bottomRight.Active {
				cave.Grid[i][k].value += 4
			}
			if cave.Grid[i][k].bottomLeft.Active {
				cave.Grid[i][k].value += 8
			}
		}
	}
}

func (cave *Cave) createMesh() {

	cave.Mesh.Clear()
	fmt.Println("create mesh")
	for y, row := range cave.Grid {
		for x := range row {
			if cave.Grid[y][x].value != 0 {
				cave.triangulate(&cave.Grid[y][x])
			}
		}
	}
}

func (cave *Cave) triangulate(square *Square) {

	switch square.value {
	// One corner - One Triangle
	case 1:
		cave.meshFromPoints(square.topLeft.Node, square.centerTop, square.centerLeft)
	case 2:
		cave.meshFromPoints(square.topRight.Node, square.centerTop, square.centerRight)
	case 4:
		cave.meshFromPoints(square.bottomRight.Node, square.centerBottom, square.centerRight)
	case 8:
		cave.meshFromPoints(square.bottomLeft.Node, square.centerBottom, square.centerLeft)

	// Two adjacent corners - Two Triangles
	case 3:
		cave.meshFromPoints(square.topLeft.Node, square.centerLeft, square.topRight.Node, square.centerRight)
	case 6:
		cave.meshFromPoints(square.topRight.Node, square.centerTop, square.bottomRight.Node, square.centerBottom)
	case 9:
		cave.meshFromPoints(square.bottomLeft.Node, square.centerBottom, square.topLeft.Node, square.centerTop)
	case 12:
		cave.meshFromPoints(square.bottomRight.Node, square.centerRight, square.bottomLeft.Node, square.centerLeft)

	// Three corners - Three Triangles
	case 7:
		cave.meshFromPoints(square.topLeft.Node, square.centerLeft, square.topRight.Node, square.centerBottom, square.bottomRight.Node)
	case 11:
		cave.meshFromPoints(square.bottomLeft.Node, square.centerBottom, square.topLeft.Node, square.centerRight, square.topRight.Node)
	case 13:
		cave.meshFromPoints(square.topLeft.Node, square.centerTop, square.bottomLeft.Node, square.centerRight, square.bottomRight.Node)
	case 14:
		cave.meshFromPoints(square.bottomLeft.Node, square.centerLeft, square.bottomRight.Node, square.centerTop, square.topRight.Node)

	// Two oposite corners - Four Triangles
	case 5:
		cave.meshFromPoints(square.centerTop, square.centerRight, square.topLeft.Node, square.bottomRight.Node, square.centerLeft, square.centerBottom)
	case 10:
		cave.meshFromPoints(square.centerTop, square.centerLeft, square.topRight.Node, square.bottomLeft.Node, square.centerRight, square.centerBottom)

	// Four corners - Two Triangles
	case 15:
		cave.meshFromPoints(square.bottomLeft.Node, square.topLeft.Node, square.bottomRight.Node, square.topRight.Node)
	}
}

func (cave *Cave) meshFromPoints(points ...Node) {
	cave.addVertices(points)

	for i := 0; i < len(points)-2; i++ {
		cave.Mesh.AddTriangle(points[i].Index, points[i+1].Index, points[i+2].Index)
	}

}

func (cave *Cave) addVertices(points []Node) {
	for i := range points {
		if points[i].Index == -1 {
			cave.Count++
			points[i].Index = len(cave.Mesh.Vertices)
			cave.Mesh.AddVertex(points[i].Position)
		}
	}
}

type Square struct {
	topLeft, topRight, bottomLeft, bottomRight       MainNode
	centerTop, centerRight, centerBottom, centerLeft Node
	value                                            int
}

type MainNode struct {
	Node
	Active bool
	Right  Node
	Top    Node
}

type Node struct {
	Position base.Vector2
	Index    int
}
