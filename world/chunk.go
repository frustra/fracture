package world

import (
	"bytes"
	"compress/zlib"
	"io"
)

const (
	ChunkWidthPerNode = 8 // Square side length of a chunk node.

	BlockHeight     = 16 // Subchunks per chunk
	BlockYSize      = 16 // Y blocks per chunk
	BlockXSize      = 16 // X blocks per chunk
	BlockZSize      = 16 // Z blocks per chunk
	BlockArraySize  = BlockHeight * BlockYSize * BlockZSize * BlockXSize
	TotalBlockYSize = int64(BlockYSize * BlockHeight)
)

func ChunkCoordsToNode(x, z int64) (nx, nz int64) {
	if x < 0 {
		x += 1 - ChunkWidthPerNode
	}
	if z < 0 {
		z += 1 - ChunkWidthPerNode
	}
	nx = (x / ChunkWidthPerNode) * ChunkWidthPerNode
	nz = (z / ChunkWidthPerNode) * ChunkWidthPerNode
	return nx, nz
}

func WorldCoordsToChunk(x, z int64) (cx, cz int64) {
	if x < 0 {
		x += 1 - BlockXSize
	}
	if z < 0 {
		z += 1 - BlockZSize
	}
	cx = x / BlockXSize
	cz = z / BlockZSize
	return cx, cz
}

func WorldCoordsToChunkInternal(x, z int64) (cx, cz int64) {
	cx = x % BlockXSize
	cz = z % BlockZSize
	if cx < 0 {
		cx += BlockXSize
	}
	if cz < 0 {
		cz += BlockZSize
	}
	return cx, cz
}

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

func indexOf(x, y, z int64) int {
	return int(y*BlockZSize*BlockXSize + z*BlockXSize + x)
}

func setInHalfArray(arr []byte, x, y, z int64, val byte) {
	index := indexOf(x, y, z)
	shift := uint(index%2) * 4
	index /= 2

	arr[index] = (arr[index] & (0xf0 >> shift)) | (val << shift)
}

func (c *Chunk) Set(x, y, z int64, val byte) {
	c.BlockTypes[indexOf(x, y, z)] = val
}

func (c *Chunk) Get(x, y, z int64) byte {
	return c.BlockTypes[indexOf(x, y, z)]
}

func (c *Chunk) SetMetadata(x, y, z int64, val byte) {
	setInHalfArray(c.BlockMeta[:], x, y, z, val)
}

func (c *Chunk) SetBlockLight(x, y, z int64, val byte) {
	setInHalfArray(c.BlockLight[:], x, y, z, val)
}

func (c *Chunk) SetSkyLight(x, y, z int64, val byte) {
	setInHalfArray(c.SkyLight[:], x, y, z, val)
}

func (c *Chunk) Clear() {
	for i := 0; i < BlockArraySize; i++ {
		c.BlockTypes[i] = 0
	}
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
