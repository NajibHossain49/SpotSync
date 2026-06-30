package main

import (
	"net/http"

	"spotsync-api/config"
	"spotsync-api/database"
	"spotsync-api/dto"
	"spotsync-api/handler"
	appmiddleware "spotsync-api/middleware"
	"spotsync-api/repository"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load config + connect database.
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)

	// 2. Manual Dependency Injection: Repository -> Service -> Handler.
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	// 3. Setup Echo.
	e := echo.New()
	e.Validator = utils.NewValidator()

	// Global middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Health check.
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, dto.NewSuccess("SpotSync API is running 🚗", nil))
	})

	// 4. Routes.
	api := e.Group("/api/v1")

	// --- Auth (public) ---
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	jwtMW := appmiddleware.JWTAuth(cfg.JWTSecret)

	// --- Zones ---
	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAll)     // public
	zones.GET("/:id", zoneHandler.GetByID) // public
	// Admin-only zone management.
	zones.POST("", zoneHandler.Create, jwtMW, appmiddleware.AdminOnly)
	zones.PUT("/:id", zoneHandler.Update, jwtMW, appmiddleware.AdminOnly)
	zones.DELETE("/:id", zoneHandler.Delete, jwtMW, appmiddleware.AdminOnly)

	// --- Reservations ---
	reservations := api.Group("/reservations")
	reservations.POST("", reservationHandler.Reserve, jwtMW)                            // authenticated
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations, jwtMW)   // authenticated
	reservations.DELETE("/:id", reservationHandler.Cancel, jwtMW)                       // owner/admin
	reservations.GET("", reservationHandler.GetAll, jwtMW, appmiddleware.AdminOnly)     // admin only

	// 5. Start.
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
