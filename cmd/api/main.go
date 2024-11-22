package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbeka02/image-service/internal/auth"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/imgstore"
	"github.com/mbeka02/image-service/internal/mailer"
	"github.com/mbeka02/image-service/internal/server"

	"github.com/mbeka02/image-service/config"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("...shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown with error: %v\n", err)
	}

	log.Println("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
	store, err := database.NewStore(conf.DB_URI)
	if err != nil {
		log.Fatalf("...unable to setup the db : %v", err)
	}
	maker, err := auth.NewJWTMaker(conf.SYMMETRIC_KEY)
	if err != nil {
		log.Fatalf("...unable to setup up the auth token maker:%v", err)
	}
	newMailer := mailer.NewMailer(conf.MAILER_HOST, conf.MAILER_PASSWORD)
	fileStorage, err := imgstore.NewGCStorage(conf.GCLOUD_PROJECT_ID, conf.GCLOUD_BUCKET_NAME)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool, 1)
	server := server.NewServer(":"+conf.PORT, store, maker, conf.ACCESS_TOKEN_DURATION, newMailer, fileStorage)
	go gracefulShutdown(server, done)
	log.Println("the server is listening on port:" + conf.PORT)
	server.ListenAndServe()
}
