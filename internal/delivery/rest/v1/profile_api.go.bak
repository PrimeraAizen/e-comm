package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/PrimeraAizen/e-comm/internal/delivery/dto"
	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// InitProfileRoutes sets up profile endpoints
func (h *Handler) InitProfileRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	profiles := rg.Group("/profiles")
	profiles.Use(authMiddleware)
	{
		profiles.GET("/me", h.GetProfile)
		profiles.PUT("/me", h.UpdateProfile)
		profiles.PUT("/me/password", h.ChangePassword)
		profiles.DELETE("/me/account", h.DeleteAccount)
		profiles.GET("/me/interactions", h.GetMyInteractions)
		profiles.GET("/me/views", h.GetMyViewHistory)
		profiles.GET("/me/likes", h.GetMyLikedProducts)
		profiles.GET("/me/purchases", h.GetMyPurchases)
		profiles.GET("/me/recommendations", h.GetRecommendations)
		profiles.GET("/me/similar", h.GetSimilarUsers)
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information with detailed profile data
// @Tags profiles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ProfileResponse
// @Router /profiles/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	// Get user and profile
	user, profile, err := h.services.UserService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithComponent("profile").WithError(err).Error("Failed to get profile")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get profile"})
		return
	}

	// Build response
	response := dto.ProfileResponse{
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	if profile != nil {
		response.ID = profile.ID
		response.UserID = profile.UserID
		response.FirstName = profile.FirstName
		response.LastName = profile.LastName
		if profile.MiddleName != nil {
			response.MiddleName = *profile.MiddleName
		}
		if profile.DateOfBirth != nil {
			response.DateOfBirth = profile.DateOfBirth.Format("2006-01-02")
		}
		if profile.Gender != nil {
			response.Gender = *profile.Gender
		}
		if profile.Phone != nil {
			response.Phone = *profile.Phone
		}
		if profile.Address != nil {
			response.Address = *profile.Address
		}
		if profile.City != nil {
			response.City = *profile.City
		}
		if profile.Country != nil {
			response.Country = *profile.Country
		}
		if profile.PostalCode != nil {
			response.PostalCode = *profile.PostalCode
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's detailed profile information
// @Tags profiles
// @Accept json
// @Produce json
// @Param profile body dto.UpdateProfileRequest true "Profile update"
// @Security BearerAuth
// @Success 200 {object} dto.ProfileResponse
// @Router /profiles/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	// Parse request
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// Validate
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Build profile data
	profileData := &domain.Profile{
		UserID: userID,
	}

	if req.FirstName != nil {
		profileData.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profileData.LastName = *req.LastName
	}
	profileData.MiddleName = req.MiddleName
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid date_of_birth format, use YYYY-MM-DD"})
			return
		}
		profileData.DateOfBirth = &dob
	}
	profileData.Gender = req.Gender
	profileData.Phone = req.Phone
	profileData.Address = req.Address
	profileData.City = req.City
	profileData.Country = req.Country
	profileData.PostalCode = req.PostalCode

	// Update profile
	profile, err := h.services.UserService.UpdateProfile(c.Request.Context(), userID, profileData)
	if err != nil {
		h.logger.WithComponent("profile").WithError(err).Error("Failed to update profile")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update profile"})
		return
	}

	// Get user for response
	user, _, err := h.services.UserService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithComponent("profile").WithError(err).Error("Failed to get user")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get user"})
		return
	}

	// Build response
	response := dto.ProfileResponse{
		ID:        profile.ID,
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: profile.CreatedAt.Format(time.RFC3339),
		UpdatedAt: profile.UpdatedAt.Format(time.RFC3339),
	}

	if profile.MiddleName != nil {
		response.MiddleName = *profile.MiddleName
	}
	if profile.DateOfBirth != nil {
		response.DateOfBirth = profile.DateOfBirth.Format("2006-01-02")
	}
	if profile.Gender != nil {
		response.Gender = *profile.Gender
	}
	if profile.Phone != nil {
		response.Phone = *profile.Phone
	}
	if profile.Address != nil {
		response.Address = *profile.Address
	}
	if profile.City != nil {
		response.City = *profile.City
	}
	if profile.Country != nil {
		response.Country = *profile.Country
	}
	if profile.PostalCode != nil {
		response.PostalCode = *profile.PostalCode
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change current user's password
// @Tags profiles
// @Accept json
// @Produce json
// @Param password body dto.ChangePasswordRequest true "Password change"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /profiles/me/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	// Parse request
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// Validate
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Change password
	if err := h.services.UserService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		h.logger.WithComponent("profile").WithError(err).Error("Failed to change password")
		if err.Error() == "invalid current password" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "password changed successfully"})
}

// DeleteAccount godoc
// @Summary Delete account
// @Description Soft delete current user's account
// @Tags profiles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /profiles/me/account [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	// Delete account
	if err := h.services.UserService.DeleteAccount(c.Request.Context(), userID); err != nil {
		h.logger.WithComponent("profile").WithError(err).Error("Failed to delete account")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "account deleted successfully"})
}

// GetMyInteractions godoc
// @Summary Get my interactions
// @Description Get summary of current user's product interactions
// @Tags profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.UserInteractionSummary
// @Router /profiles/me/interactions [get]
func (h *Handler) GetMyInteractions(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	summary, err := h.services.InteractionService.GetUserInteractionSummary(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to get interactions")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get interactions"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetMyViewHistory godoc
// @Summary Get my view history
// @Description Get products the current user has viewed
// @Tags profiles
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /profiles/me/views [get]
func (h *Handler) GetMyViewHistory(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	views, err := h.services.InteractionService.GetUserViewHistory(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to get view history")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get view history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"views": views,
		"count": len(views),
	})
}

// GetMyLikedProducts godoc
// @Summary Get my liked products
// @Description Get products the current user has liked
// @Tags profiles
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /profiles/me/likes [get]
func (h *Handler) GetMyLikedProducts(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	likes, err := h.services.InteractionService.GetUserLikedProducts(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to get liked products")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get liked products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"likes": likes,
		"count": len(likes),
	})
}

// GetRecommendations godoc
// @Summary Get personalized product recommendations
// @Description Get product recommendations based on collaborative filtering
// @Tags profiles
// @Produce json
// @Param limit query int false "Number of recommendations" default(10)
// @Security BearerAuth
// @Success 200 {object} domain.RecommendationResponse
// @Router /profiles/me/recommendations [get]
func (h *Handler) GetRecommendations(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	recommendations, err := h.services.RecommendationService.GetRecommendations(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.WithComponent("recommendation").WithError(err).Error("Failed to get recommendations")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// GetSimilarUsers godoc
// @Summary Get similar users
// @Description Get users with similar interaction patterns
// @Tags profiles
// @Produce json
// @Param limit query int false "Number of similar users" default(10)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /profiles/me/similar [get]
func (h *Handler) GetSimilarUsers(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	similarUsers, err := h.services.RecommendationService.GetSimilarUsers(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.WithComponent("recommendation").WithError(err).Error("Failed to get similar users")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get similar users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"similar_users": similarUsers,
		"count":         len(similarUsers),
	})
}

// GetMyPurchases godoc
// @Summary Get my purchases
// @Description Get products the current user has purchased
// @Tags profiles
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /profiles/me/purchases [get]
func (h *Handler) GetMyPurchases(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	purchases, err := h.services.InteractionService.GetUserPurchaseHistory(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to get purchase history")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get purchase history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases": purchases,
		"count":     len(purchases),
	})
}
