package db

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RecipeIngredient struct {
	Name        string  `json:"name" validate:"required"`
	Amount      float32 `json:"amount" validate:"required"`
	UnitType    string  `json:"unitType" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type RecipeStep struct {
	Title       *string `json:"title,omitempty"`
	Description string  `json:"description" validate:"required"`
}

type RecipeInfoYields struct {
	Value    uint   `json:"value" validate:"required"`
	UnitType string `json:"unitType" validate:"required"`
}

type CreateUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *CreateUser) IntoUser() User {
	user := User{
		Username: u.Username,
	}
	user.SetPassword(u.Password)
	return user
}

type CreateRecipeInfo RecipeInfo

type CreateRecipe struct {
	Title            string             `json:"title" validate:"required"`
	Info             CreateRecipeInfo   `json:"info,omitempty"`
	ShortDescription *string            `json:"shortDescription,omitempty"`
	LongDescription  *string            `json:"longDescription,omitempty"`
	Ingredients      []RecipeIngredient `json:"ingredients,omitempty"`
	Steps            []RecipeStep       `json:"steps,omitempty"`
}

func (r *CreateRecipe) IntoRecipe(ownerID uuid.UUID, imageID *uuid.UUID) Recipe {
	return Recipe{
		OwnerID:          ownerID,
		Title:            r.Title,
		Info:             RecipeInfo(r.Info),
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		ImageID:          imageID,
	}
}

type UpdateIngredient struct {
	Name        string  `json:"name,omitempty"`
	Amount      float32 `json:"amount,omitempty"`
	UnitType    string  `json:"unitType,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateStep struct {
	Title       *string `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
}

type UpdateRecipeInfo RecipeInfo

type UpdateRecipe struct {
	Title            string              `json:"title,omitempty"`
	Info             UpdateRecipeInfo    `json:"info,omitempty"`
	ShortDescription *string             `json:"shortDescription,omitempty"`
	LongDescription  *string             `json:"longDescription,omitempty"`
	Ingredients      *[]UpdateIngredient `json:"ingredients,omitempty"`
	Steps            *[]UpdateStep       `json:"steps,omitempty"`
	ImageID          *uuid.UUID          `json:"-"`
}

func (r *UpdateRecipe) IntoRecipe() Recipe {
	return Recipe{
		Title:            r.Title,
		Info:             RecipeInfo(r.Info),
		ShortDescription: r.ShortDescription,
		LongDescription:  r.LongDescription,
		Ingredients: func() *datatypes.JSONType[[]RecipeIngredient] {
			if r.Ingredients == nil {
				return nil
			}
			ingredients := make([]RecipeIngredient, len(*r.Ingredients))
			for i, ingredient := range *r.Ingredients {
				ingredients[i] = RecipeIngredient(ingredient)
			}
			return &datatypes.JSONType[[]RecipeIngredient]{Data: ingredients}
		}(),
		Steps: func() *datatypes.JSONType[[]RecipeStep] {
			if r.Steps == nil {
				return nil
			}
			steps := make([]RecipeStep, len(*r.Steps))
			for i, step := range *r.Steps {
				steps[i] = RecipeStep(step)
			}
			return &datatypes.JSONType[[]RecipeStep]{Data: steps}
		}(),
		ImageID: r.ImageID,
	}
}
