package handler

import (
	"backend-2/api/cmd/db/model"
	"time"

	"net/http"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func HelloWorld() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func GetQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var q []*model.Quote

		if err := db.Find(&q).Error; err != nil {
			// error handling here
			return err
		}

		return c.JSON(http.StatusOK, q)
	}
}

type QuotePost struct {
	Quote  string `json:"quote" form:"quote" query:"quote"`
}

func CreateQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		q := new(QuotePost)
		if err = c.Bind(q); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quote := model.Quote{
			Quote: q.Quote,
			CreatedAt: datatypes.Date(time.Now()),
			UpdatedAt: datatypes.Date(time.Now()),
		}
		db.Create(&quote)
		return c.JSON(http.StatusOK, quote)
	}
}
