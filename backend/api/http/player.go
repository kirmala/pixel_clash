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

// @Summary user joins game
// @Description user joins game
// @Tags user
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

	var gameId string
	gameId, player.Status = u.gameService.Find(player)

	types.ProcessError(w, err, &types.PostPlayerJoinHandlerResponse{Status : player.Status, GameId : gameId}, 200)
}

func (u *Player) WithPlayerHandlers(r chi.Router) {
	r.Route("/player", func(r chi.Router) {
		r.Post("/join", u.joinHandler)
	})
}
