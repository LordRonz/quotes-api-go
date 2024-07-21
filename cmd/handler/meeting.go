package handler

import (
	"backend-2/api/cmd/utils"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type MeetingToken struct {
	Token string `json:"token"`
}

var VIDEOSDK_API_ENDPOINT = "https://api.videosdk.live"

func getVideoSDKAPIKey() string {
	return utils.GetEnv("VIDEOSDK_API_KEY", "")
}

func getVideoSDKSecretKey() string {
	return utils.GetEnv("VIDEOSDK_SECRET_KEY", "")
}

func GetToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		var permissions [2]string
		permissions[0] = "allow_join"
		permissions[1] = "allow_mod"

		atClaims := jwt.MapClaims{}
		atClaims["apikey"] = getVideoSDKAPIKey()
		atClaims["permissions"] = permissions
		atClaims["iat"] = time.Now().Unix()
		atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		token, err := at.SignedString([]byte(getVideoSDKSecretKey()))
		if err != nil {
			log.Err(err).Msg("")
		}

		return c.JSON(http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
	}
}

func CreateMeeting() echo.HandlerFunc {
	return func(c echo.Context) error {
		m := new(MeetingToken)
		if err := c.Bind(m); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		url := VIDEOSDK_API_ENDPOINT + "/v2/rooms"
		method := "POST"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		req.Header.Add("Authorization", m.Token)
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		return c.JSON(http.StatusOK, result)
	}
}

func ValidateMeeting() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		m := new(MeetingToken)
		if err := c.Bind(m); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		url := VIDEOSDK_API_ENDPOINT + "/v2/rooms/validate/" + id
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		req.Header.Add("Authorization", m.Token)
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Err(err).Msg("")
			return err
		}
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		return c.JSON(http.StatusOK, result)
	}
}
