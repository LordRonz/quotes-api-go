package handler

import (
	"backend-2/api/cmd/db/model"
	"backend-2/api/cmd/db/plugin"
	redisclient "backend-2/api/cmd/utils/redis"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func GetNotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		page, _ := strconv.Atoi(c.QueryParam("page"))

		pagination := plugin.Pagination{
			Limit: limit,
			Page:  page,
		}

		redisKey := fmt.Sprintf("NOTES-%d-%d", limit, page)

		var res *plugin.Pagination = &plugin.Pagination{}

		val, err := redisclient.Rdb.Get(redisclient.Ctx, redisKey).Result()

		if err == redis.Nil {
			res, err = listNotes(db, pagination)

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

func listNotes(db *gorm.DB, pagination plugin.Pagination) (*plugin.Pagination, error) {
	var notes []*model.Note

	db.Scopes(plugin.Paginate(notes, &pagination, db)).Find(&notes)
	pagination.Rows = notes

	return &pagination, nil
}

type NotePost struct {
	Note        string `json:"note" form:"note" query:"note"`
	Title       string `json:"title" form:"title" query:"title"`
	Description string `json:"description" form:"description" query:"description"`
}

func CreateNotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		n := new(NotePost)
		if err = c.Bind(n); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		note := model.Note{
			Title:       n.Title,
			Description: n.Description,
			Note:        n.Note,
		}
		db.Create(&note)
		redisclient.DelByPattern("NOTES*")
		return c.JSON(http.StatusCreated, note)
	}
}

func UpdateNotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusBadRequest, "bad request")
		}
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		n := new(NotePost)
		if err = c.Bind(n); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		note := model.Note{
			BaseModel: model.BaseModel{
				ID: uint(parsedID),
			},
		}
		db.First(&note)
		if n.Title != "" {
			note.Title = n.Title
		}
		if n.Description != "" {
			note.Description = n.Description
		}
		if n.Note != "" {
			note.Note = n.Note
		}
		db.Save(note)
		redisclient.DelByPattern("NOTES*")
		return c.JSON(http.StatusOK, note)
	}
}

func DeleteNotes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusBadRequest, "bad request")
		}
		parsedID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		note := model.Note{
			BaseModel: model.BaseModel{
				ID: uint(parsedID),
			},
		}
		db.Delete(&note)
		redisclient.DelByPattern("NOTES*")
		return c.NoContent(http.StatusNoContent)
	}
}
