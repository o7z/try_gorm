package model

import "time"

type Article struct {
	ID         string `gorm:"size:50"`
	Title      string
	Content    string `gorm:"type:blob"`
	CreateUser string `gorm:"size:50"`
	CreateTime time.Time
}
