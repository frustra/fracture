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
