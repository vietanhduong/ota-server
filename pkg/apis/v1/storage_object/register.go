package storage_object

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
	"io/ioutil"
	"net/http"
	"strconv"
)

type StorageService interface {
	UploadToStorage(uploadedFile *File) (*ResponseObject, error)
	DownloadObject(objectId int) (*File, error)
	GetObject(objectId int) (*File, error)
}
type register struct {
	storageSvc StorageService
}

func Register(g *echo.Group, db *database.DB) {
	res := register{
		storageSvc: NewService(db),
	}

	storageGroup := g.Group("/storages")
	storageGroup.GET("/:id/download/*", res.download)
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
	reqObjId := ctx.Param("id")
	objectId, err := strconv.Atoi(reqObjId)
	if err != nil {
		return cerrors.NewCError(http.StatusBadRequest, errors.New("invalid object id"))
	}

	file, err := r.storageSvc.DownloadObject(objectId)

	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusOK, file.ContentType, file.Content)
}
