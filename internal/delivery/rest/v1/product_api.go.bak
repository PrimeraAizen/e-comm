package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PrimeraAizen/e-comm/internal/delivery/dto"
	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// InitProductRoutes initializes product routes
func (h *Handler) InitProductRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	products := api.Group("/products")
	products.Use(authMiddleware)
	{
		products.GET("", h.ListProducts)
		products.GET("/:id", h.GetProduct)
		products.GET("/:id/statistics", h.GetProductStatistics)
		products.POST("", h.CreateProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)

		products.POST("/:id/view", h.RecordProductView)
		products.POST("/:id/like", h.LikeProduct)
		products.DELETE("/:id/like", h.UnlikeProduct)
		products.GET("/:id/liked", h.CheckProductLiked)
		products.POST("/:id/purchase", h.PurchaseProduct)
		products.GET("/:id/purchased", h.CheckProductPurchased)
	}
}

// ListProducts godoc
// @Summary List products
// @Description Get a paginated list of products with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param category_id query string false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param search query string false "Search in name and description"
// @Param sort_by query string false "Sort by: name, price, created_at" default(created_at)
// @Param sort_order query string false "Sort order: asc, desc" default(desc)
// @Success 200 {object} dto.ProductListResponse
// @Router /products [get]
func (h *Handler) ListProducts(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Build filter
	filter := domain.ProductFilter{
		Limit:       limit,
		Offset:      offset,
		SortBy:      c.Query("sort_by"),
		SortOrder:   c.Query("sort_order"),
		SearchQuery: c.Query("search"),
	}

	// Category filter
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category_id"})
			return
		}
		filter.CategoryID = &categoryID
	}

	// Price filters
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid min_price"})
			return
		}
		filter.MinPrice = &minPrice
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid max_price"})
			return
		}
		filter.MaxPrice = &maxPrice
	}

	// Get products with categories
	products, total, err := h.services.ProductService.ListProductsWithCategories(c.Request.Context(), filter)
	if err != nil {
		h.logger.WithComponent("product").WithError(err).Error("Failed to list products")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to list products"})
		return
	}

	c.JSON(http.StatusOK, dto.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get detailed information about a specific product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} domain.ProductWithCategory
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	product, err := h.services.ProductService.GetProductWithCategory(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to get product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} domain.ProductWithCategory
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// TODO: Check if user has admin role

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}

	if err := h.services.ProductService.CreateProduct(c.Request.Context(), product); err != nil {
		h.logger.WithComponent("product").WithError(err).Error("Failed to create product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update product information (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param product body dto.UpdateProductRequest true "Updated product data"
// @Success 200 {object} domain.ProductWithCategory
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// TODO: Check if user has admin role

	// Get existing product first
	existingProduct, err := h.services.ProductService.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to get product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get product"})
		return
	}

	// Update only provided fields
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.CategoryID != nil {
		existingProduct.CategoryID = req.CategoryID
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}
	if req.Stock != nil {
		existingProduct.Stock = *req.Stock
	}
	if req.ImageURL != nil {
		existingProduct.ImageURL = *req.ImageURL
	}
	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	if err := h.services.ProductService.UpdateProduct(c.Request.Context(), existingProduct); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to update product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingProduct)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	// TODO: Check if user has admin role

	if err := h.services.ProductService.DeleteProduct(c.Request.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to delete product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProductStatistics godoc
// @Summary Get product statistics
// @Description Get view count, like count, and purchase count for a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} domain.ProductStatistics
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id}/statistics [get]
func (h *Handler) GetProductStatistics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	stats, err := h.services.ProductService.GetProductStatistics(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product statistics not found"})
			return
		}
		h.logger.WithComponent("product").WithError(err).Error("Failed to get product statistics")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// RecordProductView godoc
// @Summary Record product view
// @Description Record that a user has viewed a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /products/{id}/view [post]
func (h *Handler) RecordProductView(c *gin.Context) {
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

	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	if err := h.services.InteractionService.RecordProductView(c.Request.Context(), userID, productID); err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to record view")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to record view"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "view recorded"})
}

// LikeProduct godoc
// @Summary Like a product
// @Description Add a product to user's liked products
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /products/{id}/like [post]
func (h *Handler) LikeProduct(c *gin.Context) {
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

	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	if err := h.services.InteractionService.LikeProduct(c.Request.Context(), userID, productID); err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to like product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to like product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product liked"})
}

// UnlikeProduct godoc
// @Summary Unlike a product
// @Description Remove a product from user's liked products
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /products/{id}/like [delete]
func (h *Handler) UnlikeProduct(c *gin.Context) {
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

	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	if err := h.services.InteractionService.UnlikeProduct(c.Request.Context(), userID, productID); err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to unlike product")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to unlike product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product unliked"})
}

// CheckProductLiked godoc
// @Summary Check if product is liked
// @Description Check if the current user has liked a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Success 200 {object} map[string]bool
// @Router /products/{id}/liked [get]
func (h *Handler) CheckProductLiked(c *gin.Context) {
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

	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	liked, err := h.services.InteractionService.IsProductLiked(c.Request.Context(), userID, productID)
	if err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to check if liked")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to check like status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

// PurchaseProduct godoc
// @Summary Purchase a product
// @Description Record a product purchase and update stock
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param purchase body dto.PurchaseProductRequest true "Purchase details"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /products/{id}/purchase [post]
func (h *Handler) PurchaseProduct(c *gin.Context) {
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

	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	var req dto.PurchaseProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "quantity must be greater than 0"})
		return
	}

	if err := h.services.InteractionService.PurchaseProduct(c.Request.Context(), userID, productID, req.Quantity); err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to purchase product")
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product purchased successfully"})
}

// CheckProductPurchased godoc
// @Summary Check if product is purchased
// @Description Check if the current user has purchased a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Success 200 {object} map[string]bool
// @Router /products/{id}/purchased [get]
func (h *Handler) CheckProductPurchased(c *gin.Context) {
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

	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	purchased, err := h.services.InteractionService.HasPurchasedProduct(c.Request.Context(), userID, productID)
	if err != nil {
		h.logger.WithComponent("interaction").WithError(err).Error("Failed to check if purchased")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to check purchase status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"purchased": purchased})
}
