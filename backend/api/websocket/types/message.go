package types

import (
	"log"

	"github.com/gorilla/websocket"
)

type errorMsg struct {
	Err string `json:"error"`
}

func ProcessError(conn *websocket.Conn, resp any, err error) {
	if err != nil {
		wErr := conn.WriteJSON(errorMsg{Err : err.Error()})
		if wErr != nil {
			log.Printf("%s\n", wErr.Error())
		}
		return
	}
	wErr := conn.WriteJSON(errorMsg{Err : ""})
	if wErr != nil {
		log.Printf("%s\n", wErr.Error())
	}
	if resp != nil {
		if err := conn.WriteJSON(resp); err != nil {
			log.Println(err.Error())
		}
	}
}
