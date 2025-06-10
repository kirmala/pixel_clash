package service

import (
	"pixel_clash/model"
	"pixel_clash/repository"

	"github.com/google/uuid"
)

const (
	rows = 10
	cols = 10
	gameTime = 60
)

type Game struct {
	Repo repository.Game
	PlayerRepo repository.Player
}

func NewGame(gameRepo repository.Game, playerRepo repository.Player) *Game {
	return &Game{
		Repo:    gameRepo,
		PlayerRepo: playerRepo,
	}
}

func (g *Game) Find(player model.Player) {
	game, err := g.Repo.Find(player)

	if err != nil {
		game = &model.Game{Id: uuid.NewString(), Status: "waiting", Players: []model.Player{player}, Capacity: player.GameCapacity}
	} else {
		game.Players = append(game.Players, player)
	}

	g.Repo.Put(*game)
	player.GameId = game.Id
	player.Status = "searching"
	g.PlayerRepo.Put(player)

	if (game.Capacity == len(game.Players)) {
		g.start(*game)
	}
}

func (g *Game) start(game model.Game) {
	game.Status = "started"
	for i, player := range game.Players {
		player.Status = "playing"
		player.Color = i
		g.PlayerRepo.Put(player)
	}
	g.Repo.Put(game)
}

// func (g *Game) Move(game model.Game, User model.User, x, y int) {
// 	game.Feild[y][x] = User.Color


// }