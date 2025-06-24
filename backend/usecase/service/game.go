package service

import (
	"fmt"
	"log"
	"pixel_clash/ctypes"
	"pixel_clash/model"
	srepository "pixel_clash/repository/short"
	"pixel_clash/usecase"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ShortRepo srepository.Game
}

func NewGame(shortGameRepo srepository.Game) *Game {
	return &Game{
		ShortRepo: shortGameRepo,
	}
}

func (g *Game) Find(player model.Player) string {
	game, err := g.ShortRepo.Find(player)

	if err != nil {
		game = &model.Game{Id: uuid.NewString(), Status: "waiting", Players: make(map[model.Player]struct{}), Type: player.GameType}
		game.Players[player] = struct{}{}

		game.Field.Data = make([][]ctypes.Cell, game.Type.FieldSize)
		for i := range game.Field.Data {
			game.Field.Data[i] = make([]ctypes.Cell, game.Type.FieldSize)
			for j := range game.Field.Data[i] {
				game.Field.Data[i][j] = ctypes.Cell{
					CompSize: 0,
					Color:    "",
				}
			}
		}
		g.ShortRepo.Post(*game)
	} else {
		game.Players[player] = struct{}{}
		g.ShortRepo.Put(*game)
	}

	go g.broadcast(*game, ctypes.ServerEvent{Type : "waiting_change", Data : ctypes.WaitingChange{Waiting: len(game.Players)}})

	if game.Type.FieldSize == len(game.Players) {
		go g.start(*game)
	}

	return game.Id
}

func (g *Game) RemovePlayer(player model.Player) error {
	game, err := g.ShortRepo.Get(player.GameId)

	if err != nil {
		return fmt.Errorf("removing player from game search: %s", err)
	}

	delete(game.Players, player)
	if err := g.ShortRepo.Put(*game); err != nil {
		return fmt.Errorf("removing player from game search: %s", err)
	}

	go g.broadcast(*game, ctypes.ServerEvent{Type : "waiting change", Data : ctypes.WaitingChange{Waiting: len(game.Players)}})
	return nil
}

func (g *Game) manageTimer(game model.Game) {
	defer game.Timer.Stop()
	<-game.Timer.C

	if err := g.finish(game); err != nil {
		log.Printf("%s\n", err)
	}
}

func (g *Game) finish(game model.Game) error {
	sGame, err := g.ShortRepo.Get(game.Id)
	if err != nil {
		return fmt.Errorf("error finishing game: %s", err)
	}
	g.broadcast(*sGame, ctypes.ServerEvent{Type: "game_finished", Data: ctypes.GameFinish{Scores: game.Scores}})
	g.ShortRepo.Delete(game)
	return nil
}

func (g *Game) start(game model.Game) {
	game.Status = "started"
	game.Timer = *time.NewTimer(time.Second * time.Duration(game.Type.Time))
	go g.manageTimer(game)
	g.ShortRepo.Put(game)
	go g.broadcast(game, ctypes.ServerEvent{Type: "game_found"})
}

func (g *Game) Move(player model.Player, x, y int) error {
	game, err := g.ShortRepo.Get(player.GameId)
	if err != nil {
		return fmt.Errorf("making a move %s", err)
	}
	if game.Status != "started" {
		return usecase.ErrorGameNotStarted
	}
	if (y >= game.Type.FieldSize || y < 0) || (x >= game.Type.FieldSize || x < 0) {
		return usecase.ErrorWrongMoveCoordinates
	}
	game.Field.Data[y][x].Color = player.Id

	rows, cols := len(game.Field.Data), len(game.Field.Data[0])
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
				game.Field.Data[ni][nj].Color != "" && !visited[ni][nj] {
				visited[ni][nj] = true
				queue = append(queue, [2]int{ni, nj})
			}
		}
	}

	for _, cell := range component {
		game.Field.Data[cell[0]][cell[1]].CompSize = len(component)
	}

	if len(component) >= game.Type.ThresholdSqare {
		for _, cell := range component {
			game.Field.Data[cell[0]][cell[1]] = ctypes.Cell{}
		}
	}
	for i := range rows {
		for j := range cols {
			game.Scores.Data[game.Field.Data[i][j].Color]++
		}
	}
	if err := g.ShortRepo.Put(*game); err != nil {
		return fmt.Errorf("making a move %s", err)
	}

	go g.broadcast(*game, ctypes.ServerEvent{Type: "player_move", Data: ctypes.PlayerMove{Field: game.Field, Scores: game.Scores}})

	return nil
}

func (g *Game) broadcast(game model.Game, event ctypes.ServerEvent) {
	for player := range game.Players {
		player.Send <- event
	}
}
