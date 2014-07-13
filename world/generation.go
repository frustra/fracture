package world

import (
	"github.com/frustra/fracture/perlin"
)

func (c *Chunk) Generate(blockType byte) {
	noise := perlin.NewNoise2D(0)

	for z := int64(0); z < BlockZSize; z++ {
		for x := int64(0); x < BlockXSize; x++ {
			absx, absz := float64(x+c.OffsetX*BlockXSize), float64(z+c.OffsetZ*BlockZSize)
			r := noise.At(absx/70, absz/70) + 0.8
			r += (noise.At(absx/20, absz/20) + 0.6) / 3

			height := int64(r * 16 * 3)

			for y := int64(0); y < height; y++ {
				c.Set(x, y, z, 3)
			}
			for y := height; y < 42; y++ {
				c.Set(x, y, z, 9)
			}
			c.Set(x, height, z, blockType)
		}
	}

	c.CalculateLighting()
}

func (c *Chunk) CalculateLighting() {
	for z := int64(0); z < BlockZSize; z++ {
		for x := int64(0); x < BlockXSize; x++ {
			c.CalculateSkyLightingForColumn(x, z)
		}
	}
}

func (c *Chunk) CalculateSkyLightingForColumn(x, z int64) {
	var light byte = 15
	for y := TotalBlockYSize - 1; y >= 0; y-- {
		if light != 0 {
			block := c.Get(x, y, z)
			if block != 0 {
				if block == 9 { // Transparent.
					light -= 3
				} else { // Opaque.
					light = 0
				}
			}
		}

		c.SetSkyLight(x, y, z, light)
	}
}
