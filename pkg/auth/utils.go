package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
)

func ProcessJwtError(err error) error {
	validateErr, ok := err.(*jwt.ValidationError)
	if !ok {
		return err
	}
	if validateErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
		return cerrors.UnAuthorized("token expired")
	}
	return cerrors.UnAuthorized("token invalid")
}
