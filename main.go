package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate npm run build

//go:embed static
var static embed.FS

func main() {
	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "static",
		Browse:     true,
		IgnoreBase: true,
		HTML5:      true,
		Filesystem: http.FS(static),
	}))

	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8080"))
}
