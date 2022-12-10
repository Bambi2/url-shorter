package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bambi2/url-shorter/internal/database"
	"github.com/jmoiron/sqlx"
)

type DecoderPostgres struct {
	db *sqlx.DB
}

func NewDecoderPostgres(db *sqlx.DB) *DecoderPostgres {
	return &DecoderPostgres{db: db}
}

func (r *DecoderPostgres) GetBase63(id int64) (string, error) {
	var url string
	query := fmt.Sprintf("SELECT url FROM %s WHERE id=$1", database.Base63Table)

	row := r.db.QueryRow(query, id)
	if err := row.Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoSuchURL
		}

		return "", err
	}

	return url, nil
}
