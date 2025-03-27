package router

import (
	"github.com/go-chi/chi/v5"
	"shortcut-challenge/handlers"
	"shortcut-challenge/middleware"
	"shortcut-challenge/models"
)

func getInventoryRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.AuthMiddleware)

	router.Group(func(r chi.Router) {
		r.Use(middleware.RestrictTo(string(models.ADMIN)))

		r.Post("/", handlers.CreateItem)
		r.Post("/{itemID}/restock", handlers.RestockItem)
	})

	router.Group(func(r chi.Router) {
		r.Get("/restock", handlers.GetRestockHistory)
		r.Get("/", handlers.GetInventoryItems)
	})

	return router
}
