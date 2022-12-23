package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"awesome/smurfs/memory"
	"awesome/smurfs/resource"
)

func Mount(router chi.Router) error {
	controller := resource.NewController(&memory.Repository{})

	router.Get("/smurfs", withTimeOut(controller.List))
	router.Get("/smurfs/{id}", withTimeOut(controller.Get))

	return nil
}

func withTimeOut(h http.HandlerFunc) http.HandlerFunc {
	return http.TimeoutHandler(h, 2*time.Second, "timeout").ServeHTTP
}
