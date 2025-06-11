package model

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	Id           string
	Nickname     string
	GameCapacity int
	GameId       string
	Status       string
	Color        int
	Connection   *websocket.Conn
}