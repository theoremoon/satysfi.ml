package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You MUST specify environment variable: PORT")
	}

	e := echo.New()

	// middlewares
	e.Use(middleware.Logger())

	// handlers
	e.Static("/", "web/dist")
	e.POST("/compile", func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, string(body))
	})
	// run
	e.Logger.Fatal(e.Start(":" + port))
}
