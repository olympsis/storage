package service

import (
	"context"
	"errors"
	"net/http"

	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

func (s *Service) _uploadObject(file []byte, bucket string, name string) error {

	// upload file to bucket
	object := s.Client.Bucket(bucket).Object(name)
	wc := object.NewWriter(context.Background())
	_, err := wc.Write(file)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) _deleteObject(bucket string, name string) error {

	object := s.Client.Bucket(bucket).Object(name)
	err := object.Delete(context.Background())
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

func (s *Service) AnnotateImage(image []byte) (*visionpb.BatchAnnotateImagesResponse, error) {

	// convert image bytes to annotation request
	ann := visionpb.AnnotateImageRequest{
		Image: &visionpb.Image{
			Content: image,
		},
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_SAFE_SEARCH_DETECTION,
			},
		},
	}

	// create batch request call
	req := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{&ann},
	}

	// create request and return response
	resp, err := s.VClient.BatchAnnotateImages(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) AggregateSafetyScore(response *visionpb.AnnotateImageResponse) *int {
	annotation := response.SafeSearchAnnotation
	score := 0

	// Adult/Racy content check
	if annotation.Adult == visionpb.Likelihood_VERY_LIKELY {
		score = 5
		return &score
	} else if annotation.Racy == visionpb.Likelihood_VERY_LIKELY {
		score = 5
		return &score
	} else if annotation.Adult > visionpb.Likelihood_POSSIBLE && annotation.Racy > visionpb.Likelihood_POSSIBLE {
		score = 4
	} else if annotation.Adult <= visionpb.Likelihood_POSSIBLE && annotation.Racy <= visionpb.Likelihood_POSSIBLE {
		score = int(annotation.Racy.Number())
	}

	// Violent content check
	if annotation.Violence == visionpb.Likelihood_VERY_LIKELY {
		score = 5
		return &score
	} else if annotation.Violence <= visionpb.Likelihood_LIKELY {
		score += int(annotation.Violence.Number())
	}

	// Medical content check
	if annotation.Medical == visionpb.Likelihood_VERY_LIKELY {
		score = 5
		return &score
	} else if annotation.Medical <= visionpb.Likelihood_LIKELY {
		score += int(annotation.Medical.Number())
	}

	score /= 3

	return &score
}
