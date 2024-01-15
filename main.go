package main

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate npm run build

//go:embed static
var static embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Dict map[string]interface{}

func newDict() Dict {
	return Dict{
		"Nom": []string{},
		"Eng": []string{},
	}
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "static",
		Filesystem: http.FS(static),
	}))

	e.Use(middleware.Logger())

	data := newDict()
	nom := []string{"a", "b", "c"}
	eng := []string{"1", "2", "3"}

	data["Nom"] = nom
	data["Eng"] = eng

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
