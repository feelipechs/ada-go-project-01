package routes

import (
	"net/http"

	"orders-api/controllers"

	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	clientController *controllers.ClientController,
	productController *controllers.ProductController,
	orderController *controllers.OrderController,
) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/clients", clientController.Create)
	r.GET("/clients", clientController.FindAll)
	r.GET("/clients/:id", clientController.FindByID)

	r.POST("/products", productController.Create)
	r.GET("/products", productController.FindAll)
	r.GET("/products/:id", productController.FindByID)

	r.POST("/orders", orderController.Create)
	r.GET("/orders", orderController.FindAll)
	r.GET("/orders/:id", orderController.FindByID)
	r.POST("/orders/:id/pay", orderController.Pay)
	r.POST("/orders/:id/cancel", orderController.Cancel)
}
