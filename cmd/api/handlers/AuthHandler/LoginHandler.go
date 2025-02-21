package AuthHandler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/internal/application"

	"golang.org/x/crypto/bcrypt"
)

type LoginJsonUser struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

func LoginHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		loginUserRequest := new(LoginJsonUser)
		if err := c.Bind(loginUserRequest); err != nil {
			return c.JSON(http.StatusBadRequest, application.Response{
				Success: false,
				Message: "Error passing JSON",
			})
		}

		// Validate not empty:
		if loginUserRequest.Username == "" || loginUserRequest.Password == "" {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Username and password are required",
			})
		}

		user, err := app.DB.UserModel.GetByUsername(loginUserRequest.Username)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		// Compare password using bcrypt, failing if they don't match
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserRequest.Password)); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		jwtCookie, err := app.JWTService.CreateJwtCookie(user.ID, loginUserRequest.Username, loginUserRequest.RememberMe)

		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(jwtCookie)

		// Set the cookie:
		c.SetCookie(jwtCookie)

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Success",
			Data: application.ResponseData{
				"id":       user.ID,
				"username": user.Username,
			},
		})
	}
}
