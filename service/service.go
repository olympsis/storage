package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/v2/apiv1"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func NewStorageService(l *logrus.Logger) *Service {
	return &Service{Logger: l}
}

func (s *Service) ConnectToClient() error {
	filePath := os.Getenv("FIREBASE_CREDENTIALS_PATH")

	// Storage Client
	client, err := storage.NewClient(context.TODO(), option.WithCredentialsFile(filePath))
	if err != nil {
		s.Logger.Fatal("Failed to create client: " + err.Error())
		return err
	}
	s.Client = client

	// Computer Vision Client
	vClient, err := vision.NewImageAnnotatorClient(context.TODO(), option.WithCredentialsFile(filePath))
	if err != nil {
		s.Logger.Fatal("Failed to create client: " + err.Error())
		return err
	}
	s.VClient = vClient

	return nil
}

func (s *Service) CORSHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Filename")
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (s *Service) UploadObject() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fileBucket := r.PathValue("fileBucket")
		if len(fileBucket) < 1 {
			http.Error(rw, "invalid file bucket name", http.StatusBadRequest)
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
			s.Logger.Error(err.Error())
			http.Error(rw, "failed to read body", http.StatusBadRequest)
			return
		}

		// Get the filename from the request header.
		fileName, err := GrabFileName(&r.Header)
		if err != nil {
			s.Logger.Error(err.Error())
			http.Error(rw, "no file name", http.StatusBadRequest)
			return
		}

		resp, err := s.AnnotateImage(bodyData)
		if err != nil {
			s.Logger.Error(err.Error())
			http.Error(rw, "failed to validate image", http.StatusBadRequest)
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
			s.Logger.Error(err.Error())
			http.Error(rw, "failed to upload image", http.StatusInternalServerError)
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
		fileBucket := r.PathValue("fileBucket")

		// Get the filename from the request header.
		fileName, err := GrabFileName(&r.Header)
		if err != nil {
			s.Logger.Error(err.Error())
			http.Error(rw, "no file name", http.StatusBadRequest)
			return
		}

		if len(fileBucket) < 1 {
			http.Error(rw, "invalid file bucket name", http.StatusBadRequest)
			return
		}

		err = s._deleteObject(fileBucket, fileName)
		if err != nil {
			s.Logger.Error(err.Error())
			http.Error(rw, "failed to delete image", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}
}
