package user

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/auth"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"net/http"
	"regexp"
	"strings"
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
	res := register{
		userSvc: NewService(db),
		auth:    auth.NewAuth(redis),
	}

	authGroup := g.Group("/users")
	authGroup.POST("/login", res.login)
	authGroup.POST("/refresh-token", res.refreshToken)
	authGroup.POST("/logout", res.logout)
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
	refreshToken := extractToken(authorization)
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
	// if user not found, to be sure I have to
	// revoke the token
	if user == nil {
		_ = r.auth.RevokeToken(claims.User.Email)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	authUser := &auth.User{
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}
	// revoke both access token and refresh token,
	// and regenerate
	token, err := r.auth.GenerateToken(authUser)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, token)
}

func (r *register) logout(ctx echo.Context) error {
	// parse authorization header
	authorization := ctx.Request().Header.Get("Authorization")
	accessToken := extractToken(authorization)
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

func extractToken(authorization string) string {
	if authorization == "" {
		return ""
	}
	var validToken = regexp.MustCompile(`^((?i)bearer|(?i)token|(?i)jwt)\s`)
	if validToken.MatchString(authorization) {
		token := validToken.ReplaceAllString(authorization, "")
		return strings.Trim(token, "")
	}
	return ""
}
