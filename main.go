package main

import (
	"encoding/base64"
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
		pdf, err := ioutil.ReadFile("sample.pdf")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pdfb64 := base64.StdEncoding.EncodeToString(pdf)

		return c.String(http.StatusOK, pdfb64)
	})
	// run
	e.Logger.Fatal(e.Start(":" + port))
}
