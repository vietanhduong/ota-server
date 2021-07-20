package profile

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/database"
)

type register struct{}

func Register(g *echo.Group, db *database.DB) {
	res := register{}
	profileGroup := g.Group("/profiles")

	profileGroup.GET("/", res.home)
}

func (r *register) home(ctx echo.Context) error {
	return nil
}