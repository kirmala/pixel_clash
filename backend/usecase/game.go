package usecase

import (
	"pixel_clash/model"
)

type Game interface {
	Find(player model.Player) (string, string)
	Move(player *model.Player, y, x int) error
	RemovePlayer(player model.Player) error
}