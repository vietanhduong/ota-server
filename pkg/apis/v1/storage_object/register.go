package storage_object

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type StorageService interface {
	UploadToGoogleStorage(uploadedFile *UploadedFile) (*ResponseObject, error)
}

type register struct {
	storageSvc StorageService
}

func Register(g *echo.Group, db *database.DB) {
	res := register{
		storageSvc: NewService(db),
	}

	storageGroup := g.Group("/storages")
	storageGroup.GET("/", res.home)
	storageGroup.POST("/upload", res.upload, middlewares.BasicAuth)
}

func (r *register) home(ctx echo.Context) error {
	return nil
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

	defer func(f multipart.File) {
		_ = f.Close()
	}(f)

	// read content
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	uploadedFile := &UploadedFile{
		Filename:    file.Filename,
		Content:     content,
		ContentType: file.Header.Get("Content-Type"),
	}

	resObj, err := r.storageSvc.UploadToGoogleStorage(uploadedFile)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, resObj)
}
