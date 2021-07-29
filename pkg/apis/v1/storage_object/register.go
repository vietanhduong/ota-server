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
	UploadToStorage(ctx context.Context, uploadedFile *File) (*ResponseObject, error)
	GetObjectByKey(objectKey string) (*File, error)
	DownloadObjectAsStream(ctx context.Context, objectKey string) (*storage.Reader, error)
	DownloadObject(ctx context.Context, objectKey string) (*File, error)
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
	ctx.Response().Header().Set(echo.HeaderContentLength, fmt.Sprintf("%d", stream.Attrs.Size))
	ctx.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", object.Filename))
	ctx.Response().Header().Del("Transfer-Encoding")

	if ctx.Request().Method == http.MethodHead {
		return ctx.NoContent(http.StatusNoContent)
	}

	return ctx.Stream(http.StatusOK, object.ContentType, stream)
}
