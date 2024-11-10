package server

import (
	"net/http"

	"github.com/mbeka02/image-service/internal/database"
)

type Server struct {
	Addr  string
	Store *database.Store
}

func NewServer(addr string, store *database.Store) *http.Server {
	srv := Server{
		Addr:  addr,
		Store: store,
	}

	return &http.Server{
		Handler: srv.RegisterRoutes(),
		Addr:    srv.Addr,
	}
}
