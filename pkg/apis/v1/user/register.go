package user

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"net/http"
	"regexp"
	"strings"
)

type Service interface {
	Login(rl *RequestLogin) (*Token, error)
}

type register struct {
	userSvc Service
}

func Register(g *echo.Group, db *mysql.DB, redis *redis.Client) {
	res := register{
		userSvc: NewService(db, redis),
	}

	authGroup := g.Group("/users")
	authGroup.POST("/login", res.login)
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

	token, err := r.userSvc.Login(rl)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, token)
}

func (r *register) refreshToken(ctx echo.Context) error {
	// parse authorization header
	authorization := ctx.Request().Header.Get("Authorization")
	token := extractToken(authorization)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	return nil
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
