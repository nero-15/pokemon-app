package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mtslzr/pokeapi-go"
)

const baseURL = "https://pokeapi.co/api/v2/"

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.New("").Delims("[[", "]]").ParseGlob("views/*.html")), // vue.jsとdelimsがかぶるので変更
	}
	e.Renderer = renderer

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}"` + "\n",
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	})

	e.GET("/pokemon/:id", func(c echo.Context) error {
		id := c.Param("id")
		pokemon, err := pokeapi.Pokemon(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		fmt.Println(pokemon.Name)
		return c.JSON(http.StatusOK, pokemon)
	})

	e.GET("/type/:id", func(c echo.Context) error {
		id := c.Param("id")
		pokemontype, err := pokeapi.Type(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, pokemontype)
	})

	e.GET("/api/search", func(c echo.Context) error {
		name := c.QueryParam("name")
		pokemon, err := pokeapi.Pokemon(name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, pokemon)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
