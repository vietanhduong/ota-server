package profile

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/database"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"net/http"
	"strconv"
)

var Host = env.GetEnvAsStringOrFallback("HOST", "https://ota.anhdv.dev")

type Service interface {
	CreateProfile(reqProfile *RequestProfile) (*ResponseProfile, error)
	GetProfile(profileId int) (*ResponseProfile, error)
	GetProfiles() ([]*ResponseProfile, error)
}

type register struct {
	profileSvc Service
}

func Register(g *echo.Group, db *database.DB) {
	res := register{
		profileSvc: NewService(db),
	}
	profileGroup := g.Group("/profiles")

	profileGroup.GET("", res.home)
	profileGroup.POST("/ios", res.createProfile, middlewares.BasicAuth)
	profileGroup.GET("/ios/:id", res.getProfile)
	profileGroup.GET("/ios/:id/manifest.plist", res.getManifest)
}

func (r *register) home(ctx echo.Context) error {
	profiles, err := r.profileSvc.GetProfiles()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, profiles)
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

func (r *register) getProfile(ctx echo.Context) error {
	reqProfileId := ctx.Param("id")
	profileId, err := strconv.Atoi(reqProfileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid profile id")
	}

	profile, err := r.profileSvc.GetProfile(profileId)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (r *register) getManifest(ctx echo.Context) error {
	reqProfileId := ctx.Param("id")
	profileId, err := strconv.Atoi(reqProfileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid profile id")
	}

	profile, err := r.profileSvc.GetProfile(profileId)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"app_name":  profile.AppName,
		"bundle_id": profile.BundleIdentifier,
		// ipa_path could be download api
		"ipa_path": fmt.Sprintf("%s/api/v1/storages/%s/download", Host, profile.StorageObjectKey),
		"version":  profile.Version,
	}
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	return ctx.Render(http.StatusOK, "manifest.plist", payload)
}
