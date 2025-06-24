package model

import (
	"pixel_clash/ctypes"
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id         string
	UserId     string
	Nickname   string
	GameType   ctypes.Game
	GameId     string
	Conn *websocket.Conn
	Lock *sync.Mutex
	Send chan ctypes.ServerEvent
}
