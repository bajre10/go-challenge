package router

import (
	"github.com/go-chi/chi/v5"
	"shortcut-challenge/handlers"
	"shortcut-challenge/middleware"
	"shortcut-challenge/models"
)

func getAuthRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/register", handlers.Register)
	router.Post("/login", handlers.Login)

	router.
		With(middleware.AuthMiddleware).
		With(middleware.RestrictTo(string(models.ADMIN))).
		Post("/admin/register", handlers.RegisterAdmin)

	return router
}
