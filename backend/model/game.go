package model

import (
	"pixel_clash/ctypes"
	"time"
)

type Game struct {
	Id      string
	Players map[Player]struct{}
	Type    ctypes.Game
	Field   ctypes.Field
	Status  string
	Scores  ctypes.Scores
	Timer   time.Timer
}
