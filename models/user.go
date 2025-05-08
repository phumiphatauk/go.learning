package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email          string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	FirstName      string         `gorm:"size:100;not null" json:"first_name"`
	LastName       string         `gorm:"size:100;not null" json:"last_name"`
	HashedPassword string         `gorm:"size:255;not null" json:"-"`
	Active         bool           `gorm:"default:true" json:"active"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // soft delete
}
