package service

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

func (s *Service) _uploadObject(file []byte, bucket string, name string) error {

	// object metadata
	metadata := minio.PutObjectOptions{
		ContentType: getContentType(name),
	}

	// Upload the file to the bucket using PutObject
	_, err := s.Client.PutObject(context.Background(), bucket, name, bytes.NewReader(file), int64(len(file)), metadata)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) _deleteObject(bucket string, name string) error {
	// Delete object from bucket
	err := s.Client.RemoveObject(context.Background(), bucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GrabFileName(h *http.Header) (string, error) {
	// Get the filename from the request header.
	fileName := h.Get("X-Filename")
	if fileName == "" {
		return "", errors.New("missing X-Filename header")
	} else {
		return fileName, nil
	}
}

func getContentType(fileName string) string {
	extension := filepath.Ext(fileName)
	switch extension {
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".heic":
		return "image/heic"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".csv":
		return "text/csv"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
