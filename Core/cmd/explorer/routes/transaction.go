package routes

import (
	"github.com/gin-gonic/gin"

	"FichainCore/cmd/explorer/controllers"
)

func SetupTransactionRoute(router *gin.RouterGroup, controller *controllers.TransactionController) {
	// Add auth later, now i'm too busy to do it
	router.GET("/transaction/:address", controller.GetTransactionsByAddress)
	router.GET("/transactions", controller.GetTransactions)
}
