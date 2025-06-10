package model

import (
	"time"
)

type Game struct {
	Players []Player
	Capacity int
	Id string
	Feild   [][]int
	Status  string
	Timer time.Timer
	Winner string
}