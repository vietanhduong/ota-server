package profile

import "github.com/vietanhduong/ota-server/pkg/mysql/models"

type RequestProfile struct {
	AppName          string            `json:"app_name"`
	BundleIdentifier string            `json:"bundle_id"`
	Version          string            `json:"version"`
	Build            int               `json:"build"`
	StorageObjectID  int               `json:"object_id"`
	Metadata         map[string]string `json:"metadata"`
}

type ResponseProfile struct {
	ProfileId        uint                   `json:"profile_id"`
	AppName          string                 `json:"app_name"`
	BundleIdentifier string                 `json:"bundle_id"`
	Version          string                 `json:"version"`
	Build            uint                   `json:"build"`
	Metadata         map[string]string      `json:"metadata"`
	StorageObject    *StorageObjectResponse `json:"storage_object"`
}

type StorageObjectResponse struct {
	ObjectKey string `json:"object_key"`
	Filename  string `json:"filename"`
}

// ToResponseProfile convert profile model to profile response object
func ToResponseProfile(model *models.Profile) *ResponseProfile {
	return &ResponseProfile{
		ProfileId:        model.ID,
		AppName:          model.AppName,
		BundleIdentifier: model.BundleIdentifier,
		Version:          model.Version,
		Build:            model.Build,
		StorageObject: &StorageObjectResponse{
			ObjectKey: model.StorageObject.Key,
			Filename:  model.StorageObject.Name,
		},
	}
}
