package model

import (
	"pixel_clash/ctypes"
	"time"
)

type Game struct {
	Id string
	PlayerIds map[string]struct{}
	Type ctypes.Game 
	Feild   ctypes.Feild
	Status  string
	Timer time.Timer
}