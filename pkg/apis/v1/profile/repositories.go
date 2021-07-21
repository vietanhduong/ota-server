package profile

import (
	"errors"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/database/models"
	"gorm.io/gorm"
)

type repository struct {
	*database.DB
}

type Repository interface {
	Insert(req *RequestProfile) (*models.Profile, error)
	FindById(objectId uint) (*models.Profile, error)
	All() ([]*models.Profile, error)
}

func NewRepository(db *database.DB) *repository {
	return &repository{db}
}

func (r *repository) Insert(req *RequestProfile) (*models.Profile, error) {
	model := &models.Profile{
		AppName:          req.AppName,
		BundleIdentifier: req.BundleIdentifier,
		Version:          req.Version,
		Build:            uint(req.Build),
		StorageObjectID:  uint(req.StorageObjectID),
	}

	err := r.Create(&model).Error
	return model, err
}

func (r *repository) FindById(objectId uint) (*models.Profile, error) {
	var model *models.Profile
	err := r.Where("id = ?", objectId).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return model, err
}

func (r *repository) All() ([]*models.Profile, error) {
	var profiles []*models.Profile
	err := r.Find(&profiles).Error
	return profiles, err
}

