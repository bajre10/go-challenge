package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"shortcut-challenge/database"
	_ "shortcut-challenge/docs"
)

func SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(database.SetDBMiddleware)
	router.Use(middleware.StripSlashes)
	router.Route("/api", func(r chi.Router) {
		r.Mount("/auth", getAuthRoutes())
		r.Mount("/inventory", getInventoryRoutes())
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return router
}
