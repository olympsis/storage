package service

import (
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/v2/apiv1"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Client      *storage.Client
	VClient     *vision.ImageAnnotatorClient
	Logger      *logrus.Logger
	AccessToken *string
}

type Response struct {
	// image url
	URL *string `json:"url,omitempty"`

	// image score
	Score int `json:"score"`

	// score reason
	Reason *string `json:"reason,omitempty"`
}
