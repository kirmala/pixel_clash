package srepository

import "pixel_clash/model"

type Player interface {
	Post(player model.Player) error
	Put(player model.Player)
	Get(id string) (*model.Player, error)
}