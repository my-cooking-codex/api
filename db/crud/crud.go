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

func CreateRecipe(recipe db.CreateRecipe, userID uuid.UUID) (db.Recipe, error) {
	var newRecipe = recipe.IntoRecipe(userID, nil)
	if err := db.DB.Create(&newRecipe).Error; err != nil {
		return db.Recipe{}, err
	}
	return newRecipe, nil
}

func GetRecipesByUserID(userID uuid.UUID, offset uint, limit uint) ([]db.Recipe, error) {
	var recipes []db.Recipe
	if err := db.DB.Offset(int(offset)).Limit(int(limit)).Order("created_at DESC").Find(&recipes, "owner_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func GetRecipesByUserIDCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := db.DB.Model(&db.Recipe{}).Where("owner_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetRecipeById(id uuid.UUID) (db.Recipe, error) {
	var recipe db.Recipe
	if err := db.DB.First(&recipe, "id = ?", id).Error; err != nil {
		return db.Recipe{}, err
	}
	return recipe, nil
}

func DoesUserOwnRecipe(userID uuid.UUID, recipeId uuid.UUID) (bool, error) {
	var recipe db.Recipe

	if err := db.DB.Where("id = ?", recipeId).Where("owner_id = ?", userID).First(&recipe).Error; err != nil {
		return false, err
	}
	return true, nil
}

func UpdateRecipe(recipeID uuid.UUID, recipe db.UpdateRecipe) (db.Recipe, error) {
	var updatedRecipe db.Recipe

	if err := db.DB.Model(&updatedRecipe).Where("id = ?", recipeID).Updates(recipe.IntoRecipe()).Error; err != nil {
		return db.Recipe{}, err
	}
	return updatedRecipe, nil
}

func UpdateRecipeImage(recipeID uuid.UUID, imageID *uuid.UUID) error {
	var updatedRecipe db.Recipe
	if err := db.DB.Model(&updatedRecipe).Where("id = ?", recipeID).Updates(map[string]any{"image_id": imageID}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteRecipe(recipeID uuid.UUID) error {
	if err := db.DB.Delete(&db.Recipe{}, recipeID).Error; err != nil {
		return err
	}
	return nil
}
