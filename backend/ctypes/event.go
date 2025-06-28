package ctypes

type ServerMessage struct {
	Type string `json:"type"` // "event", "response"
	Data interface{} `json:"data"`
}

type ServerEvent struct {
	Type string      `json:"type"` // "game_start", game_finish", "opponent_move", "waiting_change"
	Data interface{} `json:"data"`
}

type GameStart struct {
	Field  Field  `json:"field"`
	Participants Participants `json:"participants"`
}

type GameFinish struct {
	Participants Participants `json:"participants"`
}

type WaitingChange struct {
	Waiting int `json:"waiting"`
}

type PlayerMove struct {
	Field  Field  `json:"field"`
	Participants Participants `json:"participants"`
}
