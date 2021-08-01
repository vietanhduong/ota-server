package profile

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/user"
	"github.com/vietanhduong/ota-server/pkg/auth"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
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

type UserService interface {
	GetUserInfo(email string) (*user.User, error)
}

type register struct {
	profileSvc Service
	userSvc    UserService
	auth       *auth.Auth
}

func Register(g *echo.Group, db *mysql.DB, redis *redis.Client) {
	reg := register{
		profileSvc: NewService(db),
		userSvc:    user.NewService(db),
		auth:       auth.NewAuth(redis),
	}
	profileGroup := g.Group("/profiles")

	profileGroup.GET("", reg.home, reg.auth.RequiredLogin())
	profileGroup.POST("/ios", reg.createProfile, reg.auth.RequiredLogin())
	profileGroup.GET("/ios/:id", reg.getProfile, reg.auth.RequiredLogin())
	profileGroup.GET("/ios/:id/manifest.plist", reg.getManifest, reg.auth.RequiredExchangeCode())
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

	// get uploaded user
	claims := r.auth.GetClaimsInContext(ctx)
	createdUser, err := r.userSvc.GetUserInfo(claims.User.Email)
	if err != nil {
		return err
	}
	// set created user id
	reqProfile.CreatedUserID = createdUser.Id

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
	code := ctx.QueryParam("code")

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
		"ipa_path":  fmt.Sprintf("%s/api/v1/storages/%s/download/%s?code=%s", Host, profile.StorageObject.ObjectKey, profile.StorageObject.Filename, code),
		"version":   profile.Version,
	}

	if v, found := profile.Metadata[AppIcon]; found {
		payload[AppIcon] = v
	}

	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	return ctx.Render(http.StatusOK, "manifest.plist", payload)
}
