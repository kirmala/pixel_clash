package usecase

import (
	"pixel_clash/model"
)

type Game interface {
	Find(player model.Player) string
	Move(player model.Player, x, y int) error
	RemovePlayer(player model.Player) error
}