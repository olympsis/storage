package test

import (
	"olympsis-storage/service"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func TestCreateClient(t *testing.T) {
	// logger
	l := logrus.New()

	// mux router
	r := mux.NewRouter()

	storage := service.NewStorageService(l, r)
	storage.ConnectToClient()
}
