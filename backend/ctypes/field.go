package ctypes

type Cell struct {
	CompSize int    `json:"component_size"`
	ParticipantID  string `json:"participant_id"`
}

type Field struct {
	Data [][]Cell `json:"data"`
}
