package chunk

import (
	"bytes"
	"compress/zlib"
	"io"
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
	c := &Chunk{
		OffsetX: offsetX,
		OffsetZ: offsetZ,
	}
	c.Generate()
	return c
}

func (c *Chunk) Set(x, y, z uint8, val byte) {
	c.BlockTypes[y*BlockZSize+z*BlockXSize+x] = val
}

func (c *Chunk) Get(x, y, z uint8) byte {
	return c.BlockTypes[y*BlockZSize+z*BlockXSize+x]
}

func (c *Chunk) Clear() {
	for i := 0; i < BlockArraySize; i++ {
		c.BlockTypes[i] = 0
	}
}

func (c *Chunk) Generate() {
	for y := uint8(0); y < BlockYSize; y++ {
		for z := uint8(0); z < BlockZSize; z++ {
			for x := uint8(0); x < BlockXSize; x++ {
				if y < 100 {
					c.Set(x, y, z, 3)
				} else if y == 100 {
					c.Set(x, y, z, 2)
				}
			}
		}
	}
}

func (c *Chunk) Cube(y uint8) []byte {
	start := y * BlockZSize * BlockXSize
	end := (y + 1) * BlockZSize * BlockXSize
	return c.BlockTypes[start:end]
}

func (c *Chunk) MarshallCompressed() []byte {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	c.WriteTo(w)
	w.Close()
	return compressed.Bytes()
}

func (c *Chunk) UnmarshallCompressed(buf []byte) error {
	w, err := zlib.NewReader(bytes.NewReader(buf))
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(c.BlockTypes[:])
	_, err = io.Copy(b, w)
	return err
}

func (c *Chunk) WriteTo(w io.Writer) (int, error) {
	return w.Write(c.BlockTypes[:])
}
