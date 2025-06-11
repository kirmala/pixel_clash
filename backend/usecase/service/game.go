package service

import (
	"go/constant"
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

type coordiante struct {
	x, y int
}

func (g *Game) Move(player model.Player, x, y int) {
	game, _ := g.Repo.Get(player.GameId)
	game.Feild[y][x] = player.Color

    rows, cols := len(game.Feild), len(game.Feild[0])
    visited := make([][]bool, rows)
    for i := range visited {
        visited[i] = make([]bool, cols)
    }
    
    directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	var component [][2]int
    queue := [][2]int{{y, x}}
    visited[y][x] = true
                
	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]
		component = append(component, cell)
		
		for _, dir := range directions {
			ni, nj := cell[0]+dir[0], cell[1]+dir[1]
			if ni >= 0 && ni < rows && nj >= 0 && nj < cols &&
				game.Feild[ni][nj] != 0 && !visited[ni][nj] {
				visited[ni][nj] = true
				queue = append(queue, [2]int{ni, nj})
			}
		}
	}
                
    // Remove if large enough
	if len(component) >= game.ThreasholdSquare {
		for _, cell := range component {
			game.Feild[cell[0]][cell[1]] = 0
			game.ColorToPlayer[game.Feild[cell[0]][cell[1]]]--
		}
	}
}