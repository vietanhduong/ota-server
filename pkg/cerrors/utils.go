package cerrors

import (
	"github.com/labstack/echo/v4"
	"github.com/vietanhduong/ota-server/pkg/logger"
	"io"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		logger.Logger.Fatal(err)
	}
}

func GetRequestID(ctx echo.Context) string {
	id := ctx.Request().Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = ctx.Response().Header().Get(echo.HeaderXRequestID)
	}
	return id
}