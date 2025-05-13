package user

import (
	"context"
	"fmt"
	"time"

	"go.learning/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *models.User) error
	GetUserList(queryParams GetUserList) ([]models.User, int64, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUser(user *models.User) error {
	err := r.db.WithContext(context.Background()).Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *repository) GetUserList(queryParams GetUserList) ([]models.User, int64, error) {
	var users []models.User
	query := r.db.WithContext(context.Background()).Where("deleted_at IS NULL")
	if queryParams.FirstName != nil {
		query = query.Where("first_name ILIKE ?", "%"+*queryParams.FirstName+"%")
	}

	// Count total records
	var total int64
	err := query.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Set sorting and pagination
	if queryParams.Sort != "" {
		if queryParams.SortDirection == "desc" {
			query = query.Order(queryParams.Sort + " DESC")
		} else {
			query = query.Order(queryParams.Sort + " ASC")
		}
	} else {
		query = query.Order("first_name ASC")
	}
	if queryParams.Page > 0 && queryParams.Limit > 0 {
		offset := (queryParams.Page - 1) * queryParams.Limit
		query = query.Offset(offset).Limit(queryParams.Limit)
	} else {
		query = query.Limit(10)
	}
	err = query.Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user list: %w", err)
	}
	return users, total, nil
}

func (r *repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *models.User) error {
	err := r.db.WithContext(context.Background()).Model(user).Updates(map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"active":     user.Active,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *repository) DeleteUser(id uint) error {
	user, err := r.GetUserByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	err = r.db.Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
