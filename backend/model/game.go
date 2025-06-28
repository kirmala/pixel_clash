package model

import (
	"pixel_clash/ctypes"
	"time"
)

type Game struct {
	ID      string
	Players map[Player]struct{}
	Type    ctypes.Game
	Field   ctypes.Field
	Status  string
	Participants  ctypes.Participants
	Timer   time.Timer
}
