package chunk

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/frustra/fracture/perlin"
)

const (
	BlockHeight     = 16 // Subchunks per chunk
	BlockYSize      = 16 // Y blocks per chunk
	BlockXSize      = 16 // X blocks per chunk
	BlockZSize      = 16 // Z blocks per chunk
	BlockArraySize  = BlockHeight * BlockYSize * BlockZSize * BlockXSize
	TotalBlockYSize = BlockYSize * BlockHeight
)

// A Chunk column really.
type Chunk struct {
	OffsetX, OffsetZ int64

	BlockTypes [BlockArraySize]byte
	BlockLight [BlockArraySize / 2]byte
	BlockMeta  [BlockArraySize / 2]byte
	SkyLight   [BlockArraySize / 2]byte
}

func NewChunk(offsetX, offsetZ int64) *Chunk {
	return &Chunk{
		OffsetX: offsetX,
		OffsetZ: offsetZ,
	}
}

func indexOf(x, y, z int) int {
	return y*BlockZSize*BlockXSize + z*BlockXSize + x
}

func setInHalfArray(arr []byte, x, y, z int, val byte) {
	index := indexOf(x, y, z)
	shift := uint(index%2) * 4
	index /= 2

	arr[index] = (arr[index] & (0xf0 >> shift)) | (val << shift)
}

func (c *Chunk) Set(x, y, z int, val byte) {
	c.BlockTypes[indexOf(x, y, z)] = val
}

func (c *Chunk) Get(x, y, z int) byte {
	return c.BlockTypes[indexOf(x, y, z)]
}

func (c *Chunk) SetBlockLight(x, y, z int, val byte) {
	setInHalfArray(c.BlockLight[:], x, y, z, val)
}

func (c *Chunk) SetSkyLight(x, y, z int, val byte) {
	setInHalfArray(c.SkyLight[:], x, y, z, val)
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

	for z := 0; z < BlockZSize; z++ {
		for x := 0; x < BlockXSize; x++ {
			var light byte = 15
			for y := TotalBlockYSize - 1; y >= 0; y-- {
				c.SetSkyLight(x, y, z, light)

				if light == 0 {
					continue
				}

				block := c.Get(x, y, z)
				if block != 0 {
					if block == 9 { // Transparent.
						light--
					} else { // Opaque.
						light = 0
					}
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
	w.Close()
	return compressed.Bytes()
}

func (c *Chunk) UnmarshallCompressed(buf []byte) error {
	r, err := zlib.NewReader(bytes.NewReader(buf))
	if err != nil {
		return err
	}
	return c.ReadFrom(r)
}

func (c *Chunk) ReadFrom(r io.Reader) error {
	b := bytes.NewBuffer(c.BlockTypes[0:0])
	_, err := io.Copy(b, r)
	if err != nil {
		return err
	}

	b = bytes.NewBuffer(c.BlockMeta[0:0])
	_, err = io.Copy(b, r)
	if err != nil {
		return err
	}

	b = bytes.NewBuffer(c.BlockLight[0:0])
	_, err = io.Copy(b, r)
	if err != nil {
		return err
	}

	b = bytes.NewBuffer(c.SkyLight[0:0])
	_, err = io.Copy(b, r)
	return err
}

func (c *Chunk) WriteTo(w io.Writer) (int, error) {
	count, err := w.Write(c.BlockTypes[:])
	if err != nil {
		return count, err
	}

	n, err := w.Write(c.BlockMeta[:])
	count += n
	if err != nil {
		return count, err
	}

	n, err = w.Write(c.BlockLight[:])
	count += n
	if err != nil {
		return count, err
	}

	n, err = w.Write(c.SkyLight[:])
	count += n
	return count, err
}
