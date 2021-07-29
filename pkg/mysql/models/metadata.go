package models

import "gorm.io/gorm"

type Metadata struct {
	*gorm.Model
	ProfileId uint
	Key       string `gorm:"index"`
	Value     string
}
