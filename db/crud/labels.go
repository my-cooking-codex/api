package crud

import (
	"github.com/google/uuid"
	"github.com/my-cooking-codex/api/db"
)

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
