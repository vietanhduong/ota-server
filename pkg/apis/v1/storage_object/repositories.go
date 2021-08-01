package storage_object

import (
	"errors"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/mysql/models"
	"gorm.io/gorm"
)

type repository struct {
	*mysql.DB
}

type Repository interface {
	Insert(uploadedFile *File) (*models.StorageObject, error)
	FindById(objectId uint) (*models.StorageObject, error)
}

func NewRepository(db *mysql.DB) *repository {
	return &repository{db}
}

func (r *repository) Insert(uploadedFile *File) (*models.StorageObject, error) {
	object := &models.StorageObject{
		Name:        uploadedFile.Filename,
		Key:         uploadedFile.Key,
		Path:        uploadedFile.AbsPath,
		ContentType: uploadedFile.ContentType,
		UserID:      uint(uploadedFile.UploadedBy),
	}

	err := r.Create(&object).Error
	return object, err
}

func (r *repository) FindById(objectId uint) (*models.StorageObject, error) {
	var model *models.StorageObject
	err := r.Where("id = ?", objectId).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return model, err
}

func (r *repository) FindByKey(objectKey string) (*models.StorageObject, error) {
	var model *models.StorageObject
	err := r.First(&model, "`key` = ?", objectKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return model, err
}
