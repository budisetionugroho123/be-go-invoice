package repository

import (
	"github.com/budisetionugroho123/be-go-invoice/internal/models"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// SearchByCode returns items whose code matches the given search string (LIKE query).
func (r *ItemRepository) SearchByCode(code string) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Where("code ILIKE ?", "%"+code+"%").Find(&items).Error
	return items, err
}

// FindByID returns a single item by its primary key.
func (r *ItemRepository) FindByID(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByIDs returns items matching the given list of IDs.
func (r *ItemRepository) FindByIDs(ids []uint) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Where("id IN ?", ids).Find(&items).Error
	return items, err
}
