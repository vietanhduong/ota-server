package profile

type RequestProfile struct {
	AppName          string `json:"app_name"`
	BundleIdentifier string `json:"bundle_id"`
	Version          string `json:"version"`
	Build            int    `json:"build"`
	StorageObjectID  int    `json:"object_id"`
}

type ResponseProfile struct {
	ProfileId        uint   `json:"profile_id"`
	AppName          string `json:"app_name"`
	BundleIdentifier string `json:"bundle_id"`
	Version          string `json:"version"`
	Build            uint    `json:"build"`
	StorageObjectID  uint    `json:"object_id"`
}
