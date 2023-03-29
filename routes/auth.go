package routes

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/config"
	"github.com/my-cooking-codex/api/core"
	"github.com/my-cooking-codex/api/db/crud"
	"gorm.io/gorm"
)

func postLogin(ctx echo.Context) error {
	appConfig := ctx.Get("AppConfig").(config.AppConfig)
	var loginData core.CreateLogin
	if err := core.BindAndValidate(ctx, &loginData); err != nil {
		return err
	}

	// validate username & password
	user, err := crud.GetUserByUsername(loginData.Username)
	if err != nil || !user.IsPasswordMatch(loginData.Password) {
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			// handle when user does not exist
			return ctx.NoContent(http.StatusUnauthorized)
		}
		// fallback, handle error in global error handler
		return err
	}

	authenticationData := core.AuthenticatedUser{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  false,
	}

	// user is valid, create a token
	if token, err := core.CreateAuthenticationToken(
		authenticationData,
		[]byte(appConfig.SecretKey),
	); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, token)
	}
}
