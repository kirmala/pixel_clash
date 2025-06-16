package ctypes


type ServerEvent struct {
	Type string      `json:"type"` // "game_found", "opponent_move"
	Data interface{} `json:"data"` // No ID needed (not a response)
}

type PlayerMove struct {
	Feild Feild `json:"feild"`
	Scores Scores `json:"scores"`
}