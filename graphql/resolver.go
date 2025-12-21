package graphql

import (
	"errors"
	"fmt"

	"member_API/models"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

// User 表示會員模型
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Resolver 包含所有 resolvers
type Resolver struct {
	DB *gorm.DB
}

// NewResolver 創建新的 resolver
func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{DB: db}
}

// GetUsers 獲取所有會員
func (r *Resolver) GetUsers(p graphql.ResolveParams) (interface{}, error) {
	if r.DB == nil {
		return []User{}, nil
	}

	var members []models.Member
	if err := r.DB.WithContext(p.Context).
		Select("id", "name", "email").
		Limit(50).
		Find(&members).Error; err != nil {
		return nil, err
	}

	users := make([]User, len(members))
	for i, member := range members {
		users[i] = User{ID: int64(member.ID), Name: member.Name, Email: member.Email}
	}

	return users, nil
}

// GetUserByID 根據 ID 獲取單個會員
func (r *Resolver) GetUserByID(p graphql.ResolveParams) (interface{}, error) {
	if r.DB == nil {
		return nil, errors.New("database connection not configured")
	}

	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("id must be an integer")
	}

	var member models.Member
	if err := r.DB.WithContext(p.Context).
		Select("id", "name", "email").
		First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return User{ID: int64(member.ID), Name: member.Name, Email: member.Email}, nil
}

// CreateUser 創建新會員
func (r *Resolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	if r.DB == nil {
		return nil, errors.New("database connection not configured")
	}

	name, ok := p.Args["name"].(string)
	if !ok {
		return nil, errors.New("name is required")
	}

	email, ok := p.Args["email"].(string)
	if !ok {
		return nil, errors.New("email is required")
	}

	member := models.Member{Name: name, Email: email}
	if err := r.DB.WithContext(p.Context).Create(&member).Error; err != nil {
		return nil, err
	}

	return User{ID: int64(member.ID), Name: member.Name, Email: member.Email}, nil
}

// UpdateUser 更新會員信息
func (r *Resolver) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	if r.DB == nil {
		return nil, errors.New("database connection not configured")
	}

	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("id must be an integer")
	}

	var member models.Member
	if err := r.DB.WithContext(p.Context).First(&member, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	if newName, ok := p.Args["name"].(string); ok && newName != "" {
		member.Name = newName
	}
	if newEmail, ok := p.Args["email"].(string); ok && newEmail != "" {
		member.Email = newEmail
	}

	if err := r.DB.WithContext(p.Context).Save(&member).Error; err != nil {
		return nil, err
	}

	return User{ID: int64(member.ID), Name: member.Name, Email: member.Email}, nil
}

// DeleteUser 刪除會員
func (r *Resolver) DeleteUser(p graphql.ResolveParams) (interface{}, error) {
	if r.DB == nil {
		return false, errors.New("database connection not configured")
	}

	id, ok := p.Args["id"].(int)
	if !ok {
		return false, errors.New("id must be an integer")
	}

	result := r.DB.WithContext(p.Context).Delete(&models.Member{}, id)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
