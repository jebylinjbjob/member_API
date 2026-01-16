package services

import (
	"fmt"
	"member_API/testutil"
	"testing"
)

func TestProductService_CreateProduct(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(db)

	service := NewProductService(db)

	tests := []struct {
		name        string
		productName string
		price       float64
		description string
		image       string
		stock       int
		creatorId   uint
		wantErr     bool
	}{
		{
			name:        "Valid product creation",
			productName: "Test Product",
			price:       99.99,
			description: "Test description",
			image:       "test.jpg",
			stock:       100,
			creatorId:   1,
			wantErr:     false,
		},
		{
			name:        "Product with zero price",
			productName: "Free Product",
			price:       0.0,
			description: "Free item",
			image:       "free.jpg",
			stock:       50,
			creatorId:   1,
			wantErr:     false,
		},
		{
			name:        "Product with zero stock",
			productName: "Out of Stock Product",
			price:       29.99,
			description: "Currently unavailable",
			image:       "unavailable.jpg",
			stock:       0,
			creatorId:   1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := service.CreateProduct(tt.productName, tt.price, tt.description, tt.image, tt.stock, tt.creatorId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if product == nil {
					t.Error("CreateProduct() returned nil product")
					return
				}
				if product.ProductName != tt.productName {
					t.Errorf("CreateProduct() ProductName = %v, want %v", product.ProductName, tt.productName)
				}
				if product.ProductPrice != tt.price {
					t.Errorf("CreateProduct() ProductPrice = %v, want %v", product.ProductPrice, tt.price)
				}
				if product.ProductStock != tt.stock {
					t.Errorf("CreateProduct() ProductStock = %v, want %v", product.ProductStock, tt.stock)
				}
			}
		})
	}
}

func TestProductService_GetProductByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(db)

	service := NewProductService(db)

	// Create a test product
	created, err := service.CreateProduct("Test Product", 99.99, "Description", "test.jpg", 100, 1)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "Existing product",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "Non-existing product",
			id:      9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := service.GetProductByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if product == nil {
					t.Error("GetProductByID() returned nil product")
					return
				}
				if product.ID != tt.id {
					t.Errorf("GetProductByID() ID = %v, want %v", product.ID, tt.id)
				}
			}
		})
	}
}

func TestProductService_GetProducts(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(db)

	service := NewProductService(db)

	// Create multiple test products
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("Product %c", 'A'+i)
		price := 99.99 + float64(i) // Different price for each product due to uniqueIndex
		stock := 100 + i            // Different stock for each product due to uniqueIndex
		_, err := service.CreateProduct(name, price, "Description", "test.jpg", stock, 1)
		if err != nil {
			t.Fatalf("Failed to create test product: %v", err)
		}
	}

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:      "Get all products",
			limit:     10,
			offset:    0,
			wantCount: 5,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "Get limited products",
			limit:     3,
			offset:    0,
			wantCount: 3,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:      "Get products with offset",
			limit:     10,
			offset:    2,
			wantCount: 3,
			wantTotal: 5,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, total, err := service.GetProducts(tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(products) != tt.wantCount {
					t.Errorf("GetProducts() count = %v, want %v", len(products), tt.wantCount)
				}
				if total != tt.wantTotal {
					t.Errorf("GetProducts() total = %v, want %v", total, tt.wantTotal)
				}
			}
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(db)

	service := NewProductService(db)

	// Create a test product
	created, err := service.CreateProduct("Test Product", 99.99, "Description", "test.jpg", 100, 1)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	tests := []struct {
		name       string
		id         uint
		updates    map[string]interface{}
		modifierId uint
		wantErr    bool
	}{
		{
			name: "Update product name",
			id:   created.ID,
			updates: map[string]interface{}{
				"product_name": "Updated Product",
			},
			modifierId: 2,
			wantErr:    false,
		},
		{
			name: "Update product price and stock",
			id:   created.ID,
			updates: map[string]interface{}{
				"product_price": 149.99,
				"product_stock": 200,
			},
			modifierId: 2,
			wantErr:    false,
		},
		{
			name: "Update non-existing product",
			id:   9999,
			updates: map[string]interface{}{
				"product_name": "Should Fail",
			},
			modifierId: 2,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := service.UpdateProduct(tt.id, tt.updates, tt.modifierId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if product == nil {
					t.Error("UpdateProduct() returned nil product")
					return
				}
				if product.LastModifierId != tt.modifierId {
					t.Errorf("UpdateProduct() LastModifierId = %v, want %v", product.LastModifierId, tt.modifierId)
				}
				if product.LastModificationTime == nil {
					t.Error("UpdateProduct() LastModificationTime is nil")
				}

				// Verify specific field updates
				for key, expectedValue := range tt.updates {
					switch key {
					case "product_name":
						if product.ProductName != expectedValue.(string) {
							t.Errorf("UpdateProduct() ProductName = %v, want %v", product.ProductName, expectedValue)
						}
					case "product_price":
						if product.ProductPrice != expectedValue.(float64) {
							t.Errorf("UpdateProduct() ProductPrice = %v, want %v", product.ProductPrice, expectedValue)
						}
					case "product_stock":
						if product.ProductStock != expectedValue.(int) {
							t.Errorf("UpdateProduct() ProductStock = %v, want %v", product.ProductStock, expectedValue)
						}
					}
				}
			}
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(db)

	service := NewProductService(db)

	// Create a test product
	created, err := service.CreateProduct("Test Product", 99.99, "Description", "test.jpg", 100, 1)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	tests := []struct {
		name      string
		id        uint
		deleterId uint
		wantErr   bool
	}{
		{
			name:      "Valid delete",
			id:        created.ID,
			deleterId: 2,
			wantErr:   false,
		},
		{
			name:      "Delete non-existing product",
			id:        9999,
			deleterId: 2,
			wantErr:   true,
		},
		{
			name:      "Delete already deleted product",
			id:        created.ID,
			deleterId: 2,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteProduct(tt.id, tt.deleterId)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the product is marked as deleted (soft delete)
			if !tt.wantErr {
				product, err := service.GetProductByID(tt.id)
				if err == nil {
					t.Errorf("DeleteProduct() product should not be retrievable after deletion, got product: %+v", product)
				}
			}
		})
	}
}
