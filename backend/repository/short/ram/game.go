package sram

import (
	"pixel_clash/model"
	"pixel_clash/ctypes"
	"pixel_clash/repository"
	"sync"

	"github.com/google/uuid"
)
type Game struct {
	data map[string]model.Game
	mu sync.RWMutex
}

func NewGame() *Game {
	return &Game{data: make(map[string]model.Game)}
}

func (g *Game) Put(game model.Game) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, ok := g.data[game.ID]
	if !ok {
		return repository.ErrorKeyNotFound
	}
	g.data[game.ID] = game
	return nil
}

func (g *Game) Post(game model.Game) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.data[game.ID]; exists {
		return repository.ErrorAlreadyExists
	}
	g.data[game.ID] = game
	return nil
}

func (g *Game) Get(ID string) (*model.Game, error){
	g.mu.RLock()
	defer g.mu.RUnlock()
	game, ok := g.data[ID]
	if !ok {
		return nil, repository.ErrorKeyNotFound
	}
	return &game, nil
}

func (g *Game) Delete(game model.Game) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, ok := g.data[game.ID]
	if !ok {
		return repository.ErrorKeyNotFound
	}
	delete(g.data, game.ID)
	return nil
}

func (g *Game) Find(player model.Player) (*model.Game, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, game := range g.data{
		if game.Status == "waiting" && game.Type == player.GameType && len(game.Players) != game.Type.Size{
			return &game, nil
		}
	}
	return nil, repository.ErrorWaitingNotFound
}

func (g *Game) Add(player model.Player) (*string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	game, ok := g.data[player.GameID]
	if !ok {
		return nil, repository.ErrorKeyNotFound
	}
	if game.Status == "started" || len(game.Players) == game.Type.Size {
		return nil, repository.ErrorGameAlreadyStarted
	}
	participant := ctypes.Participant{ID : uuid.NewString(), Score : 0, Nickname: player.Nickname}
	game.Participants.Data[participant.ID] = participant
	player.ParticipantID = participant.ID
	game.Players[player] = struct{}{}
	g.data[game.ID] = game
	return &participant.ID, nil
}

func (g *Game) Remove(player model.Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	game, ok := g.data[player.GameID]
	if !ok {
		return repository.ErrorKeyNotFound
	}
	if game.Status == "started" || len(game.Players) == game.Type.Size {
		return repository.ErrorGameAlreadyStarted
	}
	delete(game.Players, player)
	delete(game.Participants.Data, player.ParticipantID)
	g.data[game.ID] = game
	return nil
}