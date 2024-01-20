package main

import (
	"embed"
	"encoding/json"
	"fmt"
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
		"Nom":   Dict{"Key": "Nom", "Data": []string{}},
		"Eng":   Dict{"Key": "Eng", "Data": []string{}},
		"Align": [][]int{},
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

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Printf("%s\n", reqBody)
		fmt.Printf("%s\n", resBody)
	}))

	data := newDict()
	nom := []string{"a", "b", "c", "d"}
	eng := []string{"1", "2", "3"}

	align := [][]int{{0, 1}, {1, 2}, {2, 0}, {3, 0}}

	data["Nom"] = Dict{"Key": "Nom", "Data": nom}
	data["Eng"] = Dict{"Key": "Eng", "Data": eng}
	data["Align"] = align

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})
	e.PUT("/save", func(c echo.Context) error {
		data_return := c.FormValue("align")
		fmt.Printf("%s\n", data_return)
		align := [][]int{}
		err := json.Unmarshal([]byte(data_return), &align)
		message := "Saved"

		if err != nil {
			fmt.Printf("%s\n", err)
			message = "Error"
			return c.Render(500, "save", message)
		}

		data["Align"] = align
		fmt.Printf("%v\n", data)

		return c.Render(http.StatusOK, "save", message)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
