package routes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/my-cooking-codex/api/core"
	"github.com/my-cooking-codex/api/db/crud"
	"github.com/my-cooking-codex/api/db/types"
)

func postCreatePantryLocation(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	var formData types.CreatePantryLocation
	if err := core.BindAndValidate(ctx, &formData); err != nil {
		return err
	}

	if pantryLocation, err := crud.CreatePantryLocation(
		formData,
		authenticatedUser.UserID,
	); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusCreated, pantryLocation)
	}
}

func getPantryLocations(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	if pantryLocations, err := crud.GetPantryLocationsByUserID(
		authenticatedUser.UserID,
	); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, pantryLocations)
	}
}

func getPantryLocationByID(ctx echo.Context) error {
	pantryLocationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryLocation(
		authenticatedUser.UserID,
		pantryLocationID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	if pantry, err := crud.GetPantryLocationByID(pantryLocationID); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, pantry)
	}
}

func patchPantryLocationByID(ctx echo.Context) error {
	pantryLocationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryLocation(
		authenticatedUser.UserID,
		pantryLocationID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	var formData core.SelectedUpdate[types.UpdatePantryLocation]
	if err := core.BindAndValidate(ctx, &formData); err != nil {
		return err
	}

	if err := crud.UpdatePantryLocation(pantryLocationID, formData); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func deletePantryLocationByID(ctx echo.Context) error {
	pantryLocationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryLocation(
		authenticatedUser.UserID,
		pantryLocationID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	if err := crud.DeletePantryLocation(pantryLocationID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func postCreatePantryItem(ctx echo.Context) error {
	pantryLocationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryLocation(
		authenticatedUser.UserID,
		pantryLocationID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	var formData types.CreatePantryItem
	if err := core.BindAndValidate(ctx, &formData); err != nil {
		return err
	}

	pantryItem, err := crud.CreatePantryItem(formData, pantryLocationID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, pantryItem)
}

func getPantryItemByID(ctx echo.Context) error {
	pantryItemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryItem(
		authenticatedUser.UserID,
		pantryItemID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	if pantryItem, err := crud.GetPantryItemByID(pantryItemID); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, pantryItem)
	}
}

func getPantryItems(ctx echo.Context) error {
	authenticatedUser := getAuthenticatedUser(ctx)

	var filterParams core.PantryItemsFilterParams
	if err := core.BindAndValidate(ctx, &filterParams); err != nil {
		return err
	}

	// convert human page number into database offset
	rowOffset := (filterParams.Page - 1) * filterParams.PerPage

	if items, err := crud.GetPantryItemsByUserID(
		authenticatedUser.UserID,
		rowOffset,
		filterParams.PerPage,
		crud.PantryItemsFilters{
			Name:       filterParams.Name,
			Labels:     filterParams.Labels,
			LocationId: filterParams.LocationId,
			Expired:    filterParams.Expired,
		},
	); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, items)
	}
}

func patchPantryItemByID(ctx echo.Context) error {
	pantryItemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryItem(
		authenticatedUser.UserID,
		pantryItemID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	var formData core.SelectedUpdate[types.UpdatePantryItem]
	if err := core.BindAndValidate(ctx, &formData); err != nil {
		return err
	}

	// TODO check "Fields" has "LocationID" instead
	if formData.Model.LocationId != (uuid.UUID{}) {
		if isOwner, err := crud.DoesUserOwnPantryLocation(
			authenticatedUser.UserID,
			formData.Model.LocationId,
		); err != nil {
			return err
		} else if !isOwner {
			return ctx.JSON(http.StatusBadRequest, "locationId not found, are you the owner?")
		}
	}

	if err := crud.UpdatePantryItem(pantryItemID, formData); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func deletePantryItemByID(ctx echo.Context) error {
	pantryItemID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return err
	}
	authenticatedUser := getAuthenticatedUser(ctx)

	if isOwner, err := crud.DoesUserOwnPantryItem(
		authenticatedUser.UserID,
		pantryItemID,
	); err != nil {
		return err
	} else if !isOwner {
		return ctx.NoContent(http.StatusNotFound)
	}

	if err := crud.DeletePantryItem(pantryItemID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
