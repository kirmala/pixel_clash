package service

import (
	"fmt"
	"pixel_clash/ctypes"
	"pixel_clash/model"
	srepository "pixel_clash/repository/short"
	"time"

	"github.com/google/uuid"
)



type Game struct {
	ShortRepo srepository.Game
	PlayerRepo srepository.Player
}

func NewGame(shortGameRepo srepository.Game, playerRepo srepository.Player) *Game {
	return &Game{
		ShortRepo:  shortGameRepo,
		PlayerRepo: playerRepo,
	}
}

func (g *Game) Find(player model.Player) string {
	game, err := g.ShortRepo.Find(player)

	if err != nil {
		game = &model.Game{Id: uuid.NewString(), Status: "waiting", PlayerIds: make(map[string]struct{}), Type : player.GameType, Started: make(chan struct{}, 1)}
		game.PlayerIds[player.Id] = struct{}{}

		game.Feild.Data = make([][]ctypes.Cell, game.Type.FeildSize)
		for i := range game.Feild.Data {
			game.Feild.Data[i] = make([]ctypes.Cell, game.Type.FeildSize)
			for j := range game.Feild.Data[i] {
				game.Feild.Data[i][j] = ctypes.Cell{
					CompSize: 0,
					Color: "",
				}
			}
		}
	} else {
		game.PlayerIds[player.Id] = struct{}{}
	}

	g.ShortRepo.Put(*game)

	if game.Type.FeildSize == len(game.PlayerIds) {
		go g.start(*game)
	}

	return game.Id
}

func (g *Game) RemovePlayer(playerId string) error {
	player, _ := g.PlayerRepo.Get(playerId)
	game, err := g.ShortRepo.Get(player.GameId)

	if err != nil {
		return fmt.Errorf("error removing player from game search: %s", err)
	}

	delete(game.PlayerIds, player.Id)
	if err := g.ShortRepo.Put(*game); err != nil {
		return fmt.Errorf("error removing player from game search: %s", err)
	}
	return nil
}

func (g *Game) start(game model.Game) {
	game.Status = "started"
	game.Timer = *time.NewTimer(time.Second*time.Duration(game.Type.Time))
	g.ShortRepo.Put(game)
	g.broadcast(game, ctypes.ServerEvent{Type: "game_found"})
}

func (g *Game) Move(playerId string, x, y int) error {
	player, _ := g.PlayerRepo.Get(playerId)
	game, _ := g.ShortRepo.Get(player.GameId)
	if game.Status != "started" {
		return fmt.Errorf("error making a move: game did not start yet")
	}
	if (y >= game.Type.FeildSize || y < 0) || (x >= game.Type.FeildSize || x < 0) {
		return fmt.Errorf("error making a move: coordinates out of bound")
	}
	game.Feild.Data[y][x].Color = player.Id

    rows, cols := len(game.Feild.Data), len(game.Feild.Data[0])
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
				game.Feild.Data[ni][nj].Color != "" && !visited[ni][nj] {
				visited[ni][nj] = true
				queue = append(queue, [2]int{ni, nj})
			}
		}
	}

	for _, cell := range component {
		game.Feild.Data[cell[0]][cell[1]].CompSize = len(component)
	}
                
	if len(component) >= game.Type.ThreasholdSqare {
		for _, cell := range component {
			game.Feild.Data[cell[0]][cell[1]] = ctypes.Cell{}
		}
	}

	g.broadcast(*game, ctypes.ServerEvent{Type: "player_move", Data: game.Feild})

	return nil
}

func (g *Game) broadcast(game model.Game, event ctypes.ServerEvent) {
	for playerId := range game.PlayerIds {
		player, _ := g.PlayerRepo.Get(playerId)
		player.Send <- event
	}
}