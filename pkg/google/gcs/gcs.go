package gcs

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/utils/file"
	"google.golang.org/api/option"
	"io"
	"time"
)

type gcs struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGcs(credentialsPath, bucketName string) (*gcs, error) {
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

	// initial bucket
	bucket := storageClient.Bucket(bucketName)
	return &gcs{
		client: storageClient,
		bucket: bucket,
	}, nil
}

func (g *gcs) UploadObject(obj Object) error {
	// setup context with timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// upload an object with storage writer
	src := bytes.NewReader(obj.Content)

	w := g.bucket.Object(obj.OutputPath).NewWriter(ctx)

	if _, err := io.Copy(w, src); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

