package v1

import (
	"net/http"
	"strconv"

	"github.com/PrimeraAizen/e-comm/internal/delivery/dto"
	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitCategoryRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	categories := api.Group("/categories")
	categories.Use(authMiddleware)
	{
		categories.GET("", h.ListCategories)
		categories.GET("/:id", h.GetCategory)

		categories.POST("", h.CreateCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}

// Categories endpoints

// ListCategories godoc
// @Summary List categories
// @Description Get all product categories
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Category
// @Router /categories [get]
func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.services.ProductService.ListCategories(c.Request.Context())
	if err != nil {
		h.logger.WithComponent("product").WithError(err).Error("Failed to list categories")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to list categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Get detailed information about a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} domain.Category
// @Router /categories/{id} [get]
func (h *Handler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id"})
		return
	}

	category, err := h.services.ProductService.GetCategory(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to get category")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory godoc
// @Summary Create category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} domain.Category
// @Router /categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// TODO: Check if user has admin role

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := h.services.ProductService.CreateCategory(c.Request.Context(), category); err != nil {
		if err == domain.ErrAlreadyExists {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "category already exists"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to create category")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update an existing category 
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param category body dto.UpdateCategoryRequest true "Category data"
// @Success 200 {object} domain.Category
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// TODO: Check if user has admin role

	// Get existing category first
	existingCategory, err := h.services.ProductService.GetCategory(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to get category")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get category"})
		return
	}

	// Update only provided fields
	if req.Name != nil {
		existingCategory.Name = *req.Name
	}
	if req.Description != nil {
		existingCategory.Description = *req.Description
	}
	if req.ParentID != nil {
		existingCategory.ParentID = req.ParentID
	}

	if err := h.services.ProductService.UpdateCategory(c.Request.Context(), existingCategory); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to update category")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingCategory)
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category 
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 204
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category id"})
		return
	}

	if err := h.services.ProductService.DeleteCategory(c.Request.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to delete category")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to delete category"})
		return
	}

	c.Status(http.StatusNoContent)
}
