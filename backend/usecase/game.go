package usecase

import "pixel_clash/model"

type Game interface {
	Find(player model.Player) (string, string)
	// Move(game model.Game, x, y int, userId string) error
	// Stop(game model.Game) error
}