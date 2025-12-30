package services

import (
	"errors"
	"member_API/models"
	"member_API/repositories"
	"time"
)

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct 建立新產品
func (s *ProductService) CreateProduct(name string, price float64, description, image string, stock int, creatorId uint) (*models.Product, error) {
	now := time.Now()
	product := &models.Product{
		Base: models.Base{
			CreationTime: now,
			CreatorId:    creatorId,
			IsDeleted:    false,
		},
		ProductName:        name,
		ProductPrice:       price,
		ProductDescription: description,
		ProductImage:       image,
		ProductStock:       stock,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct 更新產品資訊
func (s *ProductService) UpdateProduct(id uint, updates map[string]interface{}, modifierId uint) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("產品不存在")
	}

	now := time.Now()
	updates["last_modification_time"] = &now
	updates["last_modifier_id"] = modifierId

	// 應用更新
	for key, value := range updates {
		switch key {
		case "product_name":
			if v, ok := value.(string); ok {
				product.ProductName = v
			}
		case "product_price":
			if v, ok := value.(float64); ok {
				product.ProductPrice = v
			}
		case "product_description":
			if v, ok := value.(string); ok {
				product.ProductDescription = v
			}
		case "product_image":
			if v, ok := value.(string); ok {
				product.ProductImage = v
			}
		case "product_stock":
			if v, ok := value.(int); ok {
				product.ProductStock = v
			}
		case "last_modification_time":
			if v, ok := value.(*time.Time); ok {
				product.LastModificationTime = v
			}
		case "last_modifier_id":
			if v, ok := value.(uint); ok {
				product.LastModifierId = v
			}
		}
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	// 重新載入產品資料以確保取得最新資料
	updatedProduct, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

// DeleteProduct 軟刪除產品
func (s *ProductService) DeleteProduct(id uint, deleterId uint) error {
	return s.repo.SoftDelete(id, deleterId)
}

// GetProductByID 取得單一產品
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("產品不存在")
	}
	return product, nil
}

// GetProducts 取得產品列表（支持分頁）
func (s *ProductService) GetProducts(limit, offset int) ([]models.Product, int64, error) {
	// 取得總數
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	// 取得產品列表
	products, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
