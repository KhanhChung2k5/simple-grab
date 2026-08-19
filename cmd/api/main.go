package main

import (
	"context"
	"log"

	"github.com/KhanhChung2k5/simple-grab/internal/config"
	"github.com/KhanhChung2k5/simple-grab/internal/handler"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
	"github.com/KhanhChung2k5/simple-grab/internal/server"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
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

	// initialize the repositories and services
	userRepo := repository.NewUserRepository(db)
	rideRepo := repository.NewRideRepository(db)
	driverRepo := repository.NewDriverRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWT_SECRET, cfg.JWT_EXPIRY_HOURS)
	rideService := service.NewRideService(rideRepo, driverRepo)
	driverService := service.NewDriverService(driverRepo)

	authHandler := handler.NewAuthHandler(authService)
	rideHandler := handler.NewRideHandler(rideService)
	driverHandler := handler.NewDriverHandler(driverService)
	healthHandler := handler.NewHealthHandler()

	log.Println("database connected successfully")

	gin.SetMode(cfg.GIN_MODE)
	router := server.NewRouter(authService, authHandler, rideHandler, driverHandler, healthHandler)

	log.Printf("server starting on: %s", cfg.PORT)
	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
