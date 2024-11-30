package api

import (
	"net/http"
	"time"

	"github.com/mbeka02/image-service/internal/auth"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/imgproc"
	"github.com/mbeka02/image-service/internal/imgstore"
	"github.com/mbeka02/image-service/internal/mailer"
)

type Server struct {
	Addr                string
	Store               *database.Store
	AuthMaker           auth.Maker
	Mailer              *mailer.Mailer
	FileStorage         imgstore.Storage
	ImageProcessor      imgproc.ImageProcessor
	AccessTokenDuration time.Duration
}

func NewServer(addr string, store *database.Store, maker auth.Maker, duration time.Duration, mailer *mailer.Mailer, fileStorage imgstore.Storage, imageProcessor imgproc.ImageProcessor) *http.Server {
	srv := Server{
		Addr:                addr,
		Store:               store,
		AuthMaker:           maker,
		AccessTokenDuration: duration,
		Mailer:              mailer,
		FileStorage:         fileStorage,
		ImageProcessor:      imageProcessor,
	}

	return &http.Server{
		Handler: srv.RegisterRoutes(),
		Addr:    srv.Addr,
	}
}
