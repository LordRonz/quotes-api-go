package handler

import (
	"backend-2/api/cmd/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type MeetingToken struct {
	Token string `json:"token"`
}
type tokenResponse struct {
	token string
}

var VIDEOSDK_API_KEY = utils.GetEnv("VIDEOSDK_API_KEY", "")
var VIDEOSDK_SECRET_KEY = utils.GetEnv("VIDEOSDK_SECRET_KEY", "")
var VIDEOSDK_API_ENDPOINT = "https://api.videosdk.live"

func GetToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		var permissions [2]string
		permissions[0] = "allow_join"
		permissions[1] = "allow_mod"

		atClaims := jwt.MapClaims{}
		atClaims["apikey"] = VIDEOSDK_API_KEY
		atClaims["permissions"] = permissions
		atClaims["iat"] = time.Now().Unix()
		atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		token, err := at.SignedString([]byte(VIDEOSDK_SECRET_KEY))

		if err != nil {
			fmt.Printf("%v\n", err)
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

		str, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := VIDEOSDK_API_ENDPOINT + "/api/meetings"
		method := "POST"

		// fmt.Print("\n", strings.NewReader(s))
		client := &http.Client{}
		req, err := http.NewRequest(method, url, bytes.NewBuffer(str))
		if err != nil {
			fmt.Println(err)
			return err
		}
		req.Header.Add("Authorization", m.Token)
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, body)
	}
}

func ValidateMeeting() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		m := new(MeetingToken)
		if err := c.Bind(m); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		str, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := VIDEOSDK_API_ENDPOINT + "/api/meetings/" + id
		method := "POST"

		// fmt.Print("\n", strings.NewReader(s))
		client := &http.Client{}
		req, err := http.NewRequest(method, url, bytes.NewBuffer(str))

		if err != nil {
			fmt.Println(err)
			return err
		}
		req.Header.Add("Authorization", m.Token)
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, body)
	}
}
