package server

import (
	"fmt"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test route")
	})

	return &http.Server{
		Handler: mux,
		Addr:    srv.Addr,
	}
}
