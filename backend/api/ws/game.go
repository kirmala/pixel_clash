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


func (g *Game) Handle(conn *websocket.Conn, player model.Player) {
	player.Send = make(chan ctypes.ServerEvent, 10)

	go func() {
        for msg := range player.Send {
            if err := conn.WriteJSON(msg); err != nil {
                log.Println("Write error:", err)
            }
        }
    }()

    // Main loop: handle player messages
    for {
        var msg types.Request
        if err := conn.ReadJSON(&msg); err != nil {
            conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
            continue
        }

        switch msg.Type {
		case "find_game":
			var req types.FindGameRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
        		conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
            	continue
    		}

			player.GameId = g.Service.Find(player)
			player.GameType = req.GameType
			player.Nickname = req.Nickname

			g.PlayerService.Post(player)

			resp := types.FindGameResponse{GameId: player.GameId, PlayerId: player.Id}
			
			if err := conn.WriteJSON(resp); err != nil {
                log.Println("Write error:", err)
                continue
            }
        case "stop_searching":
			var req types.StopSearchingRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
        		conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
            	continue
    		}
			if err := g.Service.RemovePlayer(req.PlayerId); err != nil {
				conn.WriteJSON(types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "error", Data: err.Error()})
				continue
			}

			if err := g.PlayerService.Delete(req.PlayerId); err != nil {
        		conn.WriteJSON(types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "error", Data: err.Error()})
				continue
    		}

			conn.WriteJSON(types.ServerResponse{Type: "stop_searching_result", ID: msg.ID, Status : "success", Data: types.StopSearchingResponse{Message: "searching stopped succefully"}})
			return
		case "move":
			var req types.MoveRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
        		conn.WriteJSON(types.ServerResponse{Status : "error", Data: "error parsing request"})
            	continue
    		}
			err := g.Service.Move(req.PlayerId, req.X, req.Y)
			if err != nil {
				conn.WriteJSON(types.ServerResponse{Type: "move_result", ID: msg.ID, Status : "error", Data: err.Error()})
            	continue
			}

			conn.WriteJSON(types.ServerResponse{Type: "move_result", ID: msg.ID, Status : "success", Data: types.MoveResponse{Message: "moved successfully"}})
		}
    }
}