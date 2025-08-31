package controllers

import (
	"net/http"
	"strconv"

	"stokq-backend/config"
	"stokq-backend/dto"
	"stokq-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Check if SKU already exists
	var existingProduct models.Product
	if err := config.DB.Where("sku = ?", req.SKU).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error: "SKU already exists",
		})
		return
	}

	// Create product
	product := models.Product{
		SKU:   req.SKU,
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create product",
		})
		return
	}

	// Return response
	response := dto.ProductResponse{
		ID:        product.ID,
		SKU:       product.SKU,
		Name:      product.Name,
		Stock:     product.Stock,
		Price:     product.Price,
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Product created successfully",
		Data:    response,
	})
}

func GetProducts(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to fetch products",
		})
		return
	}

	// Convert to response format
	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, dto.ProductResponse{
			ID:        product.ID,
			SKU:       product.SKU,
			Name:      product.Name,
			Stock:     product.Stock,
			Price:     product.Price,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Products retrieved successfully",
		Data:    responses,
	})
}

func GetProductByID(c *gin.Context) {
	// Get product ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Return response
	response := dto.ProductResponse{
		ID:        product.ID,
		SKU:       product.SKU,
		Name:      product.Name,
		Stock:     product.Stock,
		Price:     product.Price,
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Product retrieved successfully",
		Data:    response,
	})
}

func UpdateProduct(c *gin.Context) {
	// Get product ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	var req dto.UpdateProductRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Find existing product
	var product models.Product
	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Check if SKU already exists for other products
	if req.SKU != "" && req.SKU != product.SKU {
		var existingProduct models.Product
		if err := config.DB.Where("sku = ? AND id != ?", req.SKU, id).First(&existingProduct).Error; err == nil {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: "SKU already exists for another product",
			})
			return
		}
	}

	// Update fields if provided
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Price > 0 {
		product.Price = req.Price
	}

	// Save changes
	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update product",
		})
		return
	}

	// Return response
	response := dto.ProductResponse{
		ID:        product.ID,
		SKU:       product.SKU,
		Name:      product.Name,
		Stock:     product.Stock,
		Price:     product.Price,
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Product updated successfully",
		Data:    response,
	})
}

func DeleteProduct(c *gin.Context) {
	// Get product ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid product ID",
		})
		return
	}

	// Find existing product
	var product models.Product
	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Database error",
		})
		return
	}

	// Delete product
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Product deleted successfully",
	})
}
