package http

import (
	"net/http"
	"pixel_clash/api/http/types"
	"pixel_clash/model"
	"pixel_clash/usecase"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Player struct {
	service usecase.Player
	gameService usecase.Game
}


func NewPlayerHandler(playerService usecase.Player, gameService usecase.Game) *Player {
	return &Player{service: playerService, gameService: gameService}
}

// @Summary player joins game
// @Description player joins game
// @Tags player
// @Accept  json
// @Param name body types.PostPlayerJoinHandlerRequest true "user name, desired capacity"
// @Router /player/join [post]
func (u *Player) joinHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostPlayerJoinHandlerRequest(r)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	player := model.Player{Id : uuid.NewString(), Nickname: req.Nickname, Status: "searching", GameCapacity: req.Capacity}
	
	err = u.service.Post(player)

	u.gameService.Find(player)

	updatedPlayer, _ := u.service.Get(player.Id)

	types.ProcessError(w, err, &types.PostPlayerJoinHandlerResponse{Status : updatedPlayer.Status, GameId : updatedPlayer.GameId, PlayerId: updatedPlayer.Id}, 200)
}

// @Summary sends status info
// @Description sends status info
// @Tags player
// @Accept  json
// @Param name body types.PostPlayerStatusHandlerRequest true "user id"
// @Router /player/status [post]
func (u *Player) statusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostPlayerStatusHandlerRequest(r)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	player, err := u.service.Get(req.Id)

	var status string

	if err != nil {
		status = "unauthorized"
	} else {
		status = player.Status
	}

	types.ProcessError(w, err, &types.PostPlayerStatusHandlerResponse{Status : status}, 200)
}

func (u *Player) WithPlayerHandlers(r chi.Router) {
	r.Route("/player", func(r chi.Router) {
		r.Post("/join", u.joinHandler)
		r.Post("/status", u.statusHandler)
	})
}
