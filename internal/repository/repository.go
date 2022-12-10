package repository

import (
	"errors"

	"github.com/bambi2/url-shorter/internal/database"
)

// go:generate mockgen -source=repository.go -destination=mocks/mock.go
const (
	TRIES_LIMIT              = 1000
	MAX_NUMBER_OF_ROWS int64 = 984930291881790848 - 10000000000000000 //63^10-1 - ~5%
)

var (
	ErrOutOfUniqueValues = errors.New("out of uniqie values")
	ErrDuplicateId       = errors.New("such id already exists")
	ErrNoSuchURL         = errors.New("no such url")
)

type Encoder interface {
	SaveBase63(url string, id int64, counter *int) error
	IfExistsBase63(url string) (int64, error)
}

type Decoder interface {
	GetBase63(id int64) (string, error)
}

type Repository struct {
	Encoder
	Decoder
}

func NewRepository(db *database.Database) *Repository {
	switch db.DBType {
	case database.POSTGRES:
		return &Repository{
			Encoder: NewEncoderPostgres(db.Sql),
			Decoder: NewDecoderPostgres(db.Sql),
		}
	case database.REDIS:
		return &Repository{
			Encoder: NewEncoderRedis(db.Redis),
			Decoder: NewDecoderRedis(db.Redis),
		}
	}
	return nil
}
