package profile

import (
	"errors"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/metadata"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/storage_object"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/database"
	"net/http"
)

type StorageService interface {
	GetObject(objectId int) (*storage_object.File, error)
}

type MetadataService interface {
	CreateMetadata(profileId int, metadata map[string]string) ([]*metadata.Metadata, error)
	GetMetadata(profileId int) ([]*metadata.Metadata, error)
	GetMetadataByListProfileId(profileIds []uint) (map[uint][]*metadata.Metadata, error)
}

type service struct {
	repo        *repository
	storageSvc  StorageService
	metadataSvc MetadataService
}

func NewService(db *database.DB) *service {
	return &service{
		repo:        NewRepository(db),
		storageSvc:  storage_object.NewService(db),
		metadataSvc: metadata.NewService(db),
	}
}

func (s *service) GetProfiles() ([]*ResponseProfile, error) {
	profiles, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	// prepare profile ids
	var profileIds []uint
	for _, p := range profiles {
		profileIds = append(profileIds, p.ID)
	}

	// fetch metadata
	mm, err := s.metadataSvc.GetMetadataByListProfileId(profileIds)
	if err != nil {
		return nil, err
	}

	// convert to response object
	var result []*ResponseProfile
	for _, p := range profiles {
		profile := ToResponseProfile(p)
		if m, ok := mm[profile.ProfileId]; ok {
			profile.Metadata = ConvertMetadataListToMap(m)
		}
		result = append(result, profile)
	}

	return result, nil
}

func (s *service) GetProfile(profileId int) (*ResponseProfile, error) {
	model, err := s.repo.FindById(uint(profileId))
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, cerrors.NewCError(http.StatusNotFound, errors.New("profile does not exist"))
	}

	profile := ToResponseProfile(model)
	m, err := s.metadataSvc.GetMetadata(profileId)
	if err != nil {
		return nil, err
	}
	profile.Metadata = ConvertMetadataListToMap(m)
	return profile, nil
}

func (s *service) CreateProfile(reqProfile *RequestProfile) (*ResponseProfile, error) {
	// TODO: update validate before insert to database
	// validate storage object
	_, err := s.storageSvc.GetObject(reqProfile.StorageObjectID)
	if err != nil {
		return nil, err
	}
	// insert to database
	profileModel, err := s.repo.Insert(reqProfile)
	if err != nil {
		return nil, err
	}

	profile := ToResponseProfile(profileModel)

	if len(reqProfile.Metadata) > 0 {
		m, err := s.metadataSvc.CreateMetadata(int(profileModel.ID), reqProfile.Metadata)
		if err != nil {
			return nil, err
		}
		profile.Metadata = ConvertMetadataListToMap(m)
	}

	return ToResponseProfile(profileModel), err
}
