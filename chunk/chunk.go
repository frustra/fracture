package chunk

type Coord struct {
	X, Y, Z uint64
}

type Chunk struct {
	Data [16][16][16]uint8
	Meta [16][16][16]uint8 // 2 elements per byte
	// Send 0xF for light always
}

type Column struct {
	Chunks [16]Chunk
}

func (col *Column) GetChunk(y uint64, chunks *Chunk) error {
	chunks = &col.Chunks[y]
	return nil
}

func (col *Column) GetChunks(i int, chunks *[16]Chunk) error {
	chunks = &col.Chunks
	return nil
}

type Region struct {
	Columns [32][32]Column
}

