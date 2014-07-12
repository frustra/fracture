package chunk

func ChunkCoordsToNode(x, z int64) (nx, nz int64) {
	if x < 0 {
		x -= ChunkWidthPerNode
	}
	if z < 0 {
		z -= ChunkWidthPerNode
	}
	nx = (x / ChunkWidthPerNode) * ChunkWidthPerNode
	nz = (z / ChunkWidthPerNode) * ChunkWidthPerNode
	return nx, nz
}
