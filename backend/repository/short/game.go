package srepository

import "pixel_clash/model"

type Game interface {
	Put(game model.Game)
	Post(game model.Game) error
	Get(gameId string) (*model.Game, error)
	Find(player model.Player) (*model.Game, error)
}