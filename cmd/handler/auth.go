package handler

import (
	"backend-2/api/cmd/config"
	"backend-2/api/cmd/db/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthLoginStruct struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := new(AuthLoginStruct)
		if err := c.Bind(body); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		username := body.Username
		password := body.Password

		user := new(model.User)

		db.Where("username = ?", username).First(&user)

		// Throws unauthorized error
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if err != nil {
			return echo.ErrForbidden
		}

		token, err := config.GenerateJwt(username, user.IsAdmin)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token": token,
		})
	}
}
