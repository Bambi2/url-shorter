package service

import (
	"errors"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/bambi2/url-shorter/internal/repository"
)

const ()

type EncoderService struct {
	repo repository.Encoder
}

func NewEncoderService(repo repository.Encoder) *EncoderService {
	return &EncoderService{repo: repo}
}

func (s *EncoderService) Base63(urlString string) (string, error) {
	if _, err := url.ParseRequestURI(urlString); err != nil {
		return "", &ServiceError{Msg: "invalid url: " + urlString, StatusCode: http.StatusBadRequest}
	}

	id, err := s.repo.IfExistsBase63(urlString)
	if err != nil {
		return "", &ServiceError{Msg: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if id != -1 {
		return base63Encode(id), nil
	}

	for {
		id = rand.Int63n(Base63TenMaxId + 1)
		err := s.repo.SaveBase63(urlString, id)
		// might need wrapper for logging
		if err != nil {
			if !errors.Is(err, repository.ErrDuplicateId) {
				return "", &ServiceError{Msg: err.Error(), StatusCode: http.StatusInternalServerError}
			}
		} else {
			break
		}
	}

	return base63Encode(id), nil
}
