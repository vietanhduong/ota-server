package metadata

import (
	"github.com/vietanhduong/ota-server/pkg/mysql"
)

type service struct {
	repo *repository
}

func NewService(db *mysql.DB) *service {
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

func (s *service) GetMetadataByListProfileId(profileIds []uint) (map[uint][]*Metadata, error) {
	metadataModels, err := s.repo.FindByListProfileId(profileIds)
	if err != nil {
		return nil, err
	}

	result := make(map[uint][]*Metadata)
	for _, mm := range metadataModels {
		if _, ok := result[mm.ProfileId]; !ok {
			var tmp []*Metadata
			result[mm.ProfileId] = tmp
		}
		result[mm.ProfileId] = append(result[mm.ProfileId], ToMetadata(mm))
	}

	return result, nil
}

