package main

import (
	"context"
	"log"
	"net/http"

	"orders-api/config"
	"orders-api/database"

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

	r := gin.Default()

	log.Printf("server running on http://localhost:%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
