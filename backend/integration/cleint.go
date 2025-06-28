package testing

import (
	"encoding/json"
	"fmt"
	"pixel_clash/api/ws/types"
	"pixel_clash/ctypes"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Cli struct {
	Nickname string
	GamePref ctypes.Game
	Conn     *websocket.Conn
	LastReq  *string
}

func NewCli(nickname string, gamePref ctypes.Game, URL string) (*Cli, error) {
	conn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		return nil, err
	}
	return &Cli{Nickname: nickname, GamePref: gamePref, Conn: conn}, nil
}

// func (client *Cli) FindGame() error {
// 	req := types.FindGameRequest{Nickname: client.Nickname, GameType: client.GamePref}
// 	jReq, err := json.Marshal(req)
// 	if err != nil {
// 		return fmt.Errorf("marshalling find game request: %s", err)
// 	}
// 	if err := client.Conn.WriteJSON(jReq); err != nil {
// 		return fmt.Errorf("sending find game request: %s", err)
// 	}
// 	return nil
// }

// func (client *Cli) Move(x, y int) error {
// 	req := types.MoveRequest{X: x, Y: y}
// 	jReq, err := json.Marshal(req)
// 	if err != nil {
// 		return fmt.Errorf("marshalling move request: %s", err)
// 	}
// 	if err := client.Conn.WriteJSON(jReq); err != nil {
// 		return fmt.Errorf("sending move request: %s", err)
// 	}
// 	return nil
// }

// func (client *Cli) StopSearching(message string) error {
// 	req := types.StopSearchingRequest{Message: message}
// 	jReq, err := json.Marshal(req)
// 	if err != nil {
// 		return fmt.Errorf("marshalling stop searching request: %s", err)
// 	}
// 	if err := client.Conn.WriteJSON(jReq); err != nil {
// 		return fmt.Errorf("sending stop searching request: %s", err)
// 	}
// 	return nil
// }

func (client *Cli) Send(data interface{}, reqType string) (*string, error) {
	jData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshalling %s request: %s", reqType, err)
	}
	req := types.Request{Data: jData, Type: reqType, ID: uuid.NewString()}

	if err := client.Conn.WriteJSON(req); err != nil {
		return nil, fmt.Errorf("sending %s request: %s", req, err)
	}
	return &req.ID, nil
}

func (client *Cli) Receive() (*ctypes.ServerMessage, error) {
	var message ctypes.ServerMessage
	if err := client.Conn.ReadJSON(&message); err != nil {
		return nil, fmt.Errorf("receiving: %s", err)
	}
	return &message, nil
}

func (client *Cli) GetEvent() (*ctypes.ServerEvent, error) {
	var event ctypes.ServerEvent
	if err := client.Conn.ReadJSON(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (client *Cli) GetResponse() (*types.ServerResponse, error) {
	var resp types.ServerResponse
	if err := client.Conn.ReadJSON(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
