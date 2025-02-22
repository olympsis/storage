package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/olympsis/storage/middleware"
	"github.com/olympsis/storage/service"
	"github.com/olympsis/storage/utils"
	"github.com/sirupsen/logrus"
)

func main() {

	// Create logger
	l := logrus.New()

	// Environment variables
	config := utils.GetServerConfig()

	// Create Service
	// Connect to Storage & Computer Vision clients
	n := service.NewStorageService(l)
	err := n.ConnectToClient(config)
	if err != nil {
		panic(err.Error())
	}

	// Route handler
	mux := mux.NewRouter()
	mux.Handle("/v1/storage/{fileBucket}",
		middleware.Chain(
			n.UploadObject(),
			middleware.CORS(),
		),
	).Methods("POST", "OPTIONS")

	mux.Handle("/v1/storage/{fileBucket}",
		middleware.Chain(
			n.DeleteObject(),
			middleware.CORS(),
		),
	).Methods("DELETE", "OPTIONS")

	// Server
	s := &http.Server{
		Addr:         `:` + config.Port,
		Handler:      mux,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Start server
	go func() {
		l.Info(`Starting storage service at...` + config.Port)
		err := s.ListenAndServe()

		if err != nil {
			l.Info("Error starting server: ", err)
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs

	l.Printf("Received Termination(%s), graceful shutdown \n", sig)

	tc, c := context.WithTimeout(context.Background(), 30*time.Second)

	defer c()

	s.Shutdown(tc)
}
