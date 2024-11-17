package server

import (
	"net/http"
	"time"

	"github.com/mbeka02/image-service/internal/auth"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/mailer"
)

// type ServerOpts struct {}
type Server struct {
	Addr                string
	Store               *database.Store
	AuthMaker           auth.Maker
	Mailer              *mailer.Mailer
	AccessTokenDuration time.Duration
}

func NewServer(addr string, store *database.Store, maker auth.Maker, duration time.Duration, mailer *mailer.Mailer) *http.Server {
	srv := Server{
		Addr:                addr,
		Store:               store,
		AuthMaker:           maker,
		AccessTokenDuration: duration,
		Mailer:              mailer,
	}

	return &http.Server{
		Handler: srv.RegisterRoutes(),
		Addr:    srv.Addr,
	}
}
