package perlin

import (
	"math"
	"math/rand"
)

const (
	size = 0x100
	mask = size - 1
)

type Noise2D struct {
	gradients [size][2]float64
	indexes   [size]int
}

func NewNoise2D(seed int64) *Noise2D {
	noise := &Noise2D{}
	src := rand.NewSource(seed)
	rng := rand.New(src)
	copy(noise.indexes[:], rng.Perm(len(noise.indexes)))

	for i := range noise.indexes {
		noise.gradients[i][0], noise.gradients[i][1] = rand2d(rng)
	}
	return noise
}

func (n *Noise2D) At(x, y float64) float64 {
	x0, y0 := math.Floor(x), math.Floor(y)
	dx, dy := x-x0, y-y0

	s := n.weight(x0, y0, x, y)
	t := n.weight(x0+1, y0, x, y)
	u := n.weight(x0, y0+1, x, y)
	v := n.weight(x0+1, y0+1, x, y)

	sx := 3*dx*dx - 2*dx*dx*dx
	sy := 3*dy*dy - 2*dy*dy*dy
	a := s + sx*(t-s)
	b := u + sx*(v-u)

	return a + sy*(b-a)
}

func (n *Noise2D) weight(x0, y0, x, y float64) float64 {
	index := int(x0)&mask + n.indexes[int(y0)&mask]
	gradient := n.gradients[index&mask]
	return dot(gradient[0], x-x0, gradient[1], y-y0)
}

func dot(x0, y0, x1, y1 float64) float64 {
	return x0*x1 + y0*y1
}

func rand2d(rng *rand.Rand) (float64, float64) {
	return norm(2*rng.Float64()-1, 2*rng.Float64()-1)
}

func norm(x, y float64) (float64, float64) {
	length := math.Sqrt(x*x + y*y)
	return x / length, y / length
}
