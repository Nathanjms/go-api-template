package AuthHandler

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/go-api-template/internal/application"
)

type RegisterJsonUser struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	RememberMe      bool   `json:"rememberMe"`
}

func RegisterHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		newUserRequest := new(RegisterJsonUser)
		if err := c.Bind(newUserRequest); err != nil {
			return c.JSON(http.StatusBadRequest, application.Response{
				Success: false,
				Message: "Error parsing JSON",
			})
		}

		// Basic validations:
		if newUserRequest.Username == "" || newUserRequest.Password == "" {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Username and password are required",
				Errors:  map[string][]string{"username": {"Username and password are required"}, "password": {"Username and password are required"}},
			})
		}

		if len(newUserRequest.Password) < 8 {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Password must be at least 8 characters",
				Errors:  map[string][]string{"password": {"Password must be at least 8 characters"}},
			})
		}

		if newUserRequest.Password != newUserRequest.PasswordConfirm {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Passwords do not match",
				Errors:  map[string][]string{"passwordConfirm": {"Passwords do not match"}},
			})
		}

		// Define the regular expression for email validation
		validRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`

		// Compile the regular expression
		re := regexp.MustCompile(validRegex)

		// Check if the username (email) is valid
		if !re.MatchString(newUserRequest.Username) {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Invalid email address",
				Errors:  map[string][]string{"username": {"Invalid email address"}},
			})
		}

		// Ensure does not exist:
		_, err := app.DB.UserModel.GetByUsername(newUserRequest.Username)
		if err == nil {
			return c.JSON(http.StatusUnprocessableEntity, application.Response{
				Success: false,
				Message: "Username already exists",
				Errors:  map[string][]string{"username": {"Username already exists"}},
			})
		}

		var newUserId int64
		newUserId, err = app.DB.UserModel.Create(newUserRequest.Username, newUserRequest.Password)
		if err != nil {
			return err
		}

		jwtCookie, err := app.JWTService.CreateJwtCookie(newUserId, newUserRequest.Username, newUserRequest.RememberMe)

		if err != nil {
			return err
		}

		// Set the cookie:
		c.SetCookie(jwtCookie)

		return c.JSON(http.StatusOK, application.Response{
			Success: true,
			Message: "Success",
		})
	}
}
