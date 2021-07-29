package user

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"net/http"
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
	var rl *RequestLogin
	if err := ctx.Bind(&rl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
