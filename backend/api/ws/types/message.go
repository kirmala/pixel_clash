package types

import "encoding/json"

type Request struct {
	Type string          `json:"type"` // "move", "join", etc.
	ID   string          `json:"id"`   // Unique request ID (UUID or timestamp)
	Data json.RawMessage `json:"data"`
}

// Server -> Client (Response to a specific request)
type ServerResponse struct {
	Type   string      `json:"type"`   // "move_result", "join_result"
	ID     string      `json:"id"`     // Echoes the client's request ID
	Status string      `json:"status"` // "success", "error"
	Data   interface{} `json:"data"`   // Response payload
}


// Server -> Client (Unsolicited event)

