package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/core"
	"github.com/my-cooking-codex/api/db"
	"github.com/my-cooking-codex/api/db/crud"
)

func postCreateUser(ctx echo.Context) error {
	var userData db.CreateUser
	if err := core.BindAndValidate(ctx, &userData); err != nil {
		return err
	}

	user, err := crud.CreateUser(userData)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, user)
}

func getUserMe(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	user, err := crud.GetRecipeById(authenticatedUser.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}
