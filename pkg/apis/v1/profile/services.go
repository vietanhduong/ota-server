package profile

import (
	"github.com/vietanhduong/ota-server/pkg/apis/v1/storage_object"
	"github.com/vietanhduong/ota-server/pkg/database"
)

type StorageService interface {
	GetObject(objectId int) (*storage_object.File, error)
}

type service struct {
	repo       *repository
	storageSvc StorageService
}

func NewService(db *database.DB) *service {
	return &service{
		repo:       NewRepository(db),
		storageSvc: storage_object.NewService(db),
	}
}

func (s *service) CreateProfile(reqProfile *RequestProfile) (*ResponseProfile, error) {
	// TODO: update validate before insert to database
	// validate storage object
	_, err := s.storageSvc.GetObject(reqProfile.StorageObjectID)
	if err != nil {
		return nil, err
	}
	// insert to database
	model, err := s.repo.Insert(reqProfile)
	if err != nil {
		return nil, err
	}

	return &ResponseProfile{
		ProfileId:        model.ID,
		AppName:          model.AppName,
		BundleIdentifier: model.BundleIdentifier,
		Version:          model.Version,
		Build:            model.Build,
		StorageObjectID:  model.StorageObjectID,
	}, err
}
