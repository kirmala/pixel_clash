package types

import "encoding/json"

type Request struct {
	Type string          `json:"type"` // "move", "find_game", "stop_searching",
	ID   string          `json:"id"`   // Unique request ID (UUID or timestamp)
	Data json.RawMessage `json:"data"`
}

type ServerResponse struct {
	Type   string      `json:"type"`   // "move_result", "find_game_result", "stop_searching_result"
	ID     string      `json:"id"`     // Echoes the client's request ID
	Status string      `json:"status"` // "success", "error"
	Data   interface{} `json:"data"`   // Response payload
}

