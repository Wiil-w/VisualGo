package noise

import (
	"math"
	"math/rand"
	"time"
)

type Vector2 struct {
	x float64
	y float64
}

type Permutation struct {
	Perm []int
}

var (
	p Permutation = Permutation{}
)

func GeneratePermutation() {

	//Create an array (our permutation table) with the values 0 to 255 in order
	p.Perm = make([]int, 2*256)
	for i := 0; i < 256; i++ {
		p.Perm[i] = i
	}
	// shuffle it
	p.Perm = shuffle(p.Perm)
	p.Perm = append(p.Perm, p.Perm...)
}

func (vector *Vector2) Dot(otherVec *Vector2) float64 {
	return vector.x*otherVec.x + vector.y*otherVec.y
}

//return value from between -1 and 1
func Perlin(x, y float64) float64 {

	X := int(x) & 255
	Y := int(y) & 255

	xf := x - math.Floor(x)
	yf := y - math.Floor(y)

	topRight := Vector2{xf - 1.0, yf - 1.0}
	topLeft := Vector2{xf, yf - 1.0}
	bottomRight := Vector2{xf - 1.0, yf}
	bottomLeft := Vector2{xf, yf}

	//Select a value in the array for each of the 4 corners
	valueTopRight := p.Perm[p.Perm[X+1]+Y+1]
	valueTopLeft := p.Perm[p.Perm[X]+Y+1]
	valueBottomRight := p.Perm[p.Perm[X+1]+Y]
	valueBottomLeft := p.Perm[p.Perm[X]+Y]

	dotTopRight := topRight.Dot(getConstantVector(valueTopRight))
	dotTopLeft := topLeft.Dot(getConstantVector(valueTopLeft))
	dotBottomRight := bottomRight.Dot(getConstantVector(valueBottomRight))
	dotBottomLeft := bottomLeft.Dot(getConstantVector(valueBottomLeft))

	u := fade(xf)
	v := fade(yf)
	return interpolation(u,
		interpolation(v, dotBottomLeft, dotTopLeft),
		interpolation(v, dotBottomRight, dotTopRight))
}

func shuffle(arr []int) []int {
	rand.Seed(time.Now().UnixNano())
	// Fisher–Yates shuffle
	for e := 255; e > 0; e-- {
		index := rand.Intn(e + 1)
		arr[e], arr[index] = arr[index], arr[e]
	}
	return arr
}

func getConstantVector(value int) *Vector2 {
	//v is the value from the permutation table

	switch value & 3 {
	case 0:
		return &Vector2{1.0, 1.0}
	case 1:
		return &Vector2{-1.0, 1.0}
	case 2:
		return &Vector2{-1.0, -1.0}
	default:
		return &Vector2{1.0, -1.0}
	}
}

// ease curve
func fade(t float64) float64 {
	return ((6*t-15)*t + 10) * math.Pow(t, 3) // 6*t*t*t*t*t - 15*t*t*t*t + 10*t*t*t
}

func interpolation(t, a1, a2 float64) float64 {
	return a1 + t*(a2-a1)
}
