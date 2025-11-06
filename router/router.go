package router

import (
	"net/http"

	"wallet/controller"

	"github.com/gin-gonic/gin"
)

// SetupRouter set router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Wallet API is running",
		})
	})

	// API router group
	api := r.Group("/api/v1")
	{
		// users
		users := api.Group("/users")
		{
			users.POST("", controller.RegisterUser)
			users.GET("", controller.GetAllUsers)
			users.GET("/:id", controller.GetUser)
		}

		// wallets
		wallets := api.Group("/wallets")
		{
			wallets.GET("/:user_id/balance", controller.GetBalance)
			wallets.POST("/:user_id/deposit", controller.Deposit)
			wallets.POST("/:user_id/withdraw", controller.Withdraw)
			wallets.POST("/transfer", controller.Transfer)
		}

		// transactions
		transactions := api.Group("/transactions")
		{
			transactions.GET("/:user_id", controller.GetUserTransactions)
		}
	}

	return r
}
