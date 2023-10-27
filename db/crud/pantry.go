package crud

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/my-cooking-codex/api/core"
	"github.com/my-cooking-codex/api/db"
	"github.com/my-cooking-codex/api/db/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreatePantryLocation(
	newLocation types.CreatePantryLocation,
	ownerID uuid.UUID,
) (db.PantryLocation, error) {
	pantryLocation := db.PantryLocation{
		OwnerId: ownerID,
		Name:    newLocation.Name,
	}
	err := db.DB.Create(&pantryLocation).Error
	return pantryLocation, err
}

func GetPantryLocationByID(pantryID uuid.UUID) (db.PantryLocation, error) {
	var pantryLocation db.PantryLocation
	err := db.DB.First(&pantryLocation, "id = ?", pantryID).Error
	return pantryLocation, err
}

func GetPantryLocationsByUserID(userID uuid.UUID) ([]db.PantryLocation, error) {
	var pantryLocations []db.PantryLocation
	err := db.DB.Where("owner_id = ?", userID).Find(&pantryLocations).Error
	return pantryLocations, err
}

func DoesUserOwnPantryLocation(userID uuid.UUID, locationID uuid.UUID) (bool, error) {
	var count int64
	err := db.DB.
		Model(&db.PantryLocation{}).
		Where("id = ? AND owner_id = ?", locationID, userID).
		Count(&count).
		Error
	return count > 0, err
}

func UpdatePantryLocation(
	locationID uuid.UUID,
	update core.SelectedUpdate[types.UpdatePantryLocation],
) error {
	return db.DB.
		Model(&db.PantryLocation{}).
		Where("id = ?", locationID).
		Select(update.FieldsAsString()).
		Updates(update.Model).
		Error
}

func DeletePantryLocation(locationID uuid.UUID) error {
	return db.DB.Where("id = ?", locationID).Delete(&db.PantryLocation{}).Error
}

func CreatePantryItem(
	newItem types.CreatePantryItem,
	locationID uuid.UUID,
) (db.PantryItem, error) {
	pantryItem := db.PantryItem{
		Name:       newItem.Name,
		Quantity:   newItem.Quantity,
		Notes:      newItem.Notes,
		Expiry:     newItem.Expiry,
		LocationId: locationID,
	}
	labels := make([]db.Label, len(newItem.Labels))
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		for i, labelName := range newItem.Labels {
			labels[i] = db.Label{Name: labelName}
			if err := tx.FirstOrCreate(&labels[i], "name = ?", labelName).Select("id").Error; err != nil {
				return err
			}
		}
		return tx.Create(&pantryItem).Association("Labels").Append(labels)
	})
	return pantryItem, err
}

func GetPantryItemCountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	return count, db.DB.Model(&db.PantryItem{}).Count(&count).Error
}

func GetPantryItemByID(itemID uuid.UUID) (types.ReadPantryItem, error) {
	var item db.PantryItem
	err := db.DB.Preload("Labels").First(&item, "id = ?", itemID).Error
	readPantryItem := types.ReadPantryItem{
		UUIDBase:   item.UUIDBase,
		TimeBase:   item.TimeBase,
		Name:       item.Name,
		LocationId: item.LocationId,
		Quantity:   item.Quantity,
		Notes:      item.Notes,
		Expiry:     item.Expiry,
		Labels: func() []string {
			labels := make([]string, len(item.Labels))
			for i, label := range item.Labels {
				labels[i] = label.Name
			}
			return labels
		}(),
	}
	return readPantryItem, err
}

type PantryItemsFilters struct {
	Name       string
	Labels     []string
	LocationId *uuid.UUID
	Expired    *bool
}

func GetPantryItemsByUserID(
	userID uuid.UUID,
	offset uint,
	limit uint,
	filters PantryItemsFilters,
) ([]types.ReadPantryItem, error) {
	query := db.DB.
		Preload("Labels").
		Preload("Location").
		Offset(int(offset)).
		Limit(int(limit)).
		Order("expiry ASC").
		Joins("JOIN pantry_locations ON pantry_items.location_id = pantry_locations.id").
		Where("pantry_locations.owner_id = ?", userID)

	if name := strings.TrimSpace(filters.Name); name != "" {
		query = query.Where("pantry_items.name LIKE ?", "%"+name+"%")
	}

	if len(filters.Labels) != 0 {
		query = query.Joins("JOIN pantry_item_labels ON pantry_items.id = pantry_item_labels.pantry_item_id").
			Joins("JOIN labels ON pantry_item_labels.label_id = labels.id").
			Where("labels.name IN ?", filters.Labels).
			Group("pantry_items.id").
			Having("COUNT(DISTINCT labels.name) = ?", len(filters.Labels))
	}

	if filters.LocationId != nil {
		query = query.Where("location_id = ?", filters.LocationId)
	}

	if filters.Expired != nil {
		currentTime := time.Now().UTC()
		if *filters.Expired {
			query = query.Where("expiry <= ?", currentTime)
		} else {
			query = query.Where("expiry > ? OR expiry IS NULL", currentTime)
		}
	}

	var items []db.PantryItem
	err := query.Find(&items).Error
	var readItems = make([]types.ReadPantryItem, len(items))

	for i, item := range items {
		readItems[i] = types.ReadPantryItem{
			UUIDBase:   item.UUIDBase,
			TimeBase:   item.TimeBase,
			Name:       item.Name,
			LocationId: item.LocationId,
			Quantity:   item.Quantity,
			Notes:      item.Notes,
			Expiry:     item.Expiry,
			Labels: func() []string {
				labels := make([]string, len(item.Labels))
				for i, label := range item.Labels {
					labels[i] = label.Name
				}
				return labels
			}(),
		}
	}

	return readItems, err
}

func DoesUserOwnPantryItem(userID uuid.UUID, pantryItemID uuid.UUID) (bool, error) {
	var count int64
	err := db.DB.
		Model(&db.PantryItem{}).
		Preload("Location").
		Joins("JOIN pantry_locations ON pantry_items.location_id = pantry_locations.id").
		Where("pantry_items.id = ? AND pantry_locations.owner_id = ?", pantryItemID, userID).
		Count(&count).
		Error
	return count != 0, err
}

func UpdatePantryItem(
	itemID uuid.UUID,
	update core.SelectedUpdate[types.UpdatePantryItem],
) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		fields := core.PopElement("Labels", update.FieldsAsString())
		if err := tx.Model(&db.PantryItem{}).Where("id = ?", itemID).Select(fields).Updates(db.PantryItem{
			Name:       update.Model.Name,
			LocationId: update.Model.LocationId,
			Quantity:   update.Model.Quantity,
			Notes:      update.Model.Notes,
			Expiry:     update.Model.Expiry,
		}).Error; err != nil {
			return err
		}
		if len(fields) != len(update.Fields) {
			var item db.PantryItem
			if err := tx.First(&item, "id = ?", itemID).Select("id").Error; err != nil {
				return err
			}
			query := tx.Model(&item).Association("Labels")
			if length := len(update.Model.Labels); length == 0 {
				return query.Clear()
			} else {
				var labels = make([]db.Label, length)
				for i, labelName := range update.Model.Labels {
					labels[i] = db.Label{Name: labelName}
					if err := tx.FirstOrCreate(&labels[i], "name = ?", labelName).Select("id").Error; err != nil {
						return err
					}
				}
				query.Replace(&labels)
			}
		}
		return nil
	})
}

func DeletePantryItem(itemID uuid.UUID) error {
	return db.DB.Select(clause.Associations).Delete(&db.PantryItem{}, itemID).Error
}
