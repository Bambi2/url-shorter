package repository

import (
	"errors"

	"github.com/bambi2/url-shorter/internal/database"
)

var (
	ErrOutOfUniqueValues = errors.New("out of uniqie values")
	ErrDuplicateId       = errors.New("such id already exists")
	ErrNoSuchURL         = errors.New("no such url")
)

type Encoder interface {
	SaveBase63(url string, id int64) error
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
