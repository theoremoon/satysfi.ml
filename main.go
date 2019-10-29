package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/theoremoon/SATySFi-Online/repository"
	"github.com/theoremoon/SATySFi-Online/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You MUST specify the environment variable: PORT")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("You MUST specify the environment variable: DSN")
	}

	repo, err := repository.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	service := service.New(repo, 1)

	e := echo.New()

	// middleware
	e.Use(middleware.Logger())

	// handle for SPA
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "./ui/dist/",
		Index:  "index.html",
		HTML5:  true,
		Browse: false,
	}))

	// API
	e.POST("/api/register", func(c echo.Context) error {
		request := new(struct {
			Username string `json:"username"`
			Password string `json:"password"`
		})
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		user, err := service.Register(request.Username, request.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		session, err := service.CreateSession(user.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		cookie := new(http.Cookie)
		cookie.Name = "session"
		cookie.Value = strconv.Itoa(session.ID)
		cookie.Expires = time.Unix(session.ExpiredAt, 0)
		c.SetCookie(cookie)
		return c.String(http.StatusOK, "")
	})

	e.Logger.Fatal(e.Start(":" + port))
}
