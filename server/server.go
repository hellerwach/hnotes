package server

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"hellerwach.com/go/hnotes/note"
)

// Template implements the echo.Renderer interface.
type Template struct {
	templates *template.Template
}

// Render renders a template with data and sends a text/html response with a
// status code.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		logrus.Warnf("Could not render template %q: %v\n", name, err)
		return err
	}
	return nil
}

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".config/hnotes")

// Run runs the note server and makes further function calls. It will also
// terminate the process.
func Run(port int) {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	// Templates
	funcs := map[string]any{
		"filepathJoin": filepath.Join,
		"hasPrefix":    strings.HasPrefix,
		"hasSuffix":    strings.HasSuffix,
	}
	t := &Template{
		templates: template.Must(template.New("").Funcs(funcs).ParseGlob(filepath.Join(ConfigPath, "templates/*.html"))),
	}
	e.Renderer = t

	// Routes
	e.GET("/*", func(c echo.Context) error {
		path := strings.TrimPrefix(c.Request().URL.Path, "/")
		// Discriminate between Markdown, non-Markdown and directories
		if path == "" {
			return directoryView(c)
		}

		info, err := os.Stat(path)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if strings.HasSuffix(path, ".md") {
			return serveNote(c)
		} else if info.IsDir() {
			return directoryView(c)
		}

		return c.File(path)
	})

	// Start server
	go func() {
		if err := e.Start(":" + strconv.Itoa(port)); err != nil && err != http.ErrServerClosed {
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

func serveNote(c echo.Context) error {
	path := strings.TrimPrefix(c.Request().URL.Path, "/")
	note, err := note.New(path)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could create new note object: "+err.Error())
	}
	return c.Render(http.StatusOK, "single", *note)
}

type dir struct {
	Path    string
	Entries []fs.DirEntry
}

// directoryView renders the dir template with the data of the current
// directory.
func directoryView(c echo.Context) error {
	path := strings.TrimPrefix(c.Request().URL.Path, "/")
	if path == "" {
		path = "."
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not read directory: "+err.Error())
	}
	dir := dir{Entries: dirEntries, Path: path}

	return c.Render(http.StatusOK, "dir", dir)
}
