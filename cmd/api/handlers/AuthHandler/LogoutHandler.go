package AuthHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/internal/application"
)

func LogoutHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		clearCookie, err := app.JWTService.ClearCookie()

		if err == nil {
			c.SetCookie(clearCookie)
		}

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Logged out successfully",
		})
	}
}
