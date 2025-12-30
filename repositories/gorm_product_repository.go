package repositories

import (
	"errors"
	"member_API/models"
	"time"

	"gorm.io/gorm"
)

type gormProductRepository struct {
	db *gorm.DB
}

// NewGormProductRepository 創建 GORM ProductRepository 實例
func NewGormProductRepository(db *gorm.DB) ProductRepository {
	return &gormProductRepository{db: db}
}

func (r *gormProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("is_deleted = ?", false).First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *gormProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *gormProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *gormProductRepository) SoftDelete(id uint, deleterId uint) error {
	now := time.Now()
	result := r.db.Model(&models.Product{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"is_deleted":             true,
			"deleted_at":             &now,
			"last_modifier_id":       deleterId,
			"last_modification_time": &now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("產品不存在或已被刪除")
	}
	return nil
}

func (r *gormProductRepository) FindAll(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("is_deleted = ?", false).
		Order("sort ASC, id DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, err
}

func (r *gormProductRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&models.Product{}).Where("is_deleted = ?", false).Count(&total).Error
	return total, err
}

