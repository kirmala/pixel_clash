package types

import (
	"log"

	"github.com/gorilla/websocket"
)

type StatusMessage struct {
	Status string `json:"status"`
}

func SendError(conn *websocket.Conn, err error) {
	wErr := conn.WriteJSON(StatusMessage{Status : err.Error()})
	if wErr != nil {
		log.Printf("%s\n", wErr.Error())
	}
}

func SendResponse(conn *websocket.Conn, resp any) {
	wErr := conn.WriteJSON(StatusMessage{Status : "OK"})
	if wErr != nil {
		log.Printf("%s\n", wErr.Error())
	}
	if resp != nil {
		if err := conn.WriteJSON(resp); err != nil {
			log.Println(err.Error())
		}
	}
}
