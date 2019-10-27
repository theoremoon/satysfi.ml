package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You MUST specify environment variable: PORT")
	}

	e := echo.New()
	e.File("/", "web/index.html")
	e.File("/hello", "web/index.html")
	e.Logger.Fatal(e.Start(":" + port))
}
