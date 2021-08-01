package models

import "gorm.io/gorm"

type Profile struct {
	*gorm.Model
	AppName          string
	BundleIdentifier string
	Version          string
	Build            uint
	StorageObjectID  uint
	StorageObject    *StorageObject
	UserID           uint
	User             *User
}
