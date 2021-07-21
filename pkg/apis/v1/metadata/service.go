package metadata

import (
	"github.com/vietanhduong/ota-server/pkg/database"
)

type service struct {
	repo *repository
}

func NewService(db *database.DB) *service {
	return &service{
		repo: NewRepository(db),
	}
}

func (s *service) CreateMetadata(profileId int, metadata map[string]string) ([]*Metadata, error) {
	var metadataList []*Metadata
	for k, v := range metadata {
		metadata := &Metadata{
			ProfileId: profileId,
			Key:       k,
			Value:     v,
		}
		metadataList = append(metadataList, metadata)
	}

	_, err := s.repo.InsertBulk(metadataList)
	if err != nil {
		return nil, err
	}
	return metadataList, nil
}

func (s *service) GetMetadata(profileId int) ([]*Metadata, error) {
	metadataModels, err := s.repo.FindByProfileId(uint(profileId))
	if err != nil {
		return nil, err
	}

	var result []*Metadata
	for _, mm := range metadataModels {
		result = append(result, ToMetadata(mm))
	}

	return result, nil
}
