package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/db/crud"
)

type accountStats struct {
	UserCount   int64 `json:"userCount"`
	RecipeCount int64 `json:"recipeCount"`
}

func getAccountStats(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	recipeCount, err := crud.GetRecipesByUserIDCount(authenticatedUser.UserID)
	if err != nil {
		return err
	}
	userCount, err := crud.GetUserCount()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, accountStats{
		UserCount:   userCount,
		RecipeCount: recipeCount,
	})
}
