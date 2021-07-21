package profile

import (
	"errors"
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"net/http"
	"reflect"
)

func ValidateRequiredField(field string, value interface{}) error {
	if value == nil || reflect.ValueOf(value).IsZero() {
		return cerrors.NewCError(http.StatusBadRequest, errors.New(fmt.Sprintf("%s is required", field)))
	}
	return nil
}
