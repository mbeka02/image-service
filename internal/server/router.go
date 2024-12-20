package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Post("/register", s.UserHandler.handleCreateUser)
	r.Post("/login", s.UserHandler.handleLogin)

	r.Route("/images", func(r chi.Router) {
		r.Use(AuthMiddleware(s.AuthMaker))
		r.Get("/", s.ImageHandler.handleGetImages)
		r.Post("/", s.ImageHandler.handleImageUpload)
		r.Get("/{imageId}", s.ImageHandler.handleGetImage)
		r.Post("/{imageId}/transform", s.ImageHandler.handleImageTransformations)
		r.Delete("/{imageId}/delete", s.ImageHandler.handleDeleteImage)
	})

	return r
}
