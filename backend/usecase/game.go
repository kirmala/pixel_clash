package usecase

import (
	"pixel_clash/model"
)

type Game interface {
	Find(player model.Player) string
	Move(playerId string, x, y int) error
}