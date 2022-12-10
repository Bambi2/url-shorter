package service

import (
	"errors"
	"net/http"

	"github.com/bambi2/url-shorter/internal/repository"
)

type DecoderService struct {
	repo repository.Decoder
}

func NewDecoderService(repo repository.Decoder) *DecoderService {
	return &DecoderService{repo: repo}
}

func (s *DecoderService) Base63(encodedUrl string) (string, error) {
	if err := urlLengthCheck(encodedUrl, base63EncodedLength); err != nil {
		return "", err
	}

	id, err := base63Decode(encodedUrl)
	if err != nil {
		return "", err
	}

	url, err := s.repo.GetBase63(id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchURL) {
			return "", &ServiceError{Msg: "no such url", StatusCode: http.StatusBadRequest}
		}
		return "", &ServiceError{Msg: err.Error(), StatusCode: http.StatusInternalServerError}
	}

	return url, nil
}
