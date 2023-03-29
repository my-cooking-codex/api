package db

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type UUIDBase struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid" json:"id"`
}

type TimeBase struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (base *UUIDBase) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	return
}

type User struct {
	UUIDBase
	TimeBase
	Username       string   `gorm:"uniqueIndex;not null;type:varchar(30)" json:"username"`
	HashedPassword []byte   `gorm:"not null" json:"-"`
	Recipes        []Recipe `gorm:"foreignKey:OwnerID" json:"-"`
}

func (u *User) SetPassword(newPlainPassword string) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(newPlainPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.HashedPassword = hashedPw
}

func (u *User) IsPasswordMatch(plainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(plainPassword)); err == nil {
		return true
	}
	return false
}

type RecipeInfo struct {
	Yields *datatypes.JSONType[RecipeInfoYields] `gorm:"type:json" json:"yields,omitempty"`
}

type Recipe struct {
	UUIDBase
	TimeBase
	OwnerID          uuid.UUID                               `gorm:"not null;type:uuid" json:"ownerId"`
	Title            string                                  `gorm:"not null;type:varchar(30)" json:"title"`
    Info             RecipeInfo                              `gorm:"embedded;embeddedPrefix:info_" json:"info"`
	ShortDescription *string                                 `gorm:"type:varchar(256)" json:"shortDescription,omitempty"`
	LongDescription  *string                                 `json:"longDescription,omitempty"`
	Ingredients      *datatypes.JSONType[[]RecipeIngredient] `gorm:"type:json" json:"ingredients,omitempty"`
	Steps            *datatypes.JSONType[[]RecipeStep]       `gorm:"type:json" json:"steps,omitempty"`
	ImageID          *uuid.UUID                              `gorm:"type:uuid" json:"imageId"`
}
