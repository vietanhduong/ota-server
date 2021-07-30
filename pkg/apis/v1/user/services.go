package user

import (
	"github.com/vietanhduong/ota-server/pkg/auth"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"github.com/vietanhduong/ota-server/pkg/utils/crypto"
	"gopkg.in/errgo.v2/errors"
)

type service struct {
	userRepo *repository
	auth     *auth.Auth
}

func NewService(db *mysql.DB, redis *redis.Client) *service {
	return &service{
		userRepo: NewRepository(db),
		auth:     auth.NewAuth(redis),
	}
}

func (s *service) Login(rl *RequestLogin) (*auth.Token, error) {
	userModel, err := s.userRepo.FindByEmail(rl.Email, true)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if userModel == nil {
		return nil, cerrors.NotFound("user not found")
	}

	if userModel.Password != crypto.NewSHA256([]byte(rl.Password)) {
		return nil, cerrors.UnAuthorized("wrong password")
	}

	user := &auth.User{
		Email:       userModel.Email,
		DisplayName: userModel.DisplayName,
	}

	return s.auth.GenerateToken(user)
}

func (s *service) RefreshToken(refreshToken string) (*auth.Token, error) {
	// parse refresh token to token claims
	token, err := s.auth.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	// just accept with token has type is refresh
	if token.TokenType != auth.Refresh {
		return nil, cerrors.UnAuthorized("token invalid")
	}
	// make sure user is active
	userModel, err := s.userRepo.FindByEmail(token.User.Email, true)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// if user does not exist, to be sure we have to
	// revoke the token
	if userModel == nil {
		// revoke token if user does not exist
		_ = s.auth.RevokeToken(token.User.Email)
		return nil, cerrors.NotFound("user not found")
	}

	user := &auth.User{
		Email:       userModel.Email,
		DisplayName: userModel.DisplayName,
	}
	// revoke old token and regenerate new access token
	// and refresh token
	return s.auth.GenerateToken(user)
}

func (s *service) Logout(accessToken string) error {
	token, err := s.auth.ParseToken(accessToken)
	if err != nil {
		return errors.Wrap(err)
	}

	if token.TokenType != auth.Access {
		return cerrors.UnAuthorized("token invalid")
	}

	return s.auth.RevokeToken(token.User.Email)
}
