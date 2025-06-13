package sram

import (
	"pixel_clash/model"
	"pixel_clash/repository"
	"sync"
)
type Player struct {
	data map[string]model.Player
	mu sync.RWMutex
}

func NewPlayer() *Player {
	return &Player{data: make(map[string]model.Player)}
}

func (p *Player) Post(player model.Player) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.data[player.Id]; exists {
		return repository.ErrorAlreadyExists
	}
	p.data[player.Id] = player
	return nil
}

func (p *Player) Put(player model.Player) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[player.Id] = player
}

func (p *Player) Get(id string) (*model.Player, error){
	p.mu.RLock()
	defer p.mu.RUnlock()
	player, ok := p.data[id]
	if !ok {
		return nil, repository.ErrorKeyNotFound
	}
	return &player, nil
}


