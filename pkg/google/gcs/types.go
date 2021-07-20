package gcs

type Object struct {
	Content  []byte
	// OutputPath could be absolute path on GCS
	OutputPath string
}
