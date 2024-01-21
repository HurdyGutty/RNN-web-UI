package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

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
		templates: template.Must(template.New("index").Funcs(template.FuncMap{
			"add":   func(a, b int) int { return a + b },
			"minus": func(a, b int) int { return a - b },
		}).ParseGlob("views/*.html")),
	}
}

type Dict map[string]interface{}
type AlignmentValues Dict

func newDict() AlignmentValues {
	return AlignmentValues{
		"Nom":   Dict{"Key": "Nom", "Data": []string{}},
		"Eng":   Dict{"Key": "Eng", "Data": []string{}},
		"Align": [][]int{},
	}
}

func (values AlignmentValues) mockData(nom, eng []string, align [][]int) AlignmentValues {
	values["Nom"] = Dict{"Key": "Nom", "Data": nom}
	values["Eng"] = Dict{"Key": "Eng", "Data": eng}
	values["Align"] = align
	return values
}

type Page struct {
	Page      int
	TotalPage int
	Values    AlignmentValues
}

func newPage(page, totalPage int, values AlignmentValues) Page {
	return Page{
		Page:      page,
		TotalPage: totalPage,
		Values:    values,
	}
}

type Pages []Page

func newPages() Pages {
	newPages := []Page{}
	return newPages
}

func mockPages() Pages {
	newPages := []Page{
		newPage(1, 3, newDict().mockData(
			[]string{"Je", "parle", "francais"},
			[]string{"I", "speak", "French"},
			[][]int{{0, 0}, {1, 1}, {2, 2}})),
		newPage(2, 3, newDict().mockData(
			[]string{"Battre", "le", "fer", "pendant", "qu'il", "est", "chaud"},
			[]string{"Strike", "the", "iron", "while", "it", "is", "hot"},
			[][]int{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}})),
		newPage(3, 3, newDict().mockData(
			[]string{"En", "faire", "tout", "un", "fromage"},
			[]string{"To", "make", "a", "whole", "cheese"},
			[][]int{{0, 0}, {1, 1}, {2, 3}, {3, 2}, {4, 4}})),
	}
	return newPages
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "static",
		Filesystem: http.FS(static),
	}))
	// e.Use(middleware.Logger())

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Printf("%s\n", reqBody)
		fmt.Printf("%s\n", resBody)
	}))

	pages := mockPages()
	page := pages[0]

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", page)
	})
	e.GET("/page/:pageNum", func(c echo.Context) error {
		pageNum, err := strconv.Atoi(c.Param("pageNum"))
		if err != nil {
			return c.Render(http.StatusInternalServerError, "500", nil)
		}
		if pageNum < 1 || pageNum > len(pages) {
			return c.Render(http.StatusNotFound, "404", nil)
		}
		page = pages[pageNum-1]
		return c.Render(http.StatusOK, "index", page)
	})
	e.PUT("/save", func(c echo.Context) error {
		data_return := c.FormValue("align")
		align := [][]int{}
		err := json.Unmarshal([]byte(data_return), &align)
		message := "Saved"

		if err != nil {
			message = "Error"
			return c.Render(500, "save", message)
		}

		page.Values["Align"] = align

		return c.Render(http.StatusOK, "save", message)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
