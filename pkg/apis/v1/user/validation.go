package user

import (
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"reflect"
)

func ValidateRequestLogin(rl *RequestLogin) error {
	// validate request object
	if rl == nil {
		return cerrors.BadRequest("invalid request")
	}
	// validate input email
	if err := ValidateRequiredField("email", rl.Email); err != nil {
		return err
	}
	// validate input password
	return ValidateRequiredField("password", rl.Password)
}

func ValidateRequiredField(fieldName string, value interface{}) error {
	if value == nil || reflect.ValueOf(value).IsZero() {
		return cerrors.BadRequest(fmt.Sprint(fieldName, " is required"))
	}
	return nil
}
