package handler

import (
	"backend-2/api/cmd/db/model"
	"backend-2/api/cmd/db/plugin"
	redisclient "backend-2/api/cmd/utils/redis"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

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
		tag := c.QueryParam("tag")

		redisKey := fmt.Sprintf("QUOTES-%d-%d-%s", limit, page, tag)
		res := &plugin.Pagination{Limit: limit, Page: page}

		if cached, err := getFromCache(redisKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}

		quotes, err := listQuotes(db, *res, tag)
		if err != nil {
			log.Err(err).Msg("listQuotes error")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if err := setCache(redisKey, quotes, 5*time.Minute); err != nil {
			log.Err(err).Msg("Redis SET error")
		}

		return c.JSON(http.StatusOK, quotes)
	}
}

func GetRandomQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		tag := c.QueryParam("tag")
		redisKey := fmt.Sprintf("QUOTES-0-0-%s", tag)

		quote, err := getRandomQuoteFromCache(redisKey)
		if err == nil {
			return c.JSON(http.StatusOK, []*model.Quote{quote})
		}

		quote, err = getRandomQuoteFromDB(db, tag)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		go cacheAllQuotes(db, tag, redisKey)
		return c.JSON(http.StatusOK, []*model.Quote{quote})
	}
}

func CreateQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		q := new(QuotePost)
		if err := c.Bind(q); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		quote := model.Quote{Quote: q.Quote, Author: q.Author}
		if err := db.Create(&quote).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		redisclient.DelByPattern("QUOTES*")
		return c.JSON(http.StatusCreated, quote)
	}
}

func UpdateQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid id")
		}

		q := new(QuotePost)
		if err := c.Bind(q); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		quote := model.Quote{BaseModel: model.BaseModel{ID: uint(id)}}
		if err := db.First(&quote).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "quote not found")
		}

		if q.Quote != "" {
			quote.Quote = q.Quote
		}
		if q.Author != "" {
			quote.Author = q.Author
		}

		if err := db.Save(&quote).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		redisclient.DelByPattern("QUOTES*")
		return c.JSON(http.StatusOK, quote)
	}
}

func DeleteQuotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid id")
		}

		if err := db.Delete(&model.Quote{BaseModel: model.BaseModel{ID: uint(id)}}).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		redisclient.DelByPattern("QUOTES*")
		return c.NoContent(http.StatusNoContent)
	}
}

type QuotePost struct {
	Quote  string `json:"quote" form:"quote" query:"quote"`
	Author string `json:"author" form:"author" query:"author"`
}

func getFromCache(key string) (*plugin.Pagination, error) {
	val, err := redisclient.Rdb.Get(redisclient.Ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var res plugin.Pagination
	if err := json.Unmarshal([]byte(val), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func setCache(key string, value interface{}, duration time.Duration) error {
	return redisclient.Rdb.Set(redisclient.Ctx, key, value, duration).Err()
}

func listQuotes(db *gorm.DB, pagination plugin.Pagination, tag string) (*plugin.Pagination, error) {
	var quotes []*model.Quote
	query := db.Scopes(plugin.Paginate(quotes, &pagination, db))

	if tag != "" {
		query = query.Where("? = ANY(tags)", tag)
	}

	if err := query.Find(&quotes).Error; err != nil {
		return nil, err
	}

	pagination.Rows = quotes
	return &pagination, nil
}

func getRandomQuoteFromCache(redisKey string) (*model.Quote, error) {
	res := &plugin.Pagination{}
	val, err := redisclient.Rdb.Get(redisclient.Ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(val), res); err != nil {
		return nil, err
	}

	randomRow, err := res.GetRandom()
	if err != nil {
		return nil, err
	}

	randomRowJSON, err := json.Marshal(randomRow)
	if err != nil {
		return nil, err
	}

	var quote model.Quote
	if err := json.Unmarshal(randomRowJSON, &quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

func getRandomQuoteFromDB(db *gorm.DB, tag string) (*model.Quote, error) {
	var count int64
	if err := db.Table("quotes").Count(&count).Error; err != nil {
		return nil, err
	}

	randomOffset, err := rand.Int(rand.Reader, big.NewInt(count))
	if err != nil {
		return nil, err
	}

	var quotes []*model.Quote
	query := db.Limit(1).Offset(int(randomOffset.Int64()))
	if tag != "" {
		query = query.Where("? = ANY(tags)", tag)
	}

	if err := query.Find(&quotes).Error; err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes found")
	}

	return quotes[0], nil
}

func cacheAllQuotes(db *gorm.DB, tag, redisKey string) {
	res := &plugin.Pagination{Limit: 0, Page: 0}
	quotes, err := listQuotes(db, *res, tag)
	if err != nil {
		log.Err(err).Msg("listQuotes error inside cache all quotes")
		return
	}

	if quotes != nil && len(quotes.Rows.([]*model.Quote)) > 0 {
		if err := setCache(redisKey, quotes, 5*time.Minute); err != nil {
			log.Err(err).Msg("Redis SET error inside cache all quotes")
		}
	}
}
