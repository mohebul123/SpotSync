package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/mohebul123/SpotSync/config"
	"github.com/mohebul123/SpotSync/handler"
	customMiddleware "github.com/mohebul123/SpotSync/middleware"
	"github.com/mohebul123/SpotSync/repository"
	"github.com/mohebul123/SpotSync/service"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, reading from system environment variables")
	}

	config.ConnectDatabase()
	db := config.DB
	e := echo.New()
	v := validator.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	authRepo := repository.NewAuthRepository(db)
	authSrv := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authSrv, v)

	zoneRepo := repository.NewZoneRepository(db)
	zoneSrv := service.NewZoneService(zoneRepo)
	zoneHandler := handler.NewZoneHandler(zoneSrv, v)

	resRepo := repository.NewReservationRepository(db)
	resSrv := service.NewReservationService(resRepo, zoneRepo)
	resHandler := handler.NewReservationHandler(resSrv, v)

	api := e.Group("/api")
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(customMiddleware.JWTMiddleware)

	protected.GET("/zones", zoneHandler.GetAll)
	protected.GET("/zones/:id", zoneHandler.GetByID)

	adminOnly := protected.Group("")
	adminOnly.Use(customMiddleware.RequireAdmin)
	adminOnly.POST("/zones", zoneHandler.Create)

	protected.POST("/reservations", resHandler.Book)
	protected.POST("/reservations/:id/cancel", resHandler.Cancel)
	protected.GET("/reservations/my", resHandler.GetMyReservations)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server is running smoothly on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
