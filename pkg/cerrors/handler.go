// Package cerrors customize error
package cerrors

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func HTTPErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	var msg interface{}
	msg = err.Error()

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}

	if ce, ok := err.(*CError); ok {
		code = ce.Code
		msg = ce.Error()
	}

	response := &Response{
		Code:    code,
		Message: msg,
	}

	if err := ctx.JSON(code, response); err != nil {
		ctx.Logger().Error(err)
	}
	ctx.Logger().Error(err)
}
