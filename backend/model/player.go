package model

import (
	"pixel_clash/ctypes"
)

type Player struct {
	Id         string
	UserId     string
	Nickname   string
	GameType   ctypes.Game
	GameId     string
	Send chan ctypes.ServerEvent
}
