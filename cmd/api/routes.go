package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/cmd/api/handlers/AuthHandler"
	"github.com/nathanjms/go-api-template/cmd/api/handlers/UserHandler"
	"github.com/nathanjms/go-api-template/cmd/api/middleware"
	"github.com/nathanjms/go-api-template/internal/application"
)

func InitRoutes(e *echo.Echo, app *application.Application) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Hello, World!",
		})
	})
	e.GET("status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "OK",
		})
	})

	e.POST("login", AuthHandler.LoginHandler(app))
	e.POST("register", AuthHandler.RegisterHandler(app))
	e.POST("logout", AuthHandler.LogoutHandler(app))

	authed := e.Group("")
	authed.Use(middleware.JWTAuthMiddleware(app))

	// --- AUTHED ROUTES ---

	// User Routes
	authed.GET("user", UserHandler.GetAccountHandler(app))
	authed.DELETE("user", UserHandler.DeleteAccountHandler(app))

}
