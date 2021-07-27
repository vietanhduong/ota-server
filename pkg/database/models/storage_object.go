package models

import "gorm.io/gorm"

type StorageObject struct {
	*gorm.Model
	// Name the original filename
	Name string
	// Key object key
	Key string `gorm:"type:varchar(255);uniqueIndex"`
	// Path the abs path after uploaded to GCS
	Path string
	// ContentType content-type
	ContentType string
}
