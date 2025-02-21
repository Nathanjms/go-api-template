package UserHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/internal/application"
)

func GetAccountHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get("userId").(int64)

		if userId == 0 {
			return c.JSON(http.StatusUnauthorized, application.Response{
				Success: false,
				Message: "Unauthorized",
			})
		}

		// Delete the user from the database:
		user, err := app.DB.UserModel.FindUser(int64(userId))

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Account Details Retrieved",
			Data: application.ResponseData{
				"user": user,
			},
		})
	}
}
