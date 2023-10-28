package routes

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/my-cooking-codex/api/config"
	"github.com/my-cooking-codex/api/core"
)

const (
	AuthenticatedUserKey = "AuthenticatedUser"
	UserTokenKey         = "UserToken"
)

func authenticatedUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authenticatedUser, err := core.GetAuthenticatedUserFromContext(ctx)
		if err != nil {
			// invalid token contents
			return ctx.NoContent(http.StatusUnauthorized)
		}
		// TODO validate username & userID match in database
		ctx.Set(AuthenticatedUserKey, authenticatedUser)
		return next(ctx)
	}
}

func getAuthenticatedUser(ctx echo.Context) core.AuthenticatedUser {
	return ctx.Get(AuthenticatedUserKey).(core.AuthenticatedUser)
}

func InitRoutes(e *echo.Echo, appConfig config.AppConfig) {
	e.GET("/api/info/", getServerInfo)
	e.POST("/api/users/", postCreateUser)
	e.POST("/api/login/", postLogin)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(core.JWTClaims)
		},
		SigningKey: []byte(appConfig.JWTSecret),
		ContextKey: UserTokenKey,
	}
	jwtMiddleware := echojwt.WithConfig(config)

	apiRoutes := e.Group("/api/", jwtMiddleware, authenticatedUserMiddleware)
	{
		apiRoutes.GET("users/me/", getUserMe)
		apiRoutes.GET("labels/", getLabels)
		apiRoutes.POST("recipes/", postCreateRecipe)
		apiRoutes.GET("recipes/", getRecipes)
		apiRoutes.GET("recipes/:id/", getRecipe)
		apiRoutes.PATCH("recipes/:id/", patchRecipe)
		apiRoutes.DELETE("recipes/:id/", deleteRecipe)
		apiRoutes.POST("recipes/:id/image/", postSetRecipeImage, middleware.BodyLimit(appConfig.ImageUploadSizeLimit))
		apiRoutes.DELETE("recipes/:id/image/", deleteRecipeImage)
		apiRoutes.GET("pantry/", getPantryLocations)
		apiRoutes.POST("pantry/", postCreatePantryLocation)
		apiRoutes.GET("pantry/:id/", getPantryLocationByID)
		apiRoutes.PATCH("pantry/:id/", patchPantryLocationByID)
		apiRoutes.DELETE("pantry/:id/", deletePantryLocationByID)
		apiRoutes.POST("pantry/:id/items/", postCreatePantryItem)
		apiRoutes.GET("pantry-items/", getPantryItems)
		apiRoutes.GET("pantry-items/:id/", getPantryItemByID)
		apiRoutes.PATCH("pantry-items/:id/", patchPantryItemByID)
		apiRoutes.DELETE("pantry-items/:id/", deletePantryItemByID)
		apiRoutes.GET("stats/me/", getAccountStats)
	}

	mediaRoutes := e.Group("/media/")
	{
		mediaRoutes.GET("recipe-image/:id", getRecipeImageContent)
	}
}
