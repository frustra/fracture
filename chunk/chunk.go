package chunk

import (
)

// For RPC
type Block struct {
	X, Y, Z uint8
	Id uint8
}

// Actual types
type Chunk struct {
	Data [16][16][16]uint8
	//Meta [16][16][16]uint8 // 2 elements per byte
	// Send 0xF for light always
}

type Column struct {
	chunks [16]*Chunk
}

func (col *Column) GetChunk(y uint, chunk **Chunk) error {
	*chunk = col.chunks[y]
	return nil
}

func (col *Column) GetChunks(unused bool, chunks *[16]*Chunk) error {
	chunks = &col.chunks
	return nil
}

func (col *Column) SetBlock(b Block, unused *bool) error {
	col.chunks[b.Y >> 4].Data[b.X][b.Y % 16][b.Z] = b.Id
	return nil
}
