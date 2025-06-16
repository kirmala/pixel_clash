package service

import (
	"pixel_clash/model"
	srepository "pixel_clash/repository/short"
)

type Player struct {
	Repo srepository.Player
}

func NewPlayer(PlayerRepo srepository.Player) *Player {
	return &Player{
		Repo:    PlayerRepo,
	}
}

func (u *Player) Post(player model.Player) error {
	return u.Repo.Post(player)
}

func (u *Player) Get(id string) (*model.Player, error) {
	return u.Repo.Get(id)
}

func (u *Player) Delete(id string) error {
	return u.Repo.Delete(id)
}



