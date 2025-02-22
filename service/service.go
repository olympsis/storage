package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/v2/apiv1"
	"github.com/gorilla/mux"
	"github.com/olympsis/storage/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func NewStorageService(l *logrus.Logger) *Service {
	return &Service{Logger: l}
}

func (s *Service) ConnectToClient(config *utils.ServerConfig) error {
	// Storage Client
	client, err := storage.NewClient(context.TODO(), option.WithCredentialsFile(config.FirebaseFilePath))
	if err != nil {
		s.Logger.Fatal("Failed to create client: " + err.Error())
		return err
	}
	s.Client = client

	// Computer Vision Client
	vClient, err := vision.NewImageAnnotatorClient(context.TODO(), option.WithCredentialsFile(config.FirebaseFilePath))
	if err != nil {
		s.Logger.Fatal("Failed to create client: " + err.Error())
		return err
	}
	s.VClient = vClient
	s.Logger.Info("Connected to GCP Client Successfully!")

	return nil
}

func (s *Service) UploadObject() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileBucket := vars["fileBucket"]
		if len(fileBucket) == 0 {
			http.Error(rw, `{ "msg" : "invalid file bucket name" }`, http.StatusBadRequest)
			return
		}

		// Print the request headers
		fmt.Println("-----------------------------------------")
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		fmt.Println("-----------------------------------------")

		// read body data
		bodyData, err := io.ReadAll(r.Body)
		if err != nil {
			s.Logger.Error("Failed to read request body. Error: ", err.Error())
			http.Error(rw, `{ "msg" : "failed to read body" }`, http.StatusBadRequest)
			return
		}

		// Get the filename from the request header.
		fileName, err := GrabFileName(&r.Header)
		if err != nil {
			s.Logger.Error("No filename found in header")
			http.Error(rw, `{ "msg" : "no file name in header" }`, http.StatusBadRequest)
			return
		}

		resp, err := s.AnnotateImage(bodyData)
		if err != nil {
			s.Logger.Error("Failed to validate image. Error: ", err.Error())
			http.Error(rw, `{ "msg" : "failed to validate image" }`, http.StatusBadRequest)
			return
		}

		var response Response

		// Safety Score
		safety := s.AggregateSafetyScore(resp.Responses[0])
		response.Score = *safety

		if *safety == 5 {
			reason := "Unsafe image"
			response.Reason = &reason

			rw.WriteHeader(http.StatusBadRequest)
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(response)
			return
		}

		// Upload file to Storage server
		err = s._uploadObject(bodyData, fileBucket, fileName)
		if err != nil {
			s.Logger.Error("Failed to upload file. Error: ", err.Error())
			http.Error(rw, `{ "msg" : "failed to upload image" }`, http.StatusBadRequest)
			return
		}

		url := fileBucket + "/" + fileName
		response.URL = &url

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(response)
	}
}

func (s *Service) DeleteObject() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileBucket := vars["fileBucket"]
		if len(fileBucket) == 0 {
			http.Error(rw, `{ "msg" : "invalid file bucket name" }`, http.StatusBadRequest)
			return
		}

		// Get the filename from the request header.
		fileName, err := GrabFileName(&r.Header)
		if err != nil {
			s.Logger.Error("Failed to delete file. No file name in header")
			http.Error(rw, `{ "msg" : "no file name in header" }`, http.StatusBadRequest)
			return
		}

		err = s._deleteObject(fileBucket, fileName)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to delete image (%s). Error: (%s)", fileName, err.Error()))
			http.Error(rw, `{ "msg" : "failed to delete image" }`, http.StatusBadRequest)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}
}
