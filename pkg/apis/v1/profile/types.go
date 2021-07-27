package profile

import "github.com/vietanhduong/ota-server/pkg/database/models"

type RequestProfile struct {
	AppName          string            `json:"app_name"`
	BundleIdentifier string            `json:"bundle_id"`
	Version          string            `json:"version"`
	Build            int               `json:"build"`
	StorageObjectID  int               `json:"object_id"`
	Metadata         map[string]string `json:"metadata"`
}

type ResponseProfile struct {
	ProfileId        uint              `json:"profile_id"`
	AppName          string            `json:"app_name"`
	BundleIdentifier string            `json:"bundle_id"`
	Version          string            `json:"version"`
	Build            uint              `json:"build"`
	StorageObjectKey string            `json:"object_key"`
	Metadata         map[string]string `json:"metadata"`
}

// ToResponseProfile convert profile model to profile response object
func ToResponseProfile(model *models.Profile) *ResponseProfile {
	return &ResponseProfile{
		ProfileId:        model.ID,
		AppName:          model.AppName,
		BundleIdentifier: model.BundleIdentifier,
		Version:          model.Version,
		Build:            model.Build,
		StorageObjectKey: model.StorageObject.Key,
	}
}
