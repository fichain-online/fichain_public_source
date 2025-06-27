package controllers

import (
	"net/http"
	"strings"

	"FichainCore/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"FichainBridge/models"
	"FichainBridge/requests" // Your requests DTO package
)

type DepositLogController struct {
	DB *gorm.DB
}

func NewDepositLogController(
	db *gorm.DB,
) *DepositLogController {
	return &DepositLogController{
		DB: db,
	}
}

// GetDepositLogs fetches a paginated list of deposit logs for a specific Fichain address.
// It also supports optional filtering by token name.
// GET /api/v1/deposit-logs/:address?page=1&limit=20&token=USDT
func (ctrl *DepositLogController) GetDepositLogs(c *gin.Context) {
	// 1. Get and Validate Path Parameter (Address)
	fichainAddressHex := c.Param("address")
	if !common.IsHexAddress(fichainAddressHex) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Fichain address format"})
		return
	}
	fichainAddressBytes := common.HexToAddress(fichainAddressHex).Bytes()

	// 2. Build the Database Query
	var logs []*models.DepositLog

	// Start with the base query, filtering by the user's address.
	query := ctrl.DB.Model(&models.DepositLog{}).
		Where("fichain_address = ?", fichainAddressBytes)

	// 3. Apply Optional Filters (e.g., by token name)
	tokenName := c.Query("token")
	if tokenName != "" {
		// Sanitize input to be safe, although this case is simple.
		cleanTokenName := strings.ToUpper(strings.TrimSpace(tokenName))
		query = query.Where("token_name = ?", cleanTokenName)
	}

	// Add ordering to ensure consistent results across pages.
	// Ordering by creation time descending is most common for transaction histories.
	query = query.Order("created_at DESC")

	// 4. Execute the Query with Pagination
	// The .Scopes(Paginate(c)) applies the LIMIT and OFFSET from our helper.
	result := query.Scopes(Paginate(c)).Find(&logs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve deposit logs"})
		return
	}

	// To provide full pagination metadata, we should also get the total count.
	var totalRecords int64
	// We use the same query but remove Limit/Offset to count all matching records.
	countQuery := ctrl.DB.Model(&models.DepositLog{}).
		Where(query.Statement.Clauses["WHERE"].Expression).
		Count(&totalRecords)
	if countQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count deposit logs"})
		return
	}

	// 5. Convert to DTO and Return Response
	// Use the list converter we created earlier.
	responseLogs := requests.ToDepositLogListResponse(logs)

	// Return the data along with pagination info.
	c.JSON(http.StatusOK, gin.H{
		"data": responseLogs,
		"pagination": gin.H{
			"total": totalRecords,
			"page":  c.DefaultQuery("page", "1"),
			"limit": c.DefaultQuery("limit", "10"),
		},
	})
}
