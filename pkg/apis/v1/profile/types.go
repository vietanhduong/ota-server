package profile

import "github.com/vietanhduong/ota-server/pkg/mysql/models"

type RequestProfile struct {
	AppName          string            `json:"app_name"`
	BundleIdentifier string            `json:"bundle_id"`
	Version          string            `json:"version"`
	Build            int               `json:"build"`
	StorageObjectID  int               `json:"object_id"`
	Metadata         map[string]string `json:"metadata"`
	CreatedUserID    int               `json:"-"`
}

type ResponseProfile struct {
	ProfileId        uint                   `json:"profile_id"`
	AppName          string                 `json:"app_name"`
	BundleIdentifier string                 `json:"bundle_id"`
	Version          string                 `json:"version"`
	Build            uint                   `json:"build"`
	Metadata         map[string]string      `json:"metadata"`
	StorageObject    *StorageObjectResponse `json:"storage_object"`
	CreatedBy        *CreatedByResponse     `json:"created_by"`
}

type StorageObjectResponse struct {
	ObjectKey string `json:"object_key"`
	Filename  string `json:"filename"`
}

type CreatedByResponse struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

// ToResponseProfile convert profile model to profile response object
func ToResponseProfile(p *models.Profile) *ResponseProfile {
	obj := &ResponseProfile{
		ProfileId:        p.ID,
		AppName:          p.AppName,
		BundleIdentifier: p.BundleIdentifier,
		Version:          p.Version,
		Build:            p.Build,
	}

	if p.StorageObject != nil {
		obj.StorageObject = &StorageObjectResponse{
			ObjectKey: p.StorageObject.Key,
			Filename:  p.StorageObject.Name,
		}
	}

	if p.User != nil {
		obj.CreatedBy = &CreatedByResponse{
			Email:       p.User.Email,
			DisplayName: p.User.DisplayName,
		}
	}

	return obj
}
