package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (app *application) Compile(source []byte) ([]byte, string, error) {
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
	cmd := exec.CommandContext(ctx, "docker", "run", "--name", name, "--rm", "-v", dir+":/mount", app.DockerImage, "satysfi", "main.saty", "-o", "out.pdf")

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

type Application interface {
	Compile([]byte) ([]byte, string, error)
}

type application struct {
	DockerImage string `json:"dockerimage"`
	WorkDir     string `json:"workdir"`
}

type File struct {
	Name string `json:"name"`
}
type Directory struct {
	Name      string `json:"name"`
	path      string
	ChildDirs []*Directory `json:"childdirs"`
	Children  []*File      `json:"children"`
}

func isDir(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You MUST specify environment variable: PORT")
	}

	config, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	app := application{}
	json.Unmarshal(config, &app)

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

		pdf, _, err := app.Compile(source)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pdfb64 := base64.StdEncoding.EncodeToString(pdf)

		return c.String(http.StatusOK, pdfb64)
	})
	e.GET("/filetree", func(c echo.Context) error {
		if !isDir(app.WorkDir) {
			return c.String(http.StatusInternalServerError, "Working directory is not available")
		}

		// iterate working directory and construct tree
		root := Directory{
			Name:      app.WorkDir,
			path:      app.WorkDir,
			ChildDirs: []*Directory{},
			Children:  []*File{},
		}
		dirStack := make([]*Directory, 0, 1)
		dirStack = append(dirStack, &root)
		lastPrefix := root.Name
		filepath.Walk(app.WorkDir, func(p string, info os.FileInfo, err error) error {
			// this function assume the filepath.Walk search in depth first

			// skip git and itself
			if filepath.Base(p) == ".git" {
				return filepath.SkipDir
			}
			if p == app.WorkDir {
				return nil
			}

			// skip if depth > 2
			if len(dirStack) > 2 {
				return filepath.SkipDir
			}

			// pop stack until prefix matches
			for !strings.HasPrefix(p, lastPrefix) {
				dirStack = dirStack[:len(dirStack)-1]
				lastPrefix = dirStack[len(dirStack)-1].path
			}

			if info.IsDir() {
				current := &Directory{
					Name:      filepath.Base(p),
					path:      p,
					ChildDirs: []*Directory{},
					Children:  []*File{},
				}

				// push new directory to stack and update prefix
				dirStack[len(dirStack)-1].ChildDirs = append(dirStack[len(dirStack)-1].ChildDirs, current)
				dirStack = append(dirStack, current)
				lastPrefix = p
			} else if info.Mode().IsRegular() {
				dirStack[len(dirStack)-1].Children = append(dirStack[len(dirStack)-1].Children, &File{
					Name: filepath.Base(p),
				})
			}
			return nil
		})
		return c.JSON(http.StatusOK, root)
	})
	// run
	e.Logger.Fatal(e.Start(":" + port))
}
