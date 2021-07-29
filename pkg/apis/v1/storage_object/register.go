package storage_object

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
	"io/ioutil"
	"net/http"
)

type StorageService interface {
	UploadToStorage(uploadedFile *File) (*ResponseObject, error)
	GetObjectByKey(objectKey string) (*File, error)
	DownloadObjectAsStream(ctx context.Context, objectKey string) (*storage.Reader, error)
	DownloadObject(objectKey string) (*File, error)
}
type register struct {
	storageSvc StorageService
}

func Register(g *echo.Group, db *database.DB) {
	res := register{
		storageSvc: NewService(db),
	}

	storageGroup := g.Group("/storages")
	storageGroup.GET("/:key/download", res.download)
	storageGroup.HEAD("/:key/download", res.download)
	storageGroup.POST("/upload", res.upload, middlewares.BasicAuth)
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

	resObj, err := r.storageSvc.UploadToStorage(uploadedFile)
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

	if ctx.Request().Method == http.MethodHead {
		object, err := r.storageSvc.GetObjectByKey(objectKey)
		if err != nil {
			return err
		}
		stream, err := r.storageSvc.DownloadObjectAsStream(ctx.Request().Context(), objectKey)
		ctx.Response().Header().Set(echo.HeaderContentLength, fmt.Sprintf("%d", stream.Attrs.Size))
		ctx.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", object.Filename))
		ctx.Response().Header().Del("Transfer-Encoding")
		return ctx.NoContent(http.StatusNoContent)
	}

	object, err := r.storageSvc.DownloadObject(objectKey)
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, object.ContentType, object.Content)
}
