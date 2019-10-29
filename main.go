package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You MUST specify the environment variable: PORT")
	}

	e := echo.New()

	// middleware
	e.Use(middleware.Logger())

	// handle
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "./ui/dist/",
		Index:  "index.html",
		HTML5:  true,
		Browse: false,
	}))

	e.Logger.Fatal(e.Start(":" + port))
}
