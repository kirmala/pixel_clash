package main

import (
	"log"
	"pixel_clash/api/http"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	pkgHttp "pixel_clash/pkg/http"
	"pixel_clash/repository/ram"
	"pixel_clash/usecase/service"

	"github.com/go-chi/chi"
	_ "pixel_clash/docs"
)


// @title pixel_clash
// @version 1.0
// @description This is pixel_clash.

// @host localhost:8080
// @BasePath /
func main() {
	addr := "0.0.0.0:8080"

	gameRepo := ram.NewGame()
	playerRepo := ram.NewPlayer()

	gameService := service.NewGame(gameRepo, playerRepo)
	playerService := service.NewPlayer(playerRepo)

	playerHandlers := http.NewPlayerHandler(playerService, gameService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))
	playerHandlers.WithPlayerHandlers(r)

	log.Printf("Starting server on %s", addr)
	if err := pkgHttp.CreateAndRunServer(r, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}