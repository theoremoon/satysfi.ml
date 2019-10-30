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

func Compile(dockerImage, buildDir, workDir, path string) ([]byte, string, error) {
	dir := filepath.Join(buildDir, "satysfibuild"+randomName())
	defer os.RemoveAll(dir)

	copy.Copy(workDir, dir)

	name := randomName()
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	cmd := exec.CommandContext(ctx, "docker", "run", "--name", name, "--rm", "-v", dir+":/mount", dockerImage, "satysfi", filepath.Join("/mount", path), "-o", "out.pdf")

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

type application struct {
	DockerImage string `json:"dockerimage"`
	WorkDir     string `json:"workdir"`
	BuildDir    string `json:"builddir"`
	TemplateDir string `json:"templatedir"`
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

func travarseDirectory(root string) (*Directory, error) {
	tree := Directory{
		Name:      root,
		Path:      root,
		ChildDirs: []*Directory{},
		Children:  []*File{},
	}
	dirStack := make([]*Directory, 0, 1)
	dirStack = append(dirStack, &tree)
	lastPrefix := tree.Path
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		// this function assume the filepath.Walk search in depth first

		// skip git and itself
		if filepath.Base(p) == ".git" {
			return filepath.SkipDir
		}
		if p == root {
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
				Path:      strings.TrimPrefix(p, root),
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
				Path: strings.TrimPrefix(p, root),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tree, nil
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

func verifyID(id string) bool {
	return strings.Trim(id, "0123456789abcdef") == ""
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
	app.BuildDir, err = filepath.Abs(app.BuildDir)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()

	// middlewares
	e.Use(middleware.Logger())

	// handlers
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "./ui/dist/",
		Index:  "index.html",
		HTML5:  true,
		Browse: false,
	}))

	// get project file tree
	e.GET("/api/:id/list", func(c echo.Context) error {
		id := c.Param("id")
		path := filepath.Join(app.WorkDir, id)
		if !verifyID(id) || !isDir(path) {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}

		tree, err := travarseDirectory(path)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		tree.Name = "/"
		tree.Path = "/"
		return c.JSON(http.StatusOK, tree)
	})

	// get project file content
	e.GET("/api/:id/get", func(c echo.Context) error {
		id := c.Param("id")
		path := c.Request().URL.Query().Get("path")
		if !verifyID(id) || !verifyPath(path) {
			return c.String(http.StatusNotFound, "")
		}

		path = filepath.Join(app.WorkDir, id, path)
		if !isFile(path) {
			return c.String(http.StatusNotFound, "Not Found")
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		contentType := http.DetectContentType(content)
		if !strings.HasPrefix(contentType, "text/") {
			return c.String(http.StatusBadRequest, "Not a Text File")
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":    filepath.Base(path),
			"path":    strings.TrimPrefix(path, filepath.Join(app.WorkDir, id)),
			"content": string(content),
		})
	})

	// create new project
	e.POST("/api/new-project", func(c echo.Context) error {
		id := randomName()
		path := filepath.Join(app.WorkDir, id)
		for isDir(path) {
			id = randomName()
			path = filepath.Join(app.WorkDir, id)
		}
		err := copy.Copy(app.TemplateDir, path)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, id)
	})

	// save or create new file in project
	e.POST("/api/:id/save", func(c echo.Context) error {
		req := new(struct {
			Path string `json:"path"`
			Data string `json:"data"`
		})
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		id := c.Param("id")
		path := filepath.Join(app.WorkDir, id)

		if !verifyID(id) || !isDir(path) {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		if !verifyPath(req.Path) {
			return c.String(http.StatusBadRequest, "Invalid path")
		}

		content, err := base64.StdEncoding.DecodeString(req.Data)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		path = filepath.Join(path, req.Path)
		dir := filepath.Dir(path)
		if !isDir(dir) {
			if err := os.MkdirAll(dir, 0775); err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}

		err = ioutil.WriteFile(path, content, 0640)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "")
	})

	e.POST("/api/:id/compile", func(c echo.Context) error {
		req := new(struct {
			Path string `json:"path"`
		})
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		id := c.Param("id")

		if !verifyID(id) {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}

		if !verifyPath(req.Path) {
			return c.String(http.StatusBadRequest, "Bad path")
		}

		dir := filepath.Join(app.WorkDir, id)
		pdf, _, err := Compile(app.DockerImage, app.BuildDir, dir, req.Path)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pdfb64 := base64.StdEncoding.EncodeToString(pdf)

		return c.String(http.StatusOK, pdfb64)
	})

	// run
	e.Logger.Fatal(e.Start(":" + port))
}
