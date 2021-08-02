package storage_object

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/google/gcs"
	"github.com/vietanhduong/ota-server/pkg/logger"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"gopkg.in/errgo.v2/errors"
	"net/http"
	"regexp"
	"time"
)

type service struct {
	repo    *repository
	storage *gcs.GoogleStorage
}

func NewService(db *mysql.DB) *service {
	googleCredentialsPath := env.GetEnvAsStringOrFallback("GOOGLE_CREDENTIALS", gcs.DefaultCredentialsPath)
	gcsBucket := env.GetEnvAsStringOrFallback("GCS_BUCKET", gcs.DefaultBucketName)

	_storage, err := gcs.NewGcs(googleCredentialsPath, gcsBucket)
	if err != nil {
		logger.Logger.Fatalf("initial google cloud storage failed with err: %v", err)
	}

	return &service{
		repo:    NewRepository(db),
		storage: _storage,
	}
}

func (s *service) UploadToStorage(ctx context.Context, uploadedFile *File) (*ResponseObject, error) {
	// validate file extension
	if err := ValidateExtension(uploadedFile.Filename); err != nil {
		return nil, err
	}

	// normalize file name
	m := regexp.MustCompile("[^0-9a-zA-Z.]+")
	uploadedFile.Filename = m.ReplaceAllLiteralString(uploadedFile.Filename, "_")

	// generate object key
	uploadedFile.Key = shortuuid.New()

	// generate abs path
	// 2006/01/02/Cekw67uyMpBGZLRP2HFVbe_build.ipa
	uploadedFile.AbsPath = fmt.Sprintf("%s/%s_%s", time.Now().Format("2006/01/02"), uploadedFile.Key, uploadedFile.Filename)

	// upload to GCS
	obj := &gcs.Object{
		Content:    uploadedFile.Content,
		OutputPath: uploadedFile.AbsPath,
	}
	if err := s.storage.UploadObject(ctx, obj); err != nil {
		return nil, err
	}
	model, err := s.repo.Insert(uploadedFile)
	if err != nil {
		return nil, err
	}

	return &ResponseObject{
		ObjectId:  model.ID,
		ObjectKey: model.Key,
		AbsPath:   uploadedFile.AbsPath,
		Filename:  model.Name,
	}, nil
}

func (s *service) GetObjectById(objectId int) (*File, error) {
	// verify object id
	object, err := s.repo.FindById(uint(objectId))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if object == nil {
		return nil, cerrors.NewCError(http.StatusNotFound, "object does not exist")
	}
	return &File{
		Key:         object.Key,
		Filename:    object.Name,
		ContentType: object.ContentType,
		AbsPath:     object.Path,
		CreatedAt:   object.CreatedAt,
		UpdatedAt:   object.UpdatedAt,
	}, nil
}

func (s *service) GetObjectByKey(objectKey string) (*File, error) {
	// verify object id
	object, err := s.repo.FindByKey(objectKey)
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, cerrors.NewCError(http.StatusNotFound, "object does not exist")
	}
	return &File{
		Filename:    object.Name,
		ContentType: object.ContentType,
		AbsPath:     object.Path,
		CreatedAt:   object.CreatedAt,
		UpdatedAt:   object.UpdatedAt,
	}, nil
}

func (s *service) GetObjectsByKeys(objectKeys []string) (map[string]*File, error) {
	result := make(map[string]*File)

	if objectKeys == nil || len(objectKeys) == 0 {
		return result, nil
	}

	objects, err := s.repo.FindByKeys(objectKeys)
	if err != nil {
		return nil, err
	}

	for _, o := range objects {
		obj := &File{
			Filename:    o.Name,
			ContentType: o.ContentType,
			AbsPath:     o.Path,
			CreatedAt:   o.CreatedAt,
			UpdatedAt:   o.UpdatedAt,
		}
		result[o.Key] = obj
	}

	return result, nil
}

func (s *service) DownloadObject(ctx context.Context, objectKey string) (*File, error) {
	// verify object id
	object, err := s.repo.FindByKey(objectKey)
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, cerrors.NewCError(http.StatusNotFound, "object does not exist")
	}

	content, err := s.storage.ReadObject(ctx, object.Path)
	if err != nil {
		return nil, err
	}
	return &File{
		Filename:    object.Name,
		ContentType: object.ContentType,
		Content:     content,
	}, nil
}

func (s *service) DownloadObjectAsStream(ctx context.Context, objectKey string) (*storage.Reader, error) {
	// verify object id
	object, err := s.repo.FindByKey(objectKey)
	if err != nil {
		return nil, err
	}

	if object == nil {
		return nil, cerrors.NewCError(http.StatusNotFound, "object does not exist")
	}

	stream, err := s.storage.ReadObjectAsStream(ctx, object.Path)
	return stream, err
}
