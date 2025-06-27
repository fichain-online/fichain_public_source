package controllers

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"FichainCore/cmd/explorer/models"
	"FichainCore/cmd/explorer/requests"
)

// TransactionController handles transaction-related API requests.
// It holds a reference to the database connection.
type TransactionController struct {
	DB *gorm.DB
}

// NewTransactionController creates and returns a new TransactionController.
func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{DB: db}
}

func (ctrl *TransactionController) GetTransactions(c *gin.Context) {
	// 1. Pagination parameters (same as before)
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 2. Database query (same as before)
	var dbTransactions []models.TransactionDB // Note the name change for clarity
	var total int64

	if err := ctrl.DB.Model(&models.TransactionDB{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transactions"})
		return
	}

	offset := (page - 1) * pageSize
	err = ctrl.DB.Model(&models.TransactionDB{}).
		Preload("Receipt").
		Preload("Logs").
		Order("block_height desc, transaction_index desc").
		Offset(offset).
		Limit(pageSize).
		Find(&dbTransactions).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch transactions from database"},
		)
		return
	}

	// 3. CONVERT database models to response DTOs
	responseTxs := make([]*requests.TransactionResponse, len(dbTransactions))
	for i, tx := range dbTransactions {
		responseTxs[i] = requests.ToTransactionResponse(&tx)
	}

	// 4. Return the DTOs in the response
	c.JSON(http.StatusOK, gin.H{
		"data":     responseTxs, // Use the converted slice
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// GetTransactionsByAddress fetches a paginated list of transactions for a specific address.
func (ctrl *TransactionController) GetTransactionsByAddress(c *gin.Context) {
	// Note: Since you're using raw bytes now, your HexToAddress logic will need to adapt.
	// Assuming common.HexToAddress returns []byte
	addressHex := c.Param("address")
	// A simple hex validation might be needed here if common.IsHexAddress is gone
	address, err := hex.DecodeString(addressHex[2:]) // Decode hex, skipping "0x"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
		return
	}

	// 1. Pagination (same as before)
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 2. Database query (same as before)
	var dbTransactions []models.TransactionDB
	var total int64

	query := ctrl.DB.Model(&models.TransactionDB{}).
		Where("from_address = ? OR to_address = ?", address, address)

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transactions"})
		return
	}

	offset := (page - 1) * pageSize
	err = query.
		Preload("Receipt").
		Preload("Logs").
		Order("block_height desc, transaction_index desc").
		Offset(offset).
		Limit(pageSize).
		Find(&dbTransactions).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch transactions from database"},
		)
		return
	}

	// 3. CONVERT database models to response DTOs
	responseTxs := make([]*requests.TransactionResponse, len(dbTransactions))
	for i, tx := range dbTransactions {
		responseTxs[i] = requests.ToTransactionResponse(&tx)
	}

	// 4. Return the DTOs in the response
	c.JSON(http.StatusOK, gin.H{
		"data":     responseTxs, // Use the converted slice
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}
