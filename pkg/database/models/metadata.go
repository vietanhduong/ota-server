package models

import "gorm.io/gorm"

type Metadata struct {
	*gorm.Model
	Type      string
	ProfileId uint
	Key       string `gorm:"index"`
	Value     string
}
