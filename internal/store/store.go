package store

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store { return &Store{DB: db} }
