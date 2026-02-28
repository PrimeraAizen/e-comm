package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/PrimeraAizen/e-comm/internal/service"
	"github.com/PrimeraAizen/e-comm/pkg/logger"
)

type Handler struct {
	services *service.Service
	logger   *logger.Logger
}

func NewHandler(services *service.Service, appLogger *logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   appLogger,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// Public routes
	h.InitAuthRoutes(v1)

	// Protected routes (require authentication)
	// authMiddleware := middleware.AuthMiddleware(h.services.AuthService)
	// h.InitCategoryRoutes(v1, authMiddleware)
	// h.InitProductRoutes(v1, authMiddleware)
	// h.InitProfileRoutes(v1, authMiddleware)
}
