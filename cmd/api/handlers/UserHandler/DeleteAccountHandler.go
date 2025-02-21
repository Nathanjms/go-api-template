package UserHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/internal/application"
)

func DeleteAccountHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("userId").(int64)

		if userId == 0 {
			return c.JSON(http.StatusUnauthorized, application.Response{
				Success: false,
				Message: "Unauthorized",
			})
		}

		// Delete the user from the database:
		if err := app.DB.UserModel.Delete(int64(userId)); err != nil {
			return err
		}

		clearCookie, err := app.JWTService.ClearCookie()

		if err == nil {
			c.SetCookie(clearCookie)
		}

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Account Deleted",
		})
	}
}
