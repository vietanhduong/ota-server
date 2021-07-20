package models

import "gorm.io/gorm"

type StorageObject struct {
	*gorm.Model
	// Name the original filename
	Name string
	// Path the abs path after uploaded to GCS
	Path string
	// ContentType content-type
	ContentType string
}
