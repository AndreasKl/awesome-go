package api

import (
	"github.com/go-chi/chi/v5"

	"awesome/login/frontend"
)

func Mount(router chi.Router) error {
	controller, err := frontend.NewLoginController()
	if err != nil {
		return err
	}
	router.Get("/login", controller.RenderForm)
	router.Post("/login", controller.HandleForm)

	return nil
}
