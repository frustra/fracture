package player

type Player struct {
	Name string
	X, Y, Z float64
	Yaw, Pitch float32
}

type PlayerUpdate struct {
	Id int
	X, Y, Z float64
	Yaw, Pitch float32
}
