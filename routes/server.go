package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/core"
)

type serverInfo struct {
	APIVersionMajor     uint `json:"apiVersionMajor"`
	RegistrationAllowed bool `json:"registrationAllowed"`
}

func getServerInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, serverInfo{
		APIVersionMajor:     core.APIVersionMajor,
		RegistrationAllowed: true,
	})
}
