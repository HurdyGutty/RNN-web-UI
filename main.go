package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/HurdyGutty/RNN-web-UI/pkg/read"
	save "github.com/HurdyGutty/RNN-web-UI/pkg/saveAlign"

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

func main() {
	e := echo.New()
	e.Renderer = newTemplate()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "static",
		Filesystem: http.FS(static),
	}))
	e.Use(middleware.Logger())

	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	fmt.Printf("%s\n", reqBody)
	// 	fmt.Printf("%s\n", resBody)
	// }))

	pages := read.ParseData()
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
			message = "Data error"
			return c.Render(500, "save", message)
		}

		page.Values["Align"] = align

		err = save.SaveAlign("internal/DB/aligned_vie-eng.txt", page.Page-1, align)

		if err != nil {
			message = "Read File Error"
			return c.Render(500, "save", message)
		}

		return c.Render(http.StatusOK, "save", message)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
