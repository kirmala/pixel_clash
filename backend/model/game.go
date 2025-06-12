package model

import (
	"pixel_clash/ctypes"
)

type Game struct {
	Id string
	PlayerIds []string
	Type ctypes.Game 
	Feild   [][]ctypes.Cell
	Status  string
}