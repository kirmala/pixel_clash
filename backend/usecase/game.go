package usecase

import "pixel_clash/model"

type Game interface {
	Find(player model.Player)
	Move(player model.Player, x, y int) error
	// Stop(game model.Game) error
}