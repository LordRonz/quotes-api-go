package handler

import (
	"backend-2/api/cmd/db/model"
	"backend-2/api/cmd/db/plugin"
	redisclient "backend-2/api/cmd/utils/redis"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"net/http"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// Hello World
// @Summary Show the status of server.
// @Description get the hello world from server.
// @Tags Hello World
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HelloWorld() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{
			Message: "Hello, World!",
		})
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

		redisKey := fmt.Sprintf("QUOTES-%d-%d", limit, page)

		var res *plugin.Pagination = &plugin.Pagination{}

		val, err := redisclient.Rdb.Get(redisclient.Ctx, redisKey).Result()

		if err == redis.Nil {
			res, err = listQuotes(db, pagination)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			err = redisclient.Rdb.Set(redisclient.Ctx, redisKey, res, 60*time.Second).Err()

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		} else {
			err = json.Unmarshal([]byte(val), res)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
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
	Quote  string `json:"quote" form:"quote" query:"quote"`
	Author string `json:"author" form:"author" query:"author"`
}

func CreateQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		q := new(QuotePost)
		if err = c.Bind(q); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		quote := model.Quote{
			Quote:     q.Quote,
			Author:    q.Author,
		}
		db.Create(&quote)
		redisclient.DelByPattern("QUOTES*")
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
			BaseModel: model.BaseModel{
				ID: uint(parsedID),
			},
		}
		db.First(&quote)
		if q.Quote != "" {
			quote.Quote = q.Quote
		}
		if q.Author != "" {
			quote.Author = q.Author
		}
		db.Save(quote)
		redisclient.DelByPattern("QUOTES*")
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
			BaseModel: model.BaseModel{
				ID: uint(parsedID),
			},
		}
		db.Delete(&quote)
		redisclient.DelByPattern("QUOTES*")
		return c.NoContent(http.StatusNoContent)
	}
}
