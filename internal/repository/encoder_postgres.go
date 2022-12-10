package repository

import (
	"fmt"

	"github.com/bambi2/url-shorter/internal/database"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type EncoderPostgres struct {
	db *sqlx.DB
}

func NewEncoderPostgres(db *sqlx.DB) *EncoderPostgres {
	return &EncoderPostgres{db: db}
}

func (r *EncoderPostgres) SaveBase63(url string, id int64) error {
	query := fmt.Sprintf("INSERT INTO %s (id, url) values ($1, $2) RETURNING id", database.Base63Table)
	row := r.db.QueryRow(query, id, url)
	if err := row.Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return ErrDuplicateId
			}
		}
		return err
	}

	return nil
}
