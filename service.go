package main

import (
	"database/sql"
)

// db connection
type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}
