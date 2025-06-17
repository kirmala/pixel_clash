package ws

import (
	"encoding/json"
	"log"
	"pixel_clash/api/ws/types"
	"pixel_clash/ctypes"
	"pixel_clash/model"
	"pixel_clash/usecase"

	"github.com/gorilla/websocket"
)

type Game struct {
	Service usecase.Game
	PlayerService usecase.Player
}

func NewGame(gameService usecase.Game, playerService usecase.Player) *Game {
	return &Game{
		Service: gameService,
		PlayerService: playerService,
	}
}

func (g *Game) answer(player model.Player, msg types.Request, conn *websocket.Conn, done chan struct{}) {
	switch msg.Type {
	case "find_game":
		var req types.FindGameRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
			return
		}

		player.GameId = g.Service.Find(player)
		player.GameType = req.GameType
		player.Nickname = req.Nickname

		g.PlayerService.Post(player)

		resp := types.FindGameResponse{GameId: player.GameId, PlayerId: player.Id}
		
		if err := conn.WriteJSON(resp); err != nil {
			log.Println("Write error:", err)
			return
		}
	case "stop_searching":
		var req types.StopSearchingRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
			return
		}
		if err := g.Service.RemovePlayer(player.Id); err != nil {
			conn.WriteJSON(types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "error", Data: err.Error()})
			return
		}

		if err := g.PlayerService.Delete(player.Id); err != nil {
			conn.WriteJSON(types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "error", Data: err.Error()})
			return
		}

		conn.WriteJSON(types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "success", Data: types.StopSearchingResponse{Message: "searching stopped succefully"}})
		done <- struct{}{}
		return
	case "move":
		var req types.MoveRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
			return
		}
		err := g.Service.Move(player.Id, req.X, req.Y)
		if err != nil {
			conn.WriteJSON(types.ServerResponse{Type: "move_result", ID: msg.ID, Status : "error", Data: err.Error()})
			return
		}

		conn.WriteJSON(types.ServerResponse{Type: "move_result", ID: msg.ID, Status : "success", Data: types.MoveResponse{Message: "moved successfully"}})
	}
}

func (g *Game) send(msg ctypes.ServerEvent, conn *websocket.Conn, done chan struct{}) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Write error:", err)
		return
	}
	if msg.Type == "game_finish" {
		done <- struct{}{}
	}
}


func (g *Game) Handle(conn *websocket.Conn, player model.Player) {
	player.Send = make(chan ctypes.ServerEvent, 10)

	read := make(chan types.Request)
	done := make(chan struct{})

	go func() {
		for {
			var msg types.Request
			if err := conn.ReadJSON(&msg); err != nil {
				conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
				return
			}
			read <- msg
		}
	}()
	for {
		select {
		case <- done:
			return
		case msg := <- read:
			go g.answer(player, msg, conn, done)
		case msg := <- player.Send:
			go g.send(msg, conn, done)
		}
	}
}