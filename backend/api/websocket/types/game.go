package types

import (
	"pixel_clash/ctypes"
	"pixel_clash/model"
)

type JoinHandlerRequest struct {
	Nickname string `json:"nickname"`
	GameType ctypes.Game  `json:"game_type"`
}

type JoinHandlerResponse struct {
	GameId string `json:"game_id"`
	PlayerId string `json:"player_id"`
}

type MoveRequest struct {
	PlayerId string `json:"player_id"`
	X int `json:"x"`
	Y int `json:"y"`
}

type MoveResponse struct {
	Game model.Game `json:"game"`
}
