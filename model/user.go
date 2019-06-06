package model

import (
	"time"
)

type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	CreateTime time.Time
	Gender     uint8
	Articles   []Article `gorm:"foreignkey:CreateUser"`

	// LastLoginTime time.Time
}
