package models

type IosProfile struct {
	base
	AppName          string
	BundleIdentifier string
	AppVersion       string
	BuildVersion     uint  
	StorageObjectID  uint
	StorageObject    StorageObject
}
