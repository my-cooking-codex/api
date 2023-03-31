package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/my-cooking-codex/api/config"
	"github.com/my-cooking-codex/api/core"
	"github.com/my-cooking-codex/api/db"
	"github.com/my-cooking-codex/api/routes"
	"gorm.io/gorm"
)

type Validator struct {
	validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// HTTP error handler, to handle unexpected errors
func errorHandler(err error, ctx echo.Context) {
	if e, ok := err.(*echo.HTTPError); ok {
		// normal HTTP error
		ctx.JSON(e.Code, e.Message)
		return
	}
	ctx.Logger().Error(err)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.NoContent(http.StatusNotFound)
	} else {
		ctx.NoContent(http.StatusInternalServerError)
	}
}

func main() {
	// Parse config
	var appConfig config.AppConfig
	if err := appConfig.ParseConfig(); err != nil {
		log.Fatalln(err)
	}
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(appConfig.DataPath, os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	os.MkdirAll(path.Join(appConfig.DataPath, core.RecipeImagesOriginalPath), os.ModePerm)
	// Connect to database
	if err := db.InitDB(appConfig.DB); err != nil {
		log.Fatalln(err)
	}
	// Create & setup server
	e := echo.New()
	e.HTTPErrorHandler = errorHandler
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	corsConfig := middleware.DefaultCORSConfig
	{
		corsConfig.AllowOrigins = appConfig.CORSOrigins
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Validator = &Validator{validator: validator.New()}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("AppConfig", appConfig)
			return next(ctx)
		}
	})
	routes.InitRoutes(e, appConfig)
	if appConfig.StaticPath != nil {
		log.Println("Serving static files from", *appConfig.StaticPath)
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  *appConfig.StaticPath,
			HTML5: true,
		}))
	} else {
		e.GET("/", func(ctx echo.Context) error {
			return ctx.HTML(200, "<h1>API Backend Operational</h1>")
		})
	}
	// Start server
	e.Logger.Fatal(e.Start(appConfig.Bind.AsAddress()))
}
