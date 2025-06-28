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
	Colors    []string
	ShortRepo srepository.Game
}

func NewGame(shortGameRepo srepository.Game) *Game {
	return &Game{
		ShortRepo: shortGameRepo,
	}
}

func (g *Game) Find(player model.Player) (string, string) {
	var game model.Game
	var participantID string
	for {
		curGame, err := g.ShortRepo.Find(player)

		if err != nil {
			game = model.Game{ID: uuid.NewString(), Status: "waiting", Type: player.GameType}

			game.Players = make(map[model.Player]struct{})
			game.Participants.Data = make(map[string]ctypes.Participant)
			game.Field.Data = make([][]ctypes.Cell, game.Type.FieldSize)
			for i := range game.Field.Data {
				game.Field.Data[i] = make([]ctypes.Cell, game.Type.FieldSize)
				for j := range game.Field.Data[i] {
					game.Field.Data[i][j] = ctypes.Cell{
						CompSize: 0,
						ParticipantID:    "",
					}
				}
			}
			g.ShortRepo.Post(game)
			player.GameID = game.ID
			id, err := g.ShortRepo.Add(player)
			if err == nil {
				participantID = *id
				break
			}
		} else {
			game = *curGame
			player.GameID = game.ID
			id, err := g.ShortRepo.Add(player)
			if err == nil {
				participantID = *id
				break
			}
		}
	}

	go g.broadcast(game, ctypes.ServerEvent{Type: "waiting_change", Data: ctypes.WaitingChange{Waiting: len(game.Players)}})

	if game.Type.Size == len(game.Players) {
		go g.start(game)
	}
	return game.ID, participantID
}

func (g *Game) RemovePlayer(player model.Player) error {
	game, err := g.ShortRepo.Get(player.GameID)
	if err != nil {
		return fmt.Errorf("removing player from game search: %s", err)
	}
	err = g.ShortRepo.Remove(player)
	if err != nil {
		return fmt.Errorf("removing player from game search: %s", err)
	}

	go g.broadcast(*game, ctypes.ServerEvent{Type: "waiting change", Data: ctypes.WaitingChange{Waiting: len(game.Players)}})
	return nil
}

func (g *Game) manageTimer(game model.Game) {
	<-game.Timer.C

	if err := g.finish(game); err != nil {
		log.Printf("%s\n", err)
	} else {
		log.Printf("%s\n", "game finished")
	}
}

func (g *Game) finish(game model.Game) error {
	sGame, err := g.ShortRepo.Get(game.ID)
	if err != nil {
		return fmt.Errorf("error finishing game: %s", err)
	}
	g.broadcast(*sGame, ctypes.ServerEvent{Type: "game_finish", Data: ctypes.GameFinish{Participants: game.Participants}})
	err = g.ShortRepo.Delete(game)
	if err != nil {
		return fmt.Errorf("error finishing game: %s", err)
	}
	return nil
}

func (g *Game) start(game model.Game) {
	game.Status = "started"
	game.Timer = *time.NewTimer(time.Second * time.Duration(game.Type.Time))
	go g.manageTimer(game)
	if err := g.ShortRepo.Put(game); err != nil{
		log.Printf("starting game: %s", err)
	}
	go g.broadcast(game, ctypes.ServerEvent{Type: "game_start", Data: ctypes.GameStart{Participants: game.Participants, Field: game.Field}})
}

func (g *Game) Move(player *model.Player, y, x int) error {
	game, err := g.ShortRepo.Get(player.GameID)
	if time.Now().Before(player.LastMove.Add(time.Duration(game.Type.Cooldown)*time.Second)) {
		return fmt.Errorf("making a move: you are on a cooldown")
	}
	if err != nil {
		return fmt.Errorf("making a move %s", err)
	}
	if game.Status != "started" {
		return usecase.ErrorGameNotStarted
	}
	if (y >= game.Type.FieldSize || y < 0) || (x >= game.Type.FieldSize || x < 0) {
		return usecase.ErrorWrongMoveCoordinates
	}
	game.Field.Data[y][x].ParticipantID = player.ParticipantID

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
				game.Field.Data[ni][nj].ParticipantID != "" && !visited[ni][nj] {
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
	for i, participant := range game.Participants.Data {
		participant.Score = 0
		game.Participants.Data[i] = participant
	}
	for i := range rows {
		for j := range cols {
			id := game.Field.Data[i][j].ParticipantID
			if (id == "") {
				continue
			}
			participant := game.Participants.Data[id]
			participant.Score++
			game.Participants.Data[id] = participant
		}
	}
	if err := g.ShortRepo.Put(*game); err != nil {
		return fmt.Errorf("making a move %s", err)
	}

	player.LastMove = time.Now()

	go g.broadcast(*game, ctypes.ServerEvent{Type: "player_move", Data: ctypes.PlayerMove{Field: game.Field, Participants: game.Participants}})

	return nil
}

func (g *Game) broadcast(game model.Game, event ctypes.ServerEvent) {
	for player := range game.Players {
		player.Send <- event
	}
}
