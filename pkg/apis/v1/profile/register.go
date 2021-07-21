package profile

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
	"net/http"
)

type Service interface {
	CreateProfile(reqProfile *RequestProfile) (*ResponseProfile, error)
}

type register struct {
	profileSvc Service
}

func Register(g *echo.Group, db *database.DB) {
	res := register{
		profileSvc: NewService(db),
	}
	profileGroup := g.Group("/profiles")

	profileGroup.GET("/", res.home)
	profileGroup.POST("/ios", res.createProfile, middlewares.BasicAuth)
}

func (r *register) home(ctx echo.Context) error {
	return nil
}

func (r *register) createProfile(ctx echo.Context) error {
	var reqProfile *RequestProfile
	if err := ctx.Bind(&reqProfile); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := r.profileSvc.CreateProfile(reqProfile)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, res)
}
