package chunk

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
