package routes

import (
	"path"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/config"
)

func getRecipeImageContent(ctx echo.Context) error {
	appConfig := ctx.Get("AppConfig").(config.AppConfig)
	imageID := ctx.Param("id")

	return ctx.File(path.Join(
		appConfig.Data.RecipeOriginalsPath(),
		uuid.MustParse(imageID).String()+".jpg"),
	)
}
