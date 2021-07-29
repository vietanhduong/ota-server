package database

import "github.com/vietanhduong/ota-server/pkg/database/models"

func (db *DB) Migration() (err error) {
	// migrate storage object model
	if err = db.AutoMigrate(&models.StorageObject{}); err != nil {
		return
	}

	// migrate ios_profile model
	if err = db.AutoMigrate(&models.Profile{}); err != nil {
		return
	}

	// migrate metadata model
	if err = db.AutoMigrate(&models.Metadata{}); err != nil {
		return
	}

	return
}
