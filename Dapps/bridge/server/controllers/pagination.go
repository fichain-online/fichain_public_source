package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// Paginate returns a GORM scope function that can be used to apply
// pagination to a query. It reads 'page' and 'limit' from the Gin context.
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
		if page <= 0 {
			page = DefaultPage
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(DefaultLimit)))
		switch {
		case limit > MaxLimit:
			limit = MaxLimit
		case limit <= 0:
			limit = DefaultLimit
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
