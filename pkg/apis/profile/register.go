package profile

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/database"
)

type resource struct{}

func Register(g *echo.Group, db *database.DB) {
	res := resource{}
	g.GET("/", res.home)
}

func (res *resource) home(ctx echo.Context) error {
	return nil
}
