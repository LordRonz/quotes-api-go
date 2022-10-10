package handler

import (
	"backend-2/api/cmd/db/model"
	"backend-2/api/cmd/db/plugin"
	"crypto/rand"
	"math/big"
	"strconv"
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
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		page, _ := strconv.Atoi(c.QueryParam("page"))

		pagination := plugin.Pagination{
			Limit: limit,
			Page:  page,
		}

		res, err := listQuotes(db, pagination)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func listQuotes(db *gorm.DB, pagination plugin.Pagination) (*plugin.Pagination, error) {
	var quotes []*model.Quote

	db.Scopes(plugin.Paginate(quotes, &pagination, db)).Find(&quotes)
	pagination.Rows = quotes

	return &pagination, nil
}

func GetRandomQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var count int64

		db.Table("quotes").Count(&count)

		randomOffset, err := rand.Int(rand.Reader, big.NewInt(count))
		if err != nil {
			c.Logger().Error(err)
		}

		var q []*model.Quote

		if err := db.Limit(1).Offset(int(randomOffset.Int64())).Find(&q).Error; err != nil {
			// error handling here
			return err
		}

		return c.JSON(http.StatusOK, q)
	}
}

type QuotePost struct {
	Quote string `json:"quote" form:"quote" query:"quote"`
}

func CreateQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		q := new(QuotePost)
		if err = c.Bind(q); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quote := model.Quote{
			Quote:     q.Quote,
			CreatedAt: datatypes.Date(time.Now()),
			UpdatedAt: datatypes.Date(time.Now()),
		}
		db.Create(&quote)
		return c.JSON(http.StatusCreated, quote)
	}
}

func UpdateQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusBadRequest, "bad request")
		}
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		q := new(QuotePost)
		if err = c.Bind(q); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quote := model.Quote{
			ID: uint(parsedID),
		}
		db.First(&quote)
		quote.Quote = q.Quote
		quote.UpdatedAt = datatypes.Date(time.Now())
		db.Save(quote)
		return c.JSON(http.StatusOK, quote)
	}
}

func DeleteQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusBadRequest, "bad request")
		}
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		quote := model.Quote{
			ID: uint(parsedID),
		}
		db.Delete(&quote)

		return c.NoContent(http.StatusNoContent)
	}
}
