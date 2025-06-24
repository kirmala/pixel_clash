package ctypes

type ServerMessage struct {
	Type string `json:"type"` // "event", "response"
	Data interface{} `json:"data"`
}

type ServerEvent struct {
	Type string      `json:"type"` // "game_start", game_finish", "opponent_move", "waiting_change"
	Data interface{} `json:"data"`
}

type GameFinish struct {
	Scores Scores `json:"scores"`
}

type WaitingChange struct {
	Waiting int `json:"waiting"`
}

type PlayerMove struct {
	Field  Field  `json:"field"`
	Scores Scores `json:"scores"`
}
