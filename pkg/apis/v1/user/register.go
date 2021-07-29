package auth

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"net/http"
	"regexp"
	"strings"
)

type Service interface {
	Login(ctx context.Context, idToken string) (*Token, error)
}

type register struct {
	authSvc Service
}

func Register(g *echo.Group, _ *mysql.DB) {
	res := register{
		authSvc: NewService(),
	}

	authGroup := g.Group("/auth")
	authGroup.POST("/login", res.login)
}

func (r *register) login(ctx echo.Context) error {
	authorizationHeader := ctx.Request().Header.Get("Authorization")
	idToken := parseAuthorizationHeader(authorizationHeader)
	if idToken == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "token invalid")
	}

	_, err := r.authSvc.Login(ctx.Request().Context(), idToken)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, "ok")
}

func parseAuthorizationHeader(authorizationHeader string) string {
	if authorizationHeader == "" {
		return ""
	}

	var validToken = regexp.MustCompile(`^((?i)bearer|(?i)token)\s`)
	if validToken.MatchString(authorizationHeader) {
		token := validToken.ReplaceAllString(authorizationHeader, "")
		return strings.Trim(token, "")
	}
	return ""

}
