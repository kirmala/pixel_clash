package srepository

import "pixel_clash/model"

type Game interface {
	Put(game model.Game) error
	Post(game model.Game) error
	Delete(game model.Game) error
	Get(gameID string) (*model.Game, error)

	Find(player model.Player) (*model.Game, error)
	Add(player model.Player) (*string, error)
	Remove(player model.Player) error
}