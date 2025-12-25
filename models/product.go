package models

// Member represents a user stored in PostgreSQL and managed by GORM.
type Product struct {
	ProductName         string `gorm:"size:255;not null" json:"product_name"`
	ProductPrice        int    `gorm:"size:255;uniqueIndex;not null" json:"product_price"`
	ProductDescription  string `gorm:"size:255" json:"product_description"`
	Base
}
