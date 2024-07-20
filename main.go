package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"backend-2/api/cmd/config"
	"backend-2/api/cmd/db"
	"backend-2/api/cmd/db/model"
	"backend-2/api/cmd/handler"
	"backend-2/api/cmd/utils"
	redisclient "backend-2/api/cmd/utils/redis"
	"backend-2/api/graphql"

	_ "backend-2/api/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Quotes API Go
// @version 1.0
// @description This is a quotes API server.

// @contact.name API Support
// @contact.url https://lr-link.vercel.app
// @contact.email aarontanhar2000@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:1323
// @BasePath /
// @schemes http
func main() {
	utils.LoadEnv()

	e := echo.New()
	db, err := db.NewDB()
	logFatal(err)
	
	redisclient.SetClient()

	db.AutoMigrate(&model.User{}, &model.Quote{})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", handler.HelloWorld())
	e.GET("/quotes", handler.GetQuotes(db))
	e.GET("/random", handler.GetRandomQuotes(db))
	e.POST("/quotes", handler.CreateQuotes(db), config.GetJwtMiddleware())
	e.PATCH("/quotes/:id", handler.UpdateQuotes(db), config.GetJwtMiddleware())
	e.DELETE("/quotes/:id", handler.DeleteQuotes(db), config.GetJwtMiddleware())
	e.GET("/meeting/get-token", handler.GetToken())
	e.POST("/meeting/create", handler.CreateMeeting())
	e.POST("/meeting/validate/:id", handler.ValidateMeeting())
	e.GET("/notes", handler.GetNotes(db))
	e.POST("/notes", handler.CreateNotes(db), config.GetJwtMiddleware())
	e.PATCH("/notes/:id", handler.UpdateNotes(db), config.GetJwtMiddleware())
	e.DELETE("/notes/:id", handler.DeleteNotes(db), config.GetJwtMiddleware())

	e.POST("/login", handler.Login(db))

	h, err := graphql.NewHandler(db)
	logFatal(err)
	e.POST("/graphql", echo.WrapHandler(h))

	argsPort := flag.Int("port", -1, "port number")
	flag.Parse()

	port := utils.GetEnv("PORT")
	if port == "" {
		port = "1323"
	}

	if *argsPort > 0 {
		port = strconv.Itoa(*argsPort)
	}

	go func() {
		if err := e.Start(utils.ConcatStr(":", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
