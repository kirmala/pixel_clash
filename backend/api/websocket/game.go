package websocket

import (
	"net/http"
	"pixel_clash/usecase"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Game struct {
	service usecase.Game
	playerService usecase.Player
}

func NewGame(service usecase.Game, playerService usecase.Player) *Game {
	return &Game{
		service: service,
		playerService: playerService,
	}
}

func (g *Game) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	playerID := chi.URLParam(r, "playerID")
    if playerID == "" {
        http.Error(w, "Player ID is required", http.StatusBadRequest)
        return
    }

	player, err := g.playerService.Get(playerID)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer func() {
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		
		x, y := parse(message)

		g.service.Move(player, x, y)
	}
}