package service

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func NewStorageService(l *logrus.Logger, r *mux.Router) *Service {
	return &Service{Logger: l, Router: r}
}

func (s *Service) ConnectToClient() {
	endpoint := os.Getenv("STORAGE_ADDR")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL})
	if err != nil {
		s.Logger.Fatalln(err)
	} else {
		s.Client = minioClient
		s.Logger.Info("connected to minio client...")
	}
}

func (s *Service) UploadObject() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileBucket := vars["fileBucket"]
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

		// Upload file to MinIO server
		err = s._uploadObject(bodyData, fileBucket, fileName)
		if err != nil {
			s.Logger.Error(err.Error())
			http.Error(rw, "failed to upload image", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}
}

func (s *Service) DeleteObject() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileBucket := vars["fileBucket"]

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
