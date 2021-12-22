package base

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Mesh struct {
	Vertices         []Vector2
	Triangles        []Vector2
	width, height    float64
	xBorder, yBorder float64
	maxX, maxY       float64
	cellX, cellY     float64
}

type Vertex struct {
	X, Y float64
}

func NewMesh(screenWidth, screenHeight int) Mesh {
	mesh := Mesh{
		width:  float64(screenWidth),
		height: float64(screenHeight),
	}
	mesh.xBorder = mesh.height / 10.0
	mesh.yBorder = mesh.width / 10.0

	return mesh
}

func (mesh *Mesh) AddVertex(vector Vector2) {
	if vector.X > mesh.maxX {
		mesh.maxX = math.Ceil(vector.X)
	}
	if vector.Y > mesh.maxY {
		mesh.maxY = math.Ceil(vector.Y)
	}
	mesh.Vertices = append(mesh.Vertices, vector)
}

func (mesh *Mesh) AddTriangle(p1, p2, p3 int) {
	mesh.Triangles = append(mesh.Triangles, mesh.Vertices[p1])
	mesh.Triangles = append(mesh.Triangles, mesh.Vertices[p2])
	mesh.Triangles = append(mesh.Triangles, mesh.Vertices[p3])
}

func (mesh *Mesh) ClearVertices()  { mesh.Vertices = nil }
func (mesh *Mesh) ClearTriangles() { mesh.Triangles = nil }
func (mesh *Mesh) Clear() {
	mesh.ClearVertices()
	mesh.ClearTriangles()
}

// func (mesh *Mesh) Update() {}

func (mesh *Mesh) Draw(bg *ebiten.Image) error {

	// draw bg color
	bg.Fill(color.RGBA{0, 0, 0, 0xff})

	mesh.cellX = (mesh.height - (mesh.xBorder * 2)) / float64(mesh.maxX)
	mesh.cellY = (mesh.width - (mesh.yBorder * 2)) / float64(mesh.maxY)

	for i := 0; i < len(mesh.Triangles); i += 3 {
		p1 := mesh.Triangles[i]
		p2 := mesh.Triangles[i+1]
		p3 := mesh.Triangles[i+2]
		mesh.DrawTriangle(bg, p1, p2, p3)
	}
	return nil
}

func (mesh *Mesh) DrawTriangle(bg *ebiten.Image, p1, p2, p3 Vector2) {
	x1 := p1.X*mesh.cellX + mesh.xBorder
	y1 := p1.Y*mesh.cellY + mesh.yBorder
	x2 := p2.X*mesh.cellX + mesh.xBorder
	y2 := p2.Y*mesh.cellY + mesh.yBorder
	x3 := p3.X*mesh.cellX + mesh.xBorder
	y3 := p3.Y*mesh.cellY + mesh.yBorder

	ebitenutil.DrawLine(bg, x1, y1, x2, y2, color.White)
	ebitenutil.DrawLine(bg, x1, y1, x3, y3, color.White)
	ebitenutil.DrawLine(bg, x3, y3, x2, y2, color.White)
}
