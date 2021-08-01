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
		Id:          int(userModel.ID),
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
		return nil, cerrors.UnAuthorized("invalid credentials")
	}

	if userModel.Password != crypto.NewSHA256([]byte(rl.Password)) {
		return nil, cerrors.UnAuthorized("invalid credentials")
	}

	user := &User{
		Email:       userModel.Email,
		DisplayName: userModel.DisplayName,
		Active:      userModel.Active,
		CreatedAt:   userModel.CreatedAt,
	}

	return user, nil
}

func (s *service) GetUserByIds(userIds []int) (map[int]*User, error) {
	result := make(map[int]*User)
	if userIds == nil || len(userIds) == 0 {
		return result, nil
	}

	users, err := s.userRepo.FindByIds(userIds, true)
	if err != nil {
		return nil, err
	}

	for id, user := range users {
		result[id] = &User{
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Active:      user.Active,
			CreatedAt:   user.CreatedAt,
		}
	}

	return result, nil
}


func (s *service) GetUserByEmails(emails []string) (map[string]*User, error) {
	result := make(map[string]*User)
	if emails == nil || len(emails) == 0 {
		return result, nil
	}

	users, err := s.userRepo.FindByEmails(emails, true)
	if err != nil {
		return nil, err
	}

	for email, user := range users {
		result[email] = &User{
			Email:       user.Email,
			DisplayName: user.DisplayName,
			Active:      user.Active,
			CreatedAt:   user.CreatedAt,
		}
	}

	return result, nil
}
