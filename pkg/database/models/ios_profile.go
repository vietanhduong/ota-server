package models

import "gorm.io/gorm"

type IosProfile struct {
	*gorm.Model
	AppName          string
	BundleIdentifier string
	Version          string
	Build            uint
	StorageObjectID  uint
	StorageObject    StorageObject
}
