package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func CreateAndRunServer(r chi.Router, addr string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}