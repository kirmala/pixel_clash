package http

import (
	"log"
	"net/http"
	"pixel_clash/api/ws"
	"pixel_clash/model"
	"sync"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	gameConnection ws.Game
	upgrader       websocket.Upgrader
}

func NewUserHandler(gameConnection ws.Game) *User {
	return &User{
        gameConnection: gameConnection,
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
				return true
			},
        },
    }
}

func (u *User) JoinHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := u.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
        return
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	player := model.Player{ID: uuid.NewString(), Conn: conn, Lock: &sync.Mutex{}}
	u.gameConnection.Handle(player)
}

// WithUserHandlers registers user-related HTTP handlers.
func (u *User) WithUserHandlers(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/join", u.JoinHandler)
	})
}
