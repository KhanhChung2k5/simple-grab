package main

import (
	"context"
	"log"
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/config"
	"github.com/KhanhChung2k5/simple-grab/internal/handler"
	"github.com/KhanhChung2k5/simple-grab/internal/middleware"
	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
	"github.com/KhanhChung2k5/simple-grab/internal/response"
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

	authService := service.NewAuthService(userRepo, cfg.JWT_SECRET, cfg.JWT_EXPIRY_HOURS)
	rideService := service.NewRideService(rideRepo)

	authHandler := handler.NewAuthHandler(authService)
	rideHandler := handler.NewRideHandler(rideService)
	healthHandler := handler.NewHealthHandler()

	log.Println("database connected successfully")

	// initialize the router
	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.Timeout(30*time.Second),
		middleware.ErrorHandler(),
	)

	// initialize the routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)
		// auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
		}

		// protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/rider/ping", middleware.RequireRole(model.RoleRider), func(c *gin.Context) {
				response.Success(c, gin.H{"message": "rider access ok"})
			})
			protected.GET("/driver/ping", middleware.RequireRole(model.RoleDriver), func(c *gin.Context) {
				response.Success(c, gin.H{"message": "driver access ok"})
			})

			rides := protected.Group("/rides")
			{
				rides.POST("", middleware.RequireRole(model.RoleRider), rideHandler.Create)
				rides.GET("/available", middleware.RequireRole(model.RoleDriver), rideHandler.ListAvailable)
				rides.GET("/:id", rideHandler.GetByID)
				rides.GET("", rideHandler.List)
				rides.POST("/:id/accept", middleware.RequireRole(model.RoleDriver), rideHandler.Accept)
				rides.PATCH("/:id/status", middleware.RequireRole(model.RoleDriver), rideHandler.UpdateStatus)
			}
		}
	}

	log.Printf("server starting on: %s", cfg.PORT)
	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
