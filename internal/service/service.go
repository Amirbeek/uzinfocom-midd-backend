package service

import "github.com/Amirbeek/uzinfocom-midd-backend/internal/store"

type service struct {
}

type Services struct {
}

func NewServices(s *store.Store) *Services {
	return &Services{}
}
