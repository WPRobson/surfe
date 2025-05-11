package main

import (
	"fmt"
	"net/http"
	"os"
	"surfe/internal/handlers"
	"surfe/internal/repository"
	"surfe/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HideBanner = true

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return fmt.Errorf("failed to create user repository: %v", err)
	}

	actionsRepo, err := repository.NewActionRepository()
	if err != nil {
		return fmt.Errorf("failed to create action repository: %v", err)
	}

	userService := services.NewUserService(userRepo, actionsRepo)
	actionsService := services.NewActionService(actionsRepo)

	userHandler := handlers.NewUserHandler(userService)
	actionHandler := handlers.NewActionHandler(actionsService)

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/users/:id", userHandler.GetUserByID)
	v1.GET("/users/:id/actions/count", userHandler.GetUserActionCount)
	v1.GET("/actions/:type/next", actionHandler.GetNextActionProbabilities)
	v1.GET("/actions/referral", actionHandler.GetReferralIndex)

	if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}
