package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func compile(source []byte) ([]byte, string, error) {
	dir, err := ioutil.TempDir("", "satysfibuild")
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(dir)

	sourceFile := filepath.Join(dir, "main.saty")
	if err := ioutil.WriteFile(sourceFile, source, 06666); err != nil {
		return nil, "", err
	}

	randomBuf := make([]byte, 8)
	_, err = rand.Read(randomBuf)
	if err != nil {
		return nil, "", err
	}
	name := fmt.Sprintf("%x", randomBuf)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	cmd := exec.CommandContext(ctx, "docker", "run", "--name", name, "--rm", "-v", dir+":/mount", "satysfi", "satysfi", "main.saty", "-o", "out.pdf")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-ctx.Done():
		exec.Command("docker", "kill", name).Run()
	case <-done:
		// do nothing
	}

	log.Println(stdout.String())
	log.Println(stderr.String())

	pdf, err := ioutil.ReadFile(filepath.Join(dir, "out.pdf"))
	if err != nil {
		return nil, "", err
	}

	return pdf, stdout.String(), nil
}

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
		source, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		pdf, _, err := compile(source)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pdfb64 := base64.StdEncoding.EncodeToString(pdf)

		return c.String(http.StatusOK, pdfb64)
	})
	// run
	e.Logger.Fatal(e.Start(":" + port))
}
