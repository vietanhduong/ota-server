package storage_object

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/auth"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"io/ioutil"
	"net/http"
	"strconv"
)

type StorageService interface {
	UploadToStorage(ctx context.Context, uploadedFile *File) (*ResponseObject, error)
	GetObjectByKey(objectKey string) (*File, error)
	DownloadObjectAsStream(ctx context.Context, objectKey string) (*storage.Reader, error)
	DownloadObject(ctx context.Context, objectKey string) (*File, error)
}
type register struct {
	storageSvc StorageService
	auth       *auth.Auth
}

func Register(g *echo.Group, db *mysql.DB, redis *redis.Client) {
	reg := register{
		storageSvc: NewService(db),
		auth:       auth.NewAuth(redis),
	}

	storageGroup := g.Group("/storages")
	storageGroup.GET("/:key/download/*", reg.download, reg.auth.RequiredExchangeCode())
	storageGroup.HEAD("/:key/download/*", reg.download, reg.auth.RequiredExchangeCode())
	storageGroup.POST("/upload", reg.upload, reg.auth.RequiredLogin())
}

func (r *register) upload(ctx echo.Context) error {
	// receive upload file
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}

	defer cerrors.Close(f)

	// read content
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	uploadedFile := &File{
		Filename:    file.Filename,
		Content:     content,
		ContentType: file.Header.Get("Content-Type"),
	}

	resObj, err := r.storageSvc.UploadToStorage(ctx.Request().Context(), uploadedFile)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, resObj)
}

func (r *register) download(ctx echo.Context) error {
	objectKey := ctx.Param("key")
	if objectKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid object key")
	}

	object, err := r.storageSvc.GetObjectByKey(objectKey)
	if err != nil {
		return err
	}

	stream, err := r.storageSvc.DownloadObjectAsStream(ctx.Request().Context(), objectKey)
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Accept-Ranges", "bytes")
	ctx.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(int(stream.Attrs.Size)))
	ctx.Response().Header().Set(echo.HeaderContentType, object.ContentType)

	if ctx.Request().Method == http.MethodHead {
		return ctx.NoContent(http.StatusOK)
	}

	return ctx.Stream(http.StatusOK, object.ContentType, stream)
}
