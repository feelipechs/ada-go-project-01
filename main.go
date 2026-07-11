package main

import (
	"context"
	"log"
	"net/http"

	"orders-api/config"
	"orders-api/controllers"
	"orders-api/database"
	"orders-api/repository"
	"orders-api/routes"
	"orders-api/service"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer pool.Close()

	clientRepo := repository.NewClientRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)

	clientService := service.NewClientService(clientRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, clientRepo, productRepo)

	clientController := controllers.NewClientController(clientService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)

	r := gin.Default()

	routes.Register(r, clientController, productController, orderController)

	log.Printf("server running on http://localhost:%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
