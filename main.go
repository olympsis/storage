package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olympsis/storage/service"
	"github.com/sirupsen/logrus"
)

func main() {

	// Create logger
	l := logrus.New()

	// Environment variables
	port := os.Getenv("PORT")

	// Create Service
	// Connect to Storage & Computer Vision clients
	n := service.NewStorageService(l)
	err := n.ConnectToClient()
	if err != nil {
		panic(err.Error())
	}

	// Route handler
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/storage/{fileBucket}", n.UploadObject())
	mux.HandleFunc("DELETE /v1/storage/{fileBucket}", n.DeleteObject())

	// Server
	s := &http.Server{
		Addr:         `:` + port,
		Handler:      mux,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Start server
	go func() {
		l.Info(`Starting storage service at...` + port)
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
