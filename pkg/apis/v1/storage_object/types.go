package storage_object

import "time"

type File struct {
	Filename    string    `json:"filename"`
	Content     []byte    `json:"-"`
	ContentType string    `json:"content_type"`
	AbsPath     string    `json:"abs_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResponseObject struct {
	ObjectId uint   `json:"object_id"`
	Filename string `json:"filename"`
	AbsPath  string `json:"abs_path"`
}
