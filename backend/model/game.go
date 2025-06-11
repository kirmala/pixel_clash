package model

import (
	"time"
)

type Game struct {
	Id string
	Players []Player
	Capacity int
	Feild   [][]int
	ThreasholdSquare int
	Timer time.Timer

	Status  string
	
	Winner string
}