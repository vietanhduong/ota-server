package metadata

import (
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/database/models"
)

type repository struct {
	*database.DB
}

type Repository interface {
	Insert(req *Metadata) (*models.Metadata, error)
	InsertBulk(req []*Metadata) ([]*models.Metadata, error)
	FindByProfileId(profileId int) ([]*models.Metadata, error)
	FindByListProfileId(profileIds []uint) ([]*models.Metadata, error)
}

func NewRepository(db *database.DB) *repository {
	return &repository{db}
}

func (r *repository) Insert(req *Metadata) (*models.Metadata, error) {
	model := &models.Metadata{
		ProfileId: uint(req.ProfileId),
		Key:       req.Key,
		Value:     req.Value,
	}

	err := r.Create(&model).Error
	return model, err
}

func (r *repository) InsertBulk(req []*Metadata) ([]*models.Metadata, error) {
	var insertModels []*models.Metadata
	for _, r := range req {
		model := &models.Metadata{
			ProfileId: uint(r.ProfileId),
			Key:       r.Key,
			Value:     r.Value,
		}
		insertModels = append(insertModels, model)
	}
	err := r.Create(&insertModels).Error
	return insertModels, err
}

func (r *repository) FindByProfileId(profileId uint) ([]*models.Metadata, error) {
	var metadata []*models.Metadata
	err := r.Where("profile_id", profileId).Find(&metadata).Error
	return metadata, err
}

func (r *repository) FindByListProfileId(profileIds []uint) ([]*models.Metadata, error) {
	var metadata []*models.Metadata
	err := r.Where("profile_id IN ?", profileIds).Find(&metadata).Error
	return metadata, err
}
