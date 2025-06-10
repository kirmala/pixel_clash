package ram

import (
	"pixel_clash/model"
	"pixel_clash/repository"
)
type Player struct {
	data map[string]model.Player
}

func NewPlayer() *Player {
	return &Player{data: make(map[string]model.Player)}
}

func (p *Player) Post(player model.Player) error {
	if _, exists := p.data[player.Id]; exists {
		return repository.ErrorAlreadyExists
	}
	p.data[player.Id] = player
	return nil
}

func (p *Player) Put(player model.Player) {
	p.data[player.Id] = player
}

func (p *Player) Get(id string) (*model.Player, error){
	player, ok := p.data[id]
	if !ok {
		return nil, repository.ErrorKeyNotFound
	}
	return &player, nil
}


