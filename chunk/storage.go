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
	for y := 0; y < BlockYSize*BlockHeight; y++ {
		for z := 0; z < BlockZSize; z++ {
			for x := 0; x < BlockXSize; x++ {
				if y < 100 {
					c.Set(x, y, z, 3)
				} else if y == 100 {
					c.Set(x, y, z, blockType)
				}
			}
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
