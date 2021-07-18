package models

import "gorm.io/gorm"

type StorageObject struct {
	*gorm.Model
	Name        string
	Path        string
	ContentType string
}
