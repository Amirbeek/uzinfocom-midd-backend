package store

import "github.com/Amirbeek/uzinfocom-midd-backend/internal/database"

type Store struct {
	db database.Service
}

func NewStore(db database.Service) *Store {
	return &Store{db: db}
}
