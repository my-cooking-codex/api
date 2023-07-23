package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/my-cooking-codex/api/db"
)

type ReadPantryItem struct {
	db.UUIDBase
	db.TimeBase
	Name       string     `json:"name"`
	LocationId uuid.UUID  `json:"locationId"`
	Quantity   uint       `json:"quantity"`
	Notes      *string    `json:"notes,omitempty"`
	Expiry     *time.Time `json:"expiry,omitempty"`
	Labels     []string   `json:"labels,omitempty"`
}

type CreatePantryLocation struct {
	Name string `json:"name" validate:"required,min=1,max=60"`
}

type CreatePantryItem struct {
	Name     string     `json:"name" validate:"required,min=1,max=60"`
	Quantity uint       `json:"quantity"`
	Notes    *string    `json:"notes,omitempty"`
	Expiry   *time.Time `json:"expiry"`
	Labels   []string   `json:"labels,omitempty" validate:"dive,min=1,max=60"`
}

type UpdatePantryLocation struct {
	Name string `json:"name,omitempty" validate:"omitempty,min=1,max=60"`
}

type UpdatePantryItem struct {
	Name       string     `json:"name,omitempty" validate:"omitempty,required,min=1,max=60"`
	LocationId uuid.UUID  `json:"locationId,omitempty"`
	Quantity   uint       `json:"quantity,omitempty"`
	Notes      *string    `json:"notes,omitempty"`
	Expiry     *time.Time `json:"expiry,omitempty"`
	Labels     []string   `json:"labels,omitempty" validate:"omitempty,dive,min=1,max=60"`
}
