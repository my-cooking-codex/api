package crud

import (
	"github.com/google/uuid"
	"github.com/my-cooking-codex/api/db"
)

func CreateUser(user db.CreateUser) (db.User, error) {
	var newUser = user.IntoUser()
	if err := db.DB.Create(&newUser).Error; err != nil {
		return db.User{}, err
	}
	return newUser, nil
}

func GetUserById(userID uuid.UUID) (db.User, error) {
	var user db.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return db.User{}, err
	}
	return user, nil
}

func GetUserByUsername(username string) (db.User, error) {
	var user db.User
	if err := db.DB.First(&user, "username = ?", username).Error; err != nil {
		return db.User{}, err
	}
	return user, nil
}

func GetUserCount() (int64, error) {
	var count int64
	if err := db.DB.Model(&db.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
