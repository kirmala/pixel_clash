package types

import (
	"pixel_clash/ctypes"
)

type FindGameRequest struct {
	Nickname string `json:"nickname"`
	GameType ctypes.Game  `json:"game_type"`
}

type FindGameResponse struct {
	GameId string `json:"game_id"`
	PlayerId string `json:"player_id"`
}

type MoveRequest struct {
	PlayerId string `json:"player_id"`
	X int `json:"x"`
	Y int `json:"y"`
}

type MoveResponse struct {
	Message string `json:"message"`
}

type StopSearchingRequest struct {
	PlayerId string `json:"player_id"`
}

type StopSearchingResponse struct {
	Message string `json:"message"`
}
