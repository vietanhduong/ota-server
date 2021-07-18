package models

import "gorm.io/gorm"

type IosProfile struct {
	*gorm.Model
	AppName          string
	BundleIdentifier string
	AppVersion       string
	BuildVersion     uint  
	StorageObjectID  uint
	StorageObject    StorageObject
}
