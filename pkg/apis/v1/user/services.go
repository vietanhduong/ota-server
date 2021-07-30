package user

import (
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/utils/crypto"
	"gopkg.in/errgo.v2/errors"
)

type service struct {
	userRepo *repository
}

func NewService(db *mysql.DB) *service {
	return &service{
		userRepo: NewRepository(db),
	}
}

func (s *service) GetUserInfo(email string) (*User, error) {
	userModel, err := s.userRepo.FindByEmail(email, true)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if userModel == nil {
		return nil, cerrors.NotFound("user not found")
	}

	user := &User{
		Email:       userModel.Email,
		DisplayName: userModel.DisplayName,
		Active:      userModel.Active,
		CreatedAt:   userModel.CreatedAt,
	}
	return user, nil
}

func (s *service) Login(rl *RequestLogin) (*User, error) {
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

	user := &User{
		Email:       userModel.Email,
		DisplayName: userModel.DisplayName,
		Active:      userModel.Active,
		CreatedAt:   userModel.CreatedAt,
	}

	return user, nil
}
