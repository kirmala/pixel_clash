package service

import (
	"pixel_clash/model"
	"pixel_clash/repository"
)

type Player struct {
	Repo repository.Player
}

func NewPlayer(PlayerRepo repository.Player) *Player {
	return &Player{
		Repo:    PlayerRepo,
	}
}

func (u *Player) Post(player model.Player) error {
	return u.Repo.Post(player)
}

