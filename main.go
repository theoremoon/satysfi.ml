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
	"github.com/otiai10/copy"
)

func randomName() string {
	randomBuf := make([]byte, 8)
	_, err := rand.Read(randomBuf)
	if err != nil {
		panic(err)
	}
	name := fmt.Sprintf("%x", randomBuf)
	return name
}

func (app *application) Compile(main string) ([]byte, string, error) {
	dir := filepath.Join(os.TempDir(), "satysfibuild"+randomName())
	// defer os.RemoveAll(dir)

	copy.Copy(app.WorkDir, dir)

	name := randomName()
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	cmd := exec.CommandContext(ctx, "docker", "run", "--name", name, "--rm", "-v", dir+":/mount", app.DockerImage, "satysfi", filepath.Join("/mount", main), "-o", "out.pdf")

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
	Path string `json:"path"`
}
type Directory struct {
	Name      string       `json:"name"`
	Path      string       `json:"path"`
	ChildDirs []*Directory `json:"childdirs"`
	Children  []*File      `json:"children"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func isDir(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}
func isFile(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}
func verifyPath(path string) bool {
	if strings.Contains(path, "../") {
		return false
	}
	if strings.HasPrefix(path, ".git") {
		return false
	}
	return true
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
	e.File("/*", "web/dist/index.html")
	e.POST("/api/save", func(c echo.Context) error {
		request := new(struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		})
		if err := c.Bind(request); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if !verifyPath(request.Path) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"Bad path"})
		}
		path := filepath.Join(app.WorkDir, request.Path)
		if !isFile(path) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"Bad path"})
		}

		if err := ioutil.WriteFile(path, []byte(request.Content), 0); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		}
		return c.String(http.StatusOK, "")
	})
	e.POST("/api/compile", func(c echo.Context) error {
		path := new(struct {
			Path string `json:"path"`
		})
		if err := c.Bind(path); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if !verifyPath(path.Path) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"Bad path"})
		}

		pdf, _, err := app.Compile(strings.TrimPrefix(path.Path, app.WorkDir))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pdfb64 := base64.StdEncoding.EncodeToString(pdf)

		return c.String(http.StatusOK, pdfb64)
	})
	e.GET("/api/getfile", func(c echo.Context) error {
		filename := c.Request().URL.Query().Get("filename")
		if !verifyPath(filename) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"Bad path"})
		}
		path := filepath.Join(app.WorkDir, filename)
		if !isFile(path) {
			return c.JSON(http.StatusNotFound, ErrorResponse{"Not Found"})
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		}
		contentType := http.DetectContentType(content)
		if !strings.HasPrefix(contentType, "text/") {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"Not a Text File"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":    filepath.Base(filename),
			"path":    strings.TrimPrefix(path, app.WorkDir),
			"content": string(content),
		})
	})
	e.GET("/api/filetree", func(c echo.Context) error {
		if !isDir(app.WorkDir) {
			return c.String(http.StatusInternalServerError, "Working directory is not available")
		}

		// iterate working directory and construct tree
		root := Directory{
			Name:      app.WorkDir,
			Path:      app.WorkDir,
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
				lastPrefix = dirStack[len(dirStack)-1].Path
			}

			if info.IsDir() {
				current := &Directory{
					Name:      filepath.Base(p),
					Path:      p,
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
					Path: strings.TrimPrefix(p, app.WorkDir),
				})
			}
			return nil
		})
		return c.JSON(http.StatusOK, root)
	})
	// run
	e.Logger.Fatal(e.Start(":" + port))
}
