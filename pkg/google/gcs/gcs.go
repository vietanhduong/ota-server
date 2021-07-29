package gcs

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/utils/file"
	"io"
	"io/ioutil"
)

type GoogleStorage struct {
	bucket         string
	credentialPath string
}

func NewGcs(credentialsPath, bucket string) (*GoogleStorage, error) {
	// verify credentials path
	if !file.IsExist(credentialsPath) {
		return nil, errors.New(fmt.Sprintf("%s does not exist", credentialsPath))
	}
	return &GoogleStorage{
		bucket:         bucket,
		credentialPath: credentialsPath,
	}, nil
}

func (g *GoogleStorage) UploadObject(ctx context.Context, obj *Object) error {
	// upload an object with storage writer
	src := bytes.NewReader(obj.Content)
	// init gcs client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	// init bucket
	w := client.Bucket(g.bucket).Object(obj.OutputPath).NewWriter(ctx)

	if _, err := io.Copy(w, src); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (g *GoogleStorage) ReadObject(ctx context.Context, object string) ([]byte, error) {
	// init gcs client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// init bucket
	rc, err := client.Bucket(g.bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer cerrors.Close(rc)

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (g *GoogleStorage) ReadObjectAsStream(ctx context.Context, object string) (*storage.Reader, error) {
	// init gcs client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	rc, err := client.Bucket(g.bucket).Object(object).NewReader(ctx)
	return rc, err
}
