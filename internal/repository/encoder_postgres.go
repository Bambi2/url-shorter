package repository

import (
	"database/sql"
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

func (r *EncoderPostgres) IfExistsBase63(url string) (int64, error) {
	var id int64
	query := fmt.Sprintf("SELECT id FROM %s WHERE url=$1", database.Base63Table)
	row := r.db.QueryRow(query, url)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		} else {
			return -1, err
		}
	}

	return id, nil
}
