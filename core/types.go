package core

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

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

type RecipesFilterParams struct {
	Page    uint `query:"page" validate:"required,gt=0"`
	PerPage uint `query:"perPage" validate:"required,gt=0,lte=120"`
}
