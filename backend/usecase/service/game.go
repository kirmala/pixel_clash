package service

import (
	"fmt"
	"pixel_clash/api/websocket/types"
	"pixel_clash/ctypes"
	"pixel_clash/model"
	srepository "pixel_clash/repository/short"

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
		game = &model.Game{Id: uuid.NewString(), Status: "waiting", PlayerIds: []string{player.Id}, Type : player.GameType}

		game.Feild = make([][]ctypes.Cell, game.Type.FeildSize)
		for i := range game.Feild {
			game.Feild[i] = make([]ctypes.Cell, game.Type.FeildSize)
			for j := range game.Feild[i] {
				game.Feild[i][j] = ctypes.Cell{
					CompSize: 0,
					Color: "",
				}
			}
		}
	} else {
		game.PlayerIds = append(game.PlayerIds, player.Id)
	}

	g.ShortRepo.Put(*game)

	if game.Type.FeildSize == len(game.PlayerIds) {
		g.start(*game)
	}

	return game.Id
}

func (g *Game) start(game model.Game) {
	game.Status = "started"
	g.ShortRepo.Put(game)
	for _, id := range game.PlayerIds {
		player, _ := g.PlayerRepo.Get(id)
		player.Connection.WriteJSON(game)
	}
}

func (g *Game) Move(playerId string, x, y int) error {
	player, _ := g.PlayerRepo.Get(playerId)
	game, _ := g.ShortRepo.Get(player.GameId)
	if (y >= game.Type.FeildSize || y < 0) || (x >= game.Type.FeildSize || x < 0) {
		return fmt.Errorf("move coordinates out of bound")
	}
	game.Feild[y][x].Color = player.Id

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
				game.Feild[ni][nj].Color != "" && !visited[ni][nj] {
				visited[ni][nj] = true
				queue = append(queue, [2]int{ni, nj})
			}
		}
	}

	for _, cell := range component {
		game.Feild[cell[0]][cell[1]].CompSize = len(component)
	}
                
	if len(component) >= game.Type.ThreasholdSqare {
		for _, cell := range component {
			game.Feild[cell[0]][cell[1]] = ctypes.Cell{}
		}
	}

	for _, playerId := range game.PlayerIds {
		broadcastPlayer, _ := g.PlayerRepo.Get(playerId)
		types.SendResponse(broadcastPlayer.Connection, game)
	}

	return nil
}