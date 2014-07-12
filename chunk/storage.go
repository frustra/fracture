package chunk

// A Chunk column really.
type Chunk struct {
	BlockTypes [256][16][16]byte // [Y][X][Z]
}

func (c *Chunk) Clear() {
	for y, slice := range c.BlockTypes {
		for x, row := range slice {
			for z, _ := range row {
				c.BlockTypes[y][x][z] = 0
			}
		}
	}
}

func (c *Chunk) Generate() {
	for y, slice := range c.BlockTypes {
		for x, row := range slice {
			for z, _ := range row {
				if y < 100 {
					c.BlockTypes[y][x][z] = 3
				} else if y == 100 {
					c.BlockTypes[y][x][z] = 2
				}
			}
		}
	}
}
