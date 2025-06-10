package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PostPlayerJoinHandlerRequest struct {
	Nickname string `json:"nickname"`
	Capacity int `json:"capacity"`
}

type PostPlayerJoinHandlerResponse struct {
	Status string `json:"status"`
	GameId string `json:"game_id"`
	PlayerId string `json:"player_id"`
}

func CreatePostPlayerJoinHandlerRequest(r *http.Request) (*PostPlayerJoinHandlerRequest, error) {
	var req PostPlayerJoinHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type PostPlayerStatusHandlerRequest struct {
	Id string `json:"id"`
}

type PostPlayerStatusHandlerResponse struct {
	Status string `json:"status"`
}

func CreatePostPlayerStatusHandlerRequest(r *http.Request) (*PostPlayerStatusHandlerRequest, error) {
	var req PostPlayerStatusHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}