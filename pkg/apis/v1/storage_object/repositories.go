package storage_object

import "github.com/vietanhduong/ota-server/pkg/database"

type repository struct {
	*database.DB
}

func NewRepository(db *database.DB) *repository {
	return &repository{db}
}