package usecase

import "pixel_clash/model"

type Player interface {
	Post(player model.Player) error
	Get(id string) (*model.Player, error)
	Delete(id string) error
}