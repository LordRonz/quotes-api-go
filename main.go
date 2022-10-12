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
	"backend-2/api/cmd/handler"
	"backend-2/api/cmd/utils"
	_ "backend-2/api/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	e := echo.New()
	db, err := db.NewDB()
	logFatal(err)

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", handler.HelloWorld())
	e.GET("/quote", handler.GetQuotes(db))
	e.GET("/random", handler.GetRandomQuotes(db))
	e.POST("/quote", handler.CreateQuotes(db), config.GetJwtMiddleware())
	e.PATCH("/quote/:id", handler.UpdateQuotes(db), config.GetJwtMiddleware())
	e.DELETE("/quote/:id", handler.DeleteQuotes(db), config.GetJwtMiddleware())

	e.POST("/login", handler.Login(db))

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
