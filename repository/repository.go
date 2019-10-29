package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository interface {
	UserRepository
	SessionRepository
	Close() error
}

func New(dsn string) (Repository, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) Close() error {
	return r.db.Close()
}
