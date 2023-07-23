package core

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type SelectField string

func (sf *SelectField) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) != 0 {
		runes := []rune(s)
		s = strings.ToUpper(string(runes[0])) + string(runes[1:])
	}
	*sf = SelectField(s)
	return nil
}

type SelectedUpdate[T any] struct {
	Fields []SelectField `json:"fields" validate:"required"`
	Model  T             `json:"model" validate:"required"`
}

func (su *SelectedUpdate[T]) FieldsAsString() []string {
	fields := make([]string, len(su.Fields))
	for i, field := range su.Fields {
		fields[i] = string(field)
	}
	return fields
}

type AuthenticatedUser struct {
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
	IsAdmin  bool      `json:"isAdmin"`
}

type JWTClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) ToAuthenticatedUser() (AuthenticatedUser, error) {
	if userID, err := uuid.Parse(c.Subject); err != nil {
		return AuthenticatedUser{}, err
	} else {
		return AuthenticatedUser{
			UserID:   userID,
			Username: c.Username,
			IsAdmin:  c.IsAdmin,
		}, nil
	}
}

type LoginToken struct {
	Type   string    `json:"type"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

type CreateLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PaginationParams struct {
	Page    uint `query:"page" validate:"required,gt=0"`
	PerPage uint `query:"perPage" validate:"required,gt=0,lte=120"`
}

type RecipesFilterParams struct {
	PaginationParams
	Title         *string  `query:"title"`
	Labels        []string `query:"label"`
	Freezable     *bool    `query:"freezable"`
	MicrowaveOnly *bool    `query:"microwaveOnly"`
}

type PantryItemsFilterParams struct {
	PaginationParams
	Name       string     `query:"name"`
	Labels     []string   `query:"label"`
	LocationId *uuid.UUID `query:"locationId"`
	Expired    *bool      `query:"expired"`
}
