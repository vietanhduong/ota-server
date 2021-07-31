package user

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/auth"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"net/http"
)

type Service interface {
	Login(rl *RequestLogin) (*User, error)
	GetUserInfo(email string) (*User, error)
}

type register struct {
	userSvc Service
	auth    *auth.Auth
}

func Register(g *echo.Group, db *mysql.DB, redis *redis.Client) {
	reg := register{
		userSvc: NewService(db),
		auth:    auth.NewAuth(redis),
	}

	authGroup := g.Group("/users")
	authGroup.GET("/me", reg.me, reg.auth.RequiredLogin())
	authGroup.POST("/login", reg.login)
	authGroup.POST("/refresh-token", reg.refreshToken)
	authGroup.POST("/logout", reg.logout)
}

func (r *register) me(ctx echo.Context) error {
	// get claims in context
	claims := r.auth.GetClaimsInContext(ctx)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	// find user by email
	user, err := r.userSvc.GetUserInfo(claims.User.Email)
	if err != nil {
		return err
	}
	// revoke token if user not found
	if user == nil {
		_ = r.auth.RevokeToken(claims.User.Email)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return ctx.JSON(http.StatusOK, user)
}

func (r *register) login(ctx echo.Context) error {
	// parse request login
	rl := new(RequestLogin)
	if err := ctx.Bind(rl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// validate request
	if err := ValidateRequestLogin(rl); err != nil {
		return err
	}

	user, err := r.userSvc.Login(rl)
	if err != nil {
		return err
	}

	authUser := &auth.User{
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}

	token, err := r.auth.GenerateToken(authUser)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, token)
}

func (r *register) refreshToken(ctx echo.Context) error {
	// parse authorization header
	authorization := ctx.Request().Header.Get("Authorization")
	refreshToken := r.auth.ExtractToken(authorization)
	if refreshToken == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	// parse refresh token
	claims, err := r.auth.ParseToken(refreshToken)
	if err != nil {
		return err
	}
	// just accept with token has type is `refresh`
	if claims.TokenType != auth.Refresh {
		return echo.NewHTTPError(http.StatusUnauthorized, "token invalid")
	}
	// find user by email to make sure user still active
	user, err := r.userSvc.GetUserInfo(claims.User.Email)
	if err != nil {
		return err
	}


	// get old token
	oldToken, err := r.auth.GetToken(claims.Id)
	if err != nil {
		return err
	}

	// if user not found, to be sure I have to
	// revoke the token
	if user == nil {
		_ = r.auth.RevokeToken(oldToken.AccessToken)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}


	// revoke old token
	if err := r.auth.RevokeToken(oldToken.AccessToken); err != nil {
		return err
	}

	authUser := &auth.User{
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}
	// regenerate token
	token, err := r.auth.GenerateToken(authUser)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, token)
}

func (r *register) logout(ctx echo.Context) error {
	// parse authorization header
	authorization := ctx.Request().Header.Get("Authorization")
	accessToken := r.auth.ExtractToken(authorization)
	if accessToken == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	claims, err := r.auth.ParseToken(accessToken)
	if err != nil {
		return err
	}

	if claims.TokenType != auth.Access {
		return echo.NewHTTPError(http.StatusUnauthorized, "token invalid")
	}

	if err := r.auth.RevokeToken(claims.User.Email); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
