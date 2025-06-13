package sram

import (
	"pixel_clash/model"
	"pixel_clash/repository"
	"sync"
)
type Game struct {
	data map[string]model.Game
	mu sync.RWMutex
}

func NewGame() *Game {
	return &Game{data: make(map[string]model.Game)}
}

func (g *Game) Put(game model.Game) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.data[game.Id] = game
}

func (g *Game) Post(game model.Game) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.data[game.Id]; exists {
		return repository.ErrorAlreadyExists
	}
	g.data[game.Id] = game
	return nil
}

func (g *Game) Get(id string) (*model.Game, error){
	g.mu.RLock()
	defer g.mu.RUnlock()
	game, ok := g.data[id]
	if !ok {
		return nil, repository.ErrorKeyNotFound
	}
	return &game, nil
}

func (g *Game) Find(player model.Player) (*model.Game, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, game := range g.data{
		if game.Status == "waiting" && game.Type == player.GameType {
			return &game, nil
		}
	}
	return nil, repository.ErrorWaitingNotFound
}