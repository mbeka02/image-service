package server

import (
	"net/http"
	"time"

	"github.com/mbeka02/image-service/internal/auth"
	"github.com/mbeka02/image-service/internal/database"
)

type Server struct {
	Addr                string
	Store               *database.Store
	AuthMaker           auth.Maker
	AccessTokenDuration time.Duration
}

func NewServer(addr string, store *database.Store, maker auth.Maker, duration time.Duration) *http.Server {
	srv := Server{
		Addr:                addr,
		Store:               store,
		AuthMaker:           maker,
		AccessTokenDuration: duration,
	}

	return &http.Server{
		Handler: srv.RegisterRoutes(),
		Addr:    srv.Addr,
	}
}
