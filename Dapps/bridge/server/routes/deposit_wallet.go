package routes

import (
	"github.com/gin-gonic/gin"

	"FichainBridge/controllers"
)

func SetupDepositWalletRoute(
	router *gin.RouterGroup,
	controller *controllers.DepositWalletController,
) {
	// Add auth later, now i'm too busy to do it
	router.GET("/deposit-wallet/:tokenName/:address", controller.GetDepositWallet)
}
