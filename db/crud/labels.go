package crud

import (
	"github.com/google/uuid"
	"github.com/my-cooking-codex/api/db"
	"gorm.io/gorm"
)

func GetLabelNamesByUser(userID uuid.UUID) ([]string, error) {
	var labels []string
	tx := db.DB.Session(&gorm.Session{PrepareStmt: true})
	err := tx.
		Model(&db.User{}).
		Raw(`SELECT DISTINCT d1.name AS label_name
FROM (
    SELECT labels.name
    FROM users u
    LEFT JOIN recipes r ON u.id = r.owner_id
    LEFT JOIN recipe_labels r1 ON r.id = r1.recipe_id
    LEFT JOIN labels ON r1.label_id = labels.id
    WHERE u.id = ? AND labels.name IS NOT NULL

    UNION

    SELECT labels.name
    FROM users u
    LEFT JOIN pantry_locations pl ON u.id = pl.owner_id
    LEFT JOIN pantry_items pi ON pl.id = pi.location_id
    LEFT JOIN pantry_item_labels pil ON pi.id = pil.pantry_item_id
    LEFT JOIN labels ON pil.label_id = labels.id
    WHERE u.id = ? AND labels.name IS NOT NULL
) AS d1;`, userID, userID).
		Scan(&labels).Error
	return labels, err
}

func GetLabelCountByUser(userID uuid.UUID) (int64, error) {
	var count int64
	tx := db.DB.Session(&gorm.Session{PrepareStmt: true})
	err := tx.
		Model(&db.User{}).
		Raw(`SELECT COUNT(DISTINCT d1.id) AS label_count
FROM (
    SELECT l.id
    FROM users u
    LEFT JOIN recipes r ON u.id = r.owner_id
    LEFT JOIN recipe_labels r1 ON r.id = r1.recipe_id
    LEFT JOIN labels l ON r1.label_id = l.id
	WHERE u.id = ?

    UNION

    SELECT l.id
    FROM users u
    LEFT JOIN pantry_locations pl ON u.id = pl.owner_id
    LEFT JOIN pantry_items pi ON pl.id = pi.location_id
    LEFT JOIN pantry_item_labels pil ON pi.id = pil.pantry_item_id
    LEFT JOIN labels l ON pil.label_id = l.id
	WHERE u.id = ?
) AS d1;`, userID, userID).
		Scan(&count).Error
	return count, err
}
