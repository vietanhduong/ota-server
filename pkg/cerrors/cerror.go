package cerrors

import (
	"fmt"
	"net/http"
)

type CError struct {
	Code int
	Err  interface{}
}

func (c *CError) Error() string {
	return fmt.Sprintf("Error: %v", c.Err)
}

func NewCError(code int, err interface{}) *CError {
	return &CError{
		Code: code,
		Err:  err,
	}
}

func NotFound(err interface{}) *CError {
	return NewCError(http.StatusNotFound, err)
}

func UnAuthorized(err interface{}) *CError {
	return NewCError(http.StatusUnauthorized, err)
}

func BadRequest(err interface{}) *CError {
	return NewCError(http.StatusBadRequest, err)
}
