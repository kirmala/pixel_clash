package main

import (
	"log"
	"pixel_clash/api/websocket"
	pkgHttp "pixel_clash/pkg/http"
	"pixel_clash/repository/short/ram"
	"pixel_clash/usecase/service"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "pixel_clash/docs"

	"github.com/go-chi/chi"
)

// @title pixel_clash
// @version 1.0
// @description This is pixel_clash.

// @host localhost:8080
// @BasePath /
func main() {
	addr := "0.0.0.0:8080"

	gameRepo := sram.NewGame()
	playerRepo := sram.NewPlayer()

	gameService := service.NewGame(gameRepo, playerRepo)
	playerService := service.NewPlayer(playerRepo)

	joinHandler := websocket.NewGameWebsocketHandler(playerService, gameService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))
	joinHandler.WithGameHandlers(r)

	log.Printf("Starting server on %s", addr)
	if err := pkgHttp.CreateAndRunServer(r, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}