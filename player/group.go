package player

import (
	"errors"
)

type Group struct {
	players    []*Player
	maxPlayers int
}

func (group *Group) GetPlayer(i int, p **Player) error {
	*p = group.players[i]
	return nil
}

func (group *Group) GetPlayers(i int, p *[]*Player) error {
	*p = group.players
	return nil
}

func (group *Group) AddPlayer(p *Player, r *bool) error {
	if len(group.players) >= group.maxPlayers {
		*r = true
		return errors.New("too many players")
	}
	group.players = append(group.players, p)
	return nil
}

func (group *Group) UpdatePlayer(u PlayerUpdate, r *bool) error {
	if u.Id >= len(group.players) {
		*r = true
		return errors.New("invalid player id")
	}
	group.players[u.Id].X = u.X
	group.players[u.Id].Y = u.Y
	group.players[u.Id].Z = u.Z
	group.players[u.Id].Yaw = u.Yaw
	group.players[u.Id].Pitch = u.Pitch
	return nil
}
