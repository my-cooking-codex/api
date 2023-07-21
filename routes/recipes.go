package routes

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/config"
	"github.com/my-cooking-codex/api/core"
	"github.com/my-cooking-codex/api/db"
	"github.com/my-cooking-codex/api/db/crud"
)

func postCreateRecipe(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	var recipeData db.CreateRecipe
	if err := core.BindAndValidate(ctx, &recipeData); err != nil {
		return err
	}

	recipe, err := crud.CreateRecipe(recipeData, authenticatedUser.UserID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, recipe)
}

func getRecipes(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	var filterParams core.RecipesFilterParams
	if err := core.BindAndValidate(ctx, &filterParams); err != nil {
		return err
	}

	// convert human page number into database offset
	rowOffset := (filterParams.Page - 1) * filterParams.PerPage

	recipes, err := crud.GetRecipesByUserID(
		authenticatedUser.UserID,
		rowOffset,
		filterParams.PerPage,
		crud.RecipesFilterParams{
			Title:         filterParams.Title,
			Labels:        filterParams.Labels,
			Freezable:     filterParams.Freezable,
			MicrowaveOnly: filterParams.MicrowaveOnly,
		},
	)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, recipes)
}

func getRecipe(ctx echo.Context) error {
	recipeID := ctx.Param("id")
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnRecipe(
		authenticatedUser.UserID,
		uuid.MustParse(recipeID),
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	recipe, err := crud.GetRecipeById(uuid.MustParse(recipeID))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, recipe)
}

func patchRecipe(ctx echo.Context) error {
	recipeID := ctx.Param("id")
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnRecipe(
		authenticatedUser.UserID,
		uuid.MustParse(recipeID),
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	var recipeData db.UpdateRecipe
	if err := core.BindAndValidate(ctx, &recipeData); err != nil {
		return err
	}

	if _, err := crud.UpdateRecipe(uuid.MustParse(recipeID), recipeData); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func deleteRecipe(ctx echo.Context) error {
	appConfig := ctx.Get("AppConfig").(config.AppConfig)
	recipeID := ctx.Param("id")
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnRecipe(
		authenticatedUser.UserID,
		uuid.MustParse(recipeID),
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	recipe, err := crud.GetRecipeById(uuid.MustParse(recipeID))
	if err != nil {
		return err
	}

	if err := crud.DeleteRecipe(uuid.MustParse(recipeID)); err != nil {
		return err
	}

	os.Remove(path.Join(
		appConfig.Data.RecipeOriginalsPath(),
		recipe.ImageID.String()+".jpg",
	))

	return ctx.NoContent(http.StatusNoContent)
}

func postSetRecipeImage(ctx echo.Context) error {
	appConfig := ctx.Get("AppConfig").(config.AppConfig)
	recipeID := ctx.Param("id")
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnRecipe(
		authenticatedUser.UserID,
		uuid.MustParse(recipeID),
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	recipe, err := crud.GetRecipeById(uuid.MustParse(recipeID))
	if err != nil {
		return err
	}

	var content = make([]byte, ctx.Request().ContentLength)
	// TODO handle errors
	var b = bytes.Buffer{}
	io.Copy(&b, ctx.Request().Body)
	b.Read(content)
	if optimisedContent, err := core.OptimiseImageToJPEG(content, int(appConfig.OptimizedImageSize)); err == nil {
		content = optimisedContent
	} else {
		return err
	}

	imageID := uuid.New()
	imagePath := path.Join(
		appConfig.Data.RecipeOriginalsPath(),
		imageID.String()+".jpg",
	)
	if err := os.WriteFile(imagePath, content, 0644); err != nil {
		return err
	}

	if err := crud.UpdateRecipeImage(uuid.MustParse(recipeID), &imageID); err != nil {
		return err
	}

	// Remove old image if one was set
	if recipe.ImageID != nil {
		os.Remove(path.Join(
			appConfig.Data.RecipeOriginalsPath(),
			recipe.ImageID.String()+".jpg",
		))
	}

	return ctx.JSON(http.StatusCreated, imageID.String())
}

func deleteRecipeImage(ctx echo.Context) error {
	appConfig := ctx.Get("AppConfig").(config.AppConfig)
	recipeID := ctx.Param("id")
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnRecipe(
		authenticatedUser.UserID,
		uuid.MustParse(recipeID),
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	recipe, err := crud.GetRecipeById(uuid.MustParse(recipeID))
	if err != nil {
		return err
	}

	os.Remove(path.Join(
		appConfig.Data.RecipeOriginalsPath(),
		recipe.ImageID.String()+".jpg",
	))

	if err := crud.UpdateRecipeImage(uuid.MustParse(recipeID), nil); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
