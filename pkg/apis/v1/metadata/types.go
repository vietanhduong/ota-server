package metadata

import "github.com/vietanhduong/ota-server/pkg/database/models"

type Metadata struct {
	ProfileId int
	Key       string
	Value     string
}

func ToMetadata(model *models.Metadata) *Metadata {
	return &Metadata{
		ProfileId: int(model.ProfileId),
		Key:       model.Key,
		Value:     model.Value,
	}
}
