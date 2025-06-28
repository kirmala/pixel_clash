package types

import (
	"pixel_clash/ctypes"
)

type FindGameRequest struct {
	Nickname string `json:"nickname"`
	GameType ctypes.Game  `json:"game_type"`
}

type FindGameResponse struct {
	ParticipantID string `json:"participant_id"`
}

type MoveRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MoveResponse struct {
	Message string `json:"message"`
}

type StopSearchingRequest struct {
	Message string `json:"message"`
}

type StopSearchingResponse struct {
	Message string `json:"message"`
}
