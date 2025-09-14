package service

import (
	"e-comm/config"
	"e-comm/internal/repository"
)

type Service struct{}

type Deps struct {
	Repos  *repository.Repository
	Config *config.Config
}

func NewServices(deps Deps) *Service {
	return &Service{}
}
