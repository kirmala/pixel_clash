package usecase

import "pixel_clash/model"

type Player interface {
	Post(model.Player) error
	Get(string) (*model.Player, error)
}