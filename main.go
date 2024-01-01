package main

import (
	"context"
	"net/http"
	"olympsis-storage/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	// logger
	l := logrus.New()

	// mux router
	r := mux.NewRouter()

	storage := service.NewStorageService(l, r)
	storage.ConnectToClient()

	r.Handle("/storage/{fileBucket}", storage.UploadObject()).Methods("POST")
	r.Handle("/storage/{fileBucket}", storage.DeleteObject()).Methods("DELETE")

	port := os.Getenv("PORT")

	// server config
	s := &http.Server{
		Addr:         `:` + port,
		Handler:      r,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	// start server
	go func() {
		l.Info(`starting olympsis storage server at...` + port)
		err := s.ListenAndServe()

		if err != nil {
			l.Info("error starting server: ", err)
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs

	l.Printf("Recieved Termination(%s), graceful shutdown \n", sig)

	tc, c := context.WithTimeout(context.Background(), 30*time.Second)

	defer c()

	s.Shutdown(tc)
}
