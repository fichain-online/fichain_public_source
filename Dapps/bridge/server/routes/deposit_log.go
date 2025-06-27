package routes

import (
	"github.com/gin-gonic/gin"

	"FichainBridge/controllers"
)

func SetupDepositLogRoute(
	router *gin.RouterGroup,
	controller *controllers.DepositLogController,
) {
	// Add auth later, now i'm too busy to do it
	router.GET("/deposit-logs/:address", controller.GetDepositLogs)
}
