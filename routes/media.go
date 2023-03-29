package routes

import (
	"path"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/config"
	"github.com/my-cooking-codex/api/core"
)

func getRecipeImageContent(ctx echo.Context) error {
	appConfig := ctx.Get("AppConfig").(config.AppConfig)
	imageID := ctx.Param("id")

	return ctx.File(path.Join(appConfig.DataPath, core.RecipeImagesOriginalPath, uuid.MustParse(imageID).String()+".jpg"))
}
