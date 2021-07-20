package middlewares

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

var RootUser = os.Getenv("ROOT_USER")
var RootSecret = os.Getenv("ROOT_SECRET")

func constantTimeCompare(x, y string) bool {
	return subtle.ConstantTimeCompare([]byte(x), []byte(y)) == 1
}

var BasicAuth = middleware.BasicAuth(func(username, secret string, ctx echo.Context) (bool, error) {
	return constantTimeCompare(username, RootUser) && constantTimeCompare(secret, RootSecret), nil
})
