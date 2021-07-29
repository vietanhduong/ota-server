package models

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Email       string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string `gorm:"type:varchar(255)"`
	DisplayName string
	Active      bool `gorm:"default:1"`
}
