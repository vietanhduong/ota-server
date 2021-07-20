package storage_object

type UploadedFile struct {
	Filename    string
	Content     []byte
	ContentType string
	AbsPath     string
}

type ResponseObject struct {
	ObjectId uint   `json:"object_id"`
	Filename string `json:"filename"`
	AbsPath  string `json:"abs_path"`
}
