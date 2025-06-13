package websocket

import (
	"net/http"
	"pixel_clash/api/websocket/types"
	"pixel_clash/model"
	"pixel_clash/usecase"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Game struct {
	Service usecase.Game
	PlayerService usecase.Player

	upgrader websocket.Upgrader
}


func NewGameWebsocketHandler(playerService usecase.Player, gameService usecase.Game) *Game {
	return &Game{Service: gameService, PlayerService: playerService}
}

func (g *Game) JoinHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := g.upgrader.Upgrade(w, r, nil)

	if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	var joinReq types.JoinHandlerRequest
    if err := conn.ReadJSON(&joinReq); err != nil {
        types.SendError(conn, err)
        return
    }

	player := model.Player{Id : uuid.NewString(), Nickname: joinReq.Nickname, GameType: joinReq.GameType}
	
	gameId := g.Service.Find(player)
	player.GameId = gameId

	if err = g.PlayerService.Post(player); err != nil {
        types.SendError(conn, err)
        return
    }
	

	types.SendResponse(
		conn,
		types.JoinHandlerResponse{
			GameId:   player.GameId,
			PlayerId: player.Id,
    	},
	)

	for {
		var req types.MoveRequest
		if err := conn.ReadJSON(&req); err != nil {
			types.SendError(conn, err)
			return
		}

		_, err := g.PlayerService.Get(req.PlayerId)

		if err != nil {
			types.SendError(conn, err)
		}

		g.Service.Move(req.PlayerId, req.X, req.Y)
	}
}


func (g *Game) WithGameHandlers(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/join", g.JoinHandler)
	})
}
