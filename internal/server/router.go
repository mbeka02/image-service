package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.Post("/register", s.UserHandler.handleCreateUser)
	r.Post("/login", s.UserHandler.handleLogin)

	r.Route("/images", func(r chi.Router) {
		r.Use(AuthMiddleware(s.AuthMaker))
		r.Get("/get", s.ImageHandler.handleGetImages)
		r.Post("/upload", s.ImageHandler.handleImageUpload)
		r.Post("/resize", s.ImageHandler.handleImageResize)
		r.Post("/rotate", s.ImageHandler.handleImageRotation)
		r.Post("/crop", s.ImageHandler.handleImageCropping)
		r.Post("/flip", s.ImageHandler.handleImageFlip)
		r.Post("/convert", s.ImageHandler.handleImageConversion)

		r.Delete("/delete/{imageId}", s.ImageHandler.handleDeleteImage)
	})

	return r
}
