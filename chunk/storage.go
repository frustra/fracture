package chunk

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/frustra/fracture/perlin"
)

const (
	BlockHeight    = 16 // Subchunks per chunk
	BlockYSize     = 16 // Y blocks per chunk
	BlockXSize     = 16 // X blocks per chunk
	BlockZSize     = 16 // Z blocks per chunk
	BlockArraySize = BlockHeight * BlockYSize * BlockZSize * BlockXSize
)

// A Chunk column really.
type Chunk struct {
	OffsetX, OffsetZ int64

	BlockTypes [BlockArraySize]byte
}

func NewChunk(offsetX, offsetZ int64) *Chunk {
	return &Chunk{
		OffsetX: offsetX,
		OffsetZ: offsetZ,
	}
}

func (c *Chunk) Set(x, y, z int, val byte) {
	c.BlockTypes[y*BlockZSize*BlockXSize+z*BlockXSize+x] = val
}

func (c *Chunk) Get(x, y, z int) byte {
	return c.BlockTypes[y*BlockZSize*BlockXSize+z*BlockXSize+x]
}

func (c *Chunk) Clear() {
	for i := 0; i < BlockArraySize; i++ {
		c.BlockTypes[i] = 0
	}
}

func (c *Chunk) Generate(blockType byte) {
	noise := perlin.NewNoise2D(0)

	for z := 0; z < BlockZSize; z++ {
		for x := 0; x < BlockXSize; x++ {
			absx, absz := float64(x+int(c.OffsetX*BlockXSize)), float64(z+int(c.OffsetZ*BlockZSize))
			r := noise.At(absx/70, absz/70) + 0.8
			r += (noise.At(absx/20, absz/20) + 0.6) / 3

			height := int(r * 16 * 3)

			for y := 0; y < height; y++ {
				c.Set(x, y, z, 3)
			}
			for y := height; y < 42; y++ {
				c.Set(x, y, z, 9)
			}
			c.Set(x, height, z, blockType)
		}
	}
}

func (c *Chunk) Cube(y int) []byte {
	start := y * BlockZSize * BlockXSize
	end := (y + 1) * BlockZSize * BlockXSize
	return c.BlockTypes[start:end]
}

func (c *Chunk) MarshallCompressed() []byte {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	c.WriteTo(w)
	meta := make([]byte, 16*16*16*16/2)
	w.Write(meta)
	light := make([]byte, 16*16*16*16/2)
	for i := 0; i < len(light); i++ {
		light[i] = 15 | 15<<4
	}
	w.Write(light)
	w.Write(light)
	w.Close()
	return compressed.Bytes()
}

func (c *Chunk) UnmarshallCompressed(buf []byte) error {
	w, err := zlib.NewReader(bytes.NewReader(buf))
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(c.BlockTypes[0:0])
	_, err = io.Copy(b, w)
	return err
}

func (c *Chunk) WriteTo(w io.Writer) (int, error) {
	return w.Write(c.BlockTypes[:])
}
