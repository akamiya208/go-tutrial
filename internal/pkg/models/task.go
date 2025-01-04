package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint   `gorm:"primarykey;autoIncrement"`
	Name        string `gorm:"not null;size:255"`
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
