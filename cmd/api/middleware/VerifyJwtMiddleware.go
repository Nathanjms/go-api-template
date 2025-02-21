package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/internal/application"
)

// JWTAuthMiddleware is a middleware function that verifies JWT tokens.
func JWTAuthMiddleware(app *application.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("jwt")
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, application.Response{
					Success: false,
					Message: "Unauthorized",
				})
			}

			userId, err := app.JWTService.GetUserIdFromJWT(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, application.Response{
					Success: false,
					Message: "Unauthorized",
				})
			}

			c.Set("userId", userId)
			return next(c)
		}
	}
}
