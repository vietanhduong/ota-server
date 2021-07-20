package storage_object

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
)

type register struct {
}

func Register(g *echo.Group, db *database.DB) {
	res := register{}
	storageGroup := g.Group("/storages")

	storageGroup.GET("/", res.home)
	storageGroup.POST("/upload", res.upload, middlewares.BasicAuth)
}

func (r *register) home(ctx echo.Context) error {
	return nil
}

func (r *register) upload(ctx echo.Context) error {

	return nil
}
