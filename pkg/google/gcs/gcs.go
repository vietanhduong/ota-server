package gcs

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/utils/file"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"time"
)

type GoogleStorage struct {
	client *storage.Client
	bucket string
}

func NewGcs(credentialsPath, bucket string) (*GoogleStorage, error) {
	// verify credentials path
	if !file.IsExist(credentialsPath) {
		return nil, errors.New(fmt.Sprintf("%s does not exist", credentialsPath))
	}
	// initial storage client
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}

	return &GoogleStorage{
		client: storageClient,
		bucket: bucket,
	}, nil
}

func (g *GoogleStorage) UploadObject(obj *Object) error {
	// setup context with timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// upload an object with storage writer
	src := bytes.NewReader(obj.Content)

	w := g.client.Bucket(g.bucket).Object(obj.OutputPath).NewWriter(ctx)

	if _, err := io.Copy(w, src); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (g *GoogleStorage) ReadObject(object string) ([]byte, error) {
	// setup context with timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	rc, err := g.client.Bucket(g.bucket).Object(object).NewReader(ctx)
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
