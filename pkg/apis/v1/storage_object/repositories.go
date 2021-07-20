package storage_object

import (
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/database/models"
)

type repository struct {
	*database.DB
}

type Repository interface {
	Insert(uploadedFile *UploadedFile) (*models.StorageObject, error)
}

func NewRepository(db *database.DB) *repository {
	return &repository{db}
}

func (r *repository) Insert(uploadedFile *UploadedFile) (*models.StorageObject, error) {
	object := &models.StorageObject{
		Name:        uploadedFile.Filename,
		Path:        uploadedFile.AbsPath,
		ContentType: uploadedFile.ContentType,
	}

	err := r.Create(&object).Error
	return object, err
}
