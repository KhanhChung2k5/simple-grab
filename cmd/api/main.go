package main

import (
	"context"
	"log"

	"github.com/KhanhChung2k5/simple-grab/internal/config"	
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
	"github.com/KhanhChung2k5/simple-grab/internal/handler"
	"github.com/KhanhChung2k5/simple-grab/internal/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	db, err := repository.ConnectDB(ctx, cfg.DATABASE_URL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}
	defer sqlDB.Close()

	_ = repository.NewUserRepository(db)
	log.Println("database connected successfully")

	healthHandler := handler.NewHealthHandler()

	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.Timeout(30 * time.Second),
		middleware.ErrorHandler(),
	)
	
	//Group the routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)
	}

	log.Printf("server starting on: %s", cfg.PORT)
	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}