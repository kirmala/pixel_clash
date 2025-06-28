package model

import (
	"pixel_clash/ctypes"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID         string
	UserID     string
	Nickname   string
	GameType   ctypes.Game
	GameID     string
	Conn *websocket.Conn
	Lock *sync.Mutex
	Send chan ctypes.ServerEvent
	LastMove time.Time
	ParticipantID string
}
