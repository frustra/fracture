package chunk_test

import (
	"testing"

	. "github.com/frustra/fracture/chunk"
)

var chunkToNodeTests = map[int64]int64{
	0:  0,
	5:  0,
	11: 8,
	-1: -8,
	-8: -8,
	-9: -16,
}

func TestChunkCoordsToNode(t *testing.T) {
	for m, expected := range chunkToNodeTests {
		got1, got2 := ChunkCoordsToNode(m, m)

		if got1 != got2 || got1 != expected {
			t.Errorf("       got: %d, %d\n  expected: %d, %d", got1, got2, expected, expected)
		}
	}
}

var worldToChunkTests = map[int64]int64{
	0:   0,
	15:  0,
	16:  1,
	22:  1,
	-1:  -1,
	-15: -1,
	-16: -1,
	-17: -2,
}

func TestWorldCoordsToChunk(t *testing.T) {
	for m, expected := range worldToChunkTests {
		got1, got2 := WorldCoordsToChunk(m, m)

		if got1 != got2 || got1 != expected {
			t.Errorf("       got: %d, %d\n  expected: %d, %d", got1, got2, expected, expected)
		}
	}
}

var worldToChunkInternalTests = map[int64]int64{
	0:   0,
	15:  15,
	16:  0,
	22:  6,
	-1:  15,
	-15: 1,
	-16: 0,
	-17: 15,
}

func TestWorldCoordsToChunkInternal(t *testing.T) {
	for m, expected := range worldToChunkInternalTests {
		got1, got2 := WorldCoordsToChunkInternal(m, m)

		if got1 != got2 || got1 != expected {
			t.Errorf("       got: %d, %d\n  expected: %d, %d", got1, got2, expected, expected)
		}
	}
}
