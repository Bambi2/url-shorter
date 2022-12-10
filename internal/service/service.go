package service

import "github.com/bambi2/url-shorter/internal/repository"

type Encoder interface {
	Base63(urlString string) (string, error)
}

type Decoder interface {
	Base63(encodedUrl string) (string, error)
}

type Service struct {
	Encoder
	Decoder
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Encoder: NewEncoderService(repo.Encoder),
		Decoder: NewDecoderService(repo.Decoder),
	}
}
