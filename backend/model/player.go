package model

import (
	"pixel_clash/ctypes"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id         string
	UserId     string
	Nickname   string
	GameType   ctypes.Game
	GameId     string
	Connection *websocket.Conn
}
