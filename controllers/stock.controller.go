package controllers

import (
	"net/http"

	"stokq-backend/config"
	"stokq-backend/dto"
	"stokq-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StockIn(c *gin.Context) {
	var req dto.StockTransactionRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	// Find product
	var product models.Product
	if err := tx.First(&product, req.ProductID).Error; err != nil {
		tx.Rollback()
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

	// Add stock
	product.Stock += req.Quantity

	// Save changes
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update stock",
		})
		return
	}

	// Commit transaction
	tx.Commit()

	// Return response
	response := dto.StockTransactionResponse{
		Message: "Stock added successfully",
		Product: dto.ProductResponse{
			ID:        product.ID,
			SKU:       product.SKU,
			Name:      product.Name,
			Stock:     product.Stock,
			Price:     product.Price,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}

func StockOut(c *gin.Context) {
	var req dto.StockTransactionRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	// Find product
	var product models.Product
	if err := tx.First(&product, req.ProductID).Error; err != nil {
		tx.Rollback()
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

	// Check if stock is sufficient
	if product.Stock < req.Quantity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Insufficient stock available",
		})
		return
	}

	// Reduce stock
	product.Stock -= req.Quantity

	// Save changes
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update stock",
		})
		return
	}

	// Commit transaction
	tx.Commit()

	// Return response
	response := dto.StockTransactionResponse{
		Message: "Stock reduced successfully",
		Product: dto.ProductResponse{
			ID:        product.ID,
			SKU:       product.SKU,
			Name:      product.Name,
			Stock:     product.Stock,
			Price:     product.Price,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, response)
}
