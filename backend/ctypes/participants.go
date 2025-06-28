package ctypes

type Participant struct {
	ID string `json:"data"`
	Score int `json:"score"`
	Nickname string `json:"nickname"`
}

type Participants struct {
	Data map[string]Participant `json:"data"`
}