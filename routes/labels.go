package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/db/crud"
)

func getLabels(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	labels, err := crud.GetLabelNamesByUser(authenticatedUser.UserID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, labels)
}
