package crud

import (
	"strings"

	"github.com/google/uuid"
	"github.com/my-cooking-codex/api/db"
	"gorm.io/gorm"
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

func GetLabelNamesByUser(userID uuid.UUID) ([]string, error) {
	var labels []string
	if err := db.DB.Model(&db.Recipe{}).
		Distinct("labels.name").
		Joins("JOIN recipe_labels ON recipes.id = recipe_labels.recipe_id").
		Joins("JOIN labels ON recipe_labels.label_id = labels.id").
		Where("recipes.owner_id = ?", userID).
		Pluck("name", &labels).Error; err != nil {
		return nil, err
	}
	return labels, nil
}

func GetLabelCountByUser(userID uuid.UUID) (int64, error) {
	var count int64
	if err := db.DB.Model(&db.Recipe{}).
		Distinct("recipe_labels.label_id").
		Joins("JOIN recipe_labels ON recipes.id = recipe_labels.recipe_id").
		Where("recipes.owner_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CreateRecipe(recipe db.CreateRecipe, userID uuid.UUID) (db.ReadRecipe, error) {
	var newRecipe = recipe.IntoRecipe(userID, nil)
	labels := make([]db.Label, len(recipe.Labels))

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for i, label := range recipe.Labels {
			labels[i] = db.Label{Name: label}
			if err := db.DB.FirstOrCreate(&labels[i], "name = ?", label).Select("id").Error; err != nil {
				return err
			}
		}

		return db.DB.Create(&newRecipe).Association("Labels").Append(labels)
	})

	return newRecipe.IntoReadRecipe(), err
}

type RecipesFilterParams struct {
	Title         *string
	Labels        []string
	Freezable     *bool
	MicrowaveOnly *bool
}

func GetRecipesByUserID(userID uuid.UUID, offset uint, limit uint, filters RecipesFilterParams) ([]db.ReadRecipe, error) {
	var recipes []db.Recipe

	// build base query
	query := db.DB.Preload("Labels").
		Offset(int(offset)).
		Limit(int(limit)).
		Order("created_at DESC").
		Where("owner_id = ?", userID)

	// add title filter if present
	if filters.Title != nil {
		titleFilter := strings.TrimSpace(*filters.Title)
		if titleFilter != "" {
			query = query.Where("title LIKE ?", "%"+titleFilter+"%")
		}
	}

	// add labels filter if present
	if len(filters.Labels) > 0 {
		query = query.Joins("JOIN recipe_labels ON recipes.id = recipe_labels.recipe_id").
			Joins("JOIN labels ON recipe_labels.label_id = labels.id").
			Where("labels.name IN ?", filters.Labels).
			Group("recipes.id").
			Having("COUNT(DISTINCT labels.name) = ?", len(filters.Labels))
	}

	// add freezable filter if present
	if filters.Freezable != nil {
		query = query.Where("info_freezable = ?", *filters.Freezable)
	}

	// add microwave-only filter if present
	if filters.MicrowaveOnly != nil {
		query = query.Where("info_microwave_only = ?", *filters.MicrowaveOnly)
	}

	if err := query.Find(&recipes).Error; err != nil {
		return nil, err
	}
	readRecipes := make([]db.ReadRecipe, len(recipes))
	for i, recipe := range recipes {
		readRecipes[i] = recipe.IntoReadRecipe()
	}

	return readRecipes, nil
}

func GetRecipesByUserIDCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := db.DB.Model(&db.Recipe{}).Where("owner_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetRecipeById(id uuid.UUID) (db.ReadRecipe, error) {
	var recipe db.Recipe
	if err := db.DB.Preload("Labels").First(&recipe, "id = ?", id).Error; err != nil {
		return db.ReadRecipe{}, err
	}
	return recipe.IntoReadRecipe(), nil
}

func DoesUserOwnRecipe(userID uuid.UUID, recipeId uuid.UUID) (bool, error) {
	var recipe db.Recipe

	if err := db.DB.Where("id = ?", recipeId).Where("owner_id = ?", userID).First(&recipe).Error; err != nil {
		return false, err
	}
	return true, nil
}

func UpdateRecipe(recipeID uuid.UUID, recipe db.UpdateRecipe) (db.ReadRecipe, error) {
	var updatedRecipe db.Recipe

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := db.DB.Model(&updatedRecipe).Where("id = ?", recipeID).Updates(recipe.IntoRecipe()).Error; err != nil {
			return err
		}

		if recipe.Labels != nil {
			labels := make([]db.Label, len(*recipe.Labels))
			for i, label := range *recipe.Labels {
				labels[i] = db.Label{Name: label}
				if err := db.DB.FirstOrCreate(&labels[i], "name = ?", label).Select("id").Error; err != nil {
					return err
				}
			}
			var foundRecipe db.Recipe
			if err := db.DB.First(&foundRecipe, recipeID).Select("id").Error; err != nil {
				return err
			}
			return db.DB.Model(&foundRecipe).Association("Labels").Replace(&labels)
		}

		return nil
	})

	return updatedRecipe.IntoReadRecipe(), err
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
