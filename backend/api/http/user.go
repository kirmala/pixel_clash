package http

import (
	"net/http"
	"pixel_clash/api/ws"
	"pixel_clash/model"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	gameConnection ws.Game
	upgrader websocket.Upgrader
}

func NewUserHandler(gameConnection ws.Game) *User {
	return &User{gameConnection: gameConnection}
}

func (u *User) JoinHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := u.upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	player := model.Player{Id: uuid.NewString()}
	u.gameConnection.Handle(conn, player)
}

// WithUserHandlers registers user-related HTTP handlers.
func (u *User) WithUserHandlers(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/join", u.JoinHandler)
	})
}
