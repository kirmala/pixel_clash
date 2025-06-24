package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"pixel_clash/api/ws/types"
	"pixel_clash/ctypes"
	"pixel_clash/model"
	"pixel_clash/usecase"

	"github.com/gorilla/websocket"
)

type Game struct {
	Service usecase.Game
}

func NewGame(gameService usecase.Game) *Game {
	return &Game{
		Service: gameService,
	}
}


func (g *Game) response(player model.Player, msg types.Request, cancel context.CancelFunc, send chan ctypes.ServerMessage) {
	switch msg.Type {
	case "find_game":
		var req types.FindGameRequest

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			send <- ctypes.ServerMessage{Type : "response", Data : types.ServerResponse{Status : "error", Data: "error parsing request"}}
			return
		}
		player.GameType = req.GameType
		player.Nickname = req.Nickname
		player.GameId = g.Service.Find(player)
		send <- ctypes.ServerMessage{
			Type: "response",
			Data: types.ServerResponse{
				Type: "find_game_result", 
				ID: msg.ID,
				Status : "success",
				Data: types.FindGameResponse{Message: "searching stopped succefully"},
			},
		}
	case "stop_searching":
		var req types.StopSearchingRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			send <- ctypes.ServerMessage{Type : "response", Data : types.ServerResponse{Status : "error", Data: "error parsing request"}}
			return
		}
		if err := g.Service.RemovePlayer(player); err != nil {
			send <- ctypes.ServerMessage{Type : "response", Data : types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "error", Data: fmt.Sprintf("stopping search: %s", err)}}
			return
		}
		send <- ctypes.ServerMessage{
			Type: "response",
			Data: types.ServerResponse{
				Type: "stop_searching_result", 
				ID: msg.ID,
				Status : "success",
				Data: types.StopSearchingResponse{Message: "searching stopped succefully"},
			},
		}
		cancel()
	case "move":
		var req types.MoveRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			send <- ctypes.ServerMessage{Type : "response", Data : types.ServerResponse{Status : "error", Data: "error parsing request"}}
			return
		}
		err := g.Service.Move(&player, req.X, req.Y)
		if err != nil {
			send <- ctypes.ServerMessage{Type : "response", Data : types.ServerResponse{Type: "move_result", ID: msg.ID, Status : "error", Data: fmt.Sprintf("making a move: %s", err)}}
			return
		}
		send <- ctypes.ServerMessage{
			Type: "response",
			Data: types.ServerResponse{
				Type: "move_result",
				ID: msg.ID,
				Status : "success",
				Data: types.MoveResponse{Message: "moved successfully"},
			},
		}
	}
}

func (g *Game) event(event ctypes.ServerEvent, cancel context.CancelFunc, send chan ctypes.ServerMessage) {
	msg := ctypes.ServerMessage{Type: "event", Data: event}
	send <- msg
	if msg.Type == "game_finish" {
		cancel()
	}
}


func (g *Game) Handle(player model.Player) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	player.Send = make(chan ctypes.ServerEvent, 10)
	read := make(chan types.Request, 10)

	send := make(chan ctypes.ServerMessage, 10)

	go func() {
        for {
            select {
            case msg := <- send:
                if err := player.Conn.WriteJSON(msg); err != nil {
                    log.Printf("Write error: %s", err)
					cancel()
					return
                }
            case <- ctx.Done():
                return
            }
        }
    }()


	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			default:
				var msg types.Request
				err := player.Conn.ReadJSON(&msg)

				if err != nil {
					if websocket.IsCloseError(err) || 
					errors.Is(err, net.ErrClosed) || 
					errors.Is(err, io.EOF) {
						log.Printf("Read error: %s", err)
						cancel()// Trigger graceful shutdown
						return
					}
					
					send <- ctypes.ServerMessage{
						Type: "response",
						Data: types.ServerResponse{
							Status: "error", 
							Data: err.Error(),
						},
					}
				}
				read <- msg
			}
		}
	}()

	for {
		select {
		case <- ctx.Done():
			return
		case msg := <- read:
			go g.response(player, msg, cancel, send)
		case msg := <- player.Send:
			go g.event(msg, cancel, send)
		}
	}
}