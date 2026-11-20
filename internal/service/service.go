package service

import (
	"github.com/PrimeraAizen/e-comm/config"
	"github.com/PrimeraAizen/e-comm/internal/repository"
	"github.com/PrimeraAizen/e-comm/pkg/logger"
)

type Service struct {
	AuthService           AuthService
	UserService           UserService
	ProductService        ProductService
	InteractionService    InteractionService
	RecommendationService RecommendationService
}

type Deps struct {
	Repos  *repository.Repository
	Config *config.Config
	Logger *logger.Logger
}

func NewServices(deps Deps) *Service {
	authService, err := NewAuthService(deps.Repos.User, deps.Logger, &deps.Config.JWT)
	if err != nil {
		panic("failed to create auth service: " + err.Error())
	}

	return &Service{
		AuthService:           authService,
		UserService:           NewUserService(deps.Repos.User, deps.Repos.Profile),
		ProductService:        NewProductService(deps.Repos.Product),
		InteractionService:    NewInteractionService(deps.Repos.Interaction, deps.Repos.Product),
		RecommendationService: NewRecommendationService(deps.Repos.Interaction, deps.Repos.Product),
	}
}
