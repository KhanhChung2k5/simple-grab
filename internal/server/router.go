package server

import (
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/handler"
	"github.com/KhanhChung2k5/simple-grab/internal/middleware"
	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/gin-gonic/gin"
)

// NewRouter builds the HTTP API independently from process startup so the
// exact production routes can also be exercised by integration tests.
func NewRouter(
	auth middleware.AuthProvider,
	authHandler *handler.AuthHandler,
	rideHandler *handler.RideHandler,
	driverHandler *handler.DriverHandler,
	healthHandler *handler.HealthHandler,
) *gin.Engine {
	router := gin.New()
	// allow all proxies
	_ = router.SetTrustedProxies(nil)
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.CORS(),
		middleware.Timeout(30*time.Second),
		middleware.ErrorHandler(),
	)
	// version 1 routes
	v1 := router.Group("/api/v1")
	v1.GET("/health", healthHandler.Health)

	authRoutes := v1.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.GET("/me", middleware.AuthMiddleware(auth), authHandler.Me)
	// protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(auth))
	protected.GET("/rider/ping", middleware.RequireRole(model.RoleRider), func(c *gin.Context) {
		response.Success(c, gin.H{"message": "rider access ok"})
	})
	protected.GET("/driver/ping", middleware.RequireRole(model.RoleDriver), func(c *gin.Context) {
		response.Success(c, gin.H{"message": "driver access ok"})
	})
	// rides routes
	rides := protected.Group("/rides")
	rides.POST("", middleware.RequireRole(model.RoleRider), rideHandler.Create)
	rides.GET("/available", middleware.RequireRole(model.RoleDriver), rideHandler.ListAvailable)
	rides.GET("/:id", rideHandler.GetByID)
	rides.GET("", rideHandler.List)
	rides.POST("/:id/accept", middleware.RequireRole(model.RoleDriver), rideHandler.Accept)
	rides.POST("/:id/cancel", rideHandler.Cancel)
	rides.PATCH("/:id/status", middleware.RequireRole(model.RoleDriver), rideHandler.UpdateStatus)

	// drivers routes
	drivers := protected.Group("/drivers/me")
	drivers.Use(middleware.RequireRole(model.RoleDriver))
	drivers.PATCH("/online", driverHandler.UpdateOnline)
	drivers.PATCH("/location", driverHandler.UpdateLocation)

	return router
}
