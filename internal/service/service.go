package service

import "github.com/Amirbeek/uzinfocom-midd-backend/internal/store"

type service struct {
}

type Services struct {
	Store store.Service
}

func NewServices(s store.Service) *Services {
	return &Services{Store: s}
}
