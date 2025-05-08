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
