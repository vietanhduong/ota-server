package user

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"github.com/vietanhduong/ota-server/pkg/utils/crypto"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"gopkg.in/errgo.v2/errors"
	"math/rand"
	"time"
)

const RefreshTokenValidTime = 12 * time.Hour
const AccessTokenValidTime = time.Hour
const Refresh = "refresh"
const Access = "access"

type service struct {
	userRepo *repository
	redis    *redis.Client
}

var secret = env.GetEnvAsStringOrFallback("SECRET", "some-thing-very-secret")

func NewService(db *mysql.DB, redis *redis.Client) *service {
	return &service{
		userRepo: NewRepository(db),
		redis:    redis,
	}
}

func (s *service) Login(rl *RequestLogin) (*Token, error) {
	userModel, err := s.userRepo.FindByEmail(rl.Email)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if userModel == nil {
		return nil, cerrors.NotFound("userModel not found")
	}

	if userModel.Password != crypto.NewSHA256([]byte(rl.Password)) {
		return nil, cerrors.UnAuthorized("wrong password")
	}

	user := &User{
		Email:       userModel.Email,
		DisplayName: userModel.DisplayName,
	}

	return s.GenerateToken(user)
}

func (s *service) GenerateToken(user *User) (*Token, error) {
	// generate Jwt Token Id
	jti, err := s.GenerateJti()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// generate refresh token
	refreshToken, err := s.GenerateRefreshToken(user, jti)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// generate access token
	accessToken, err := s.GenerateAccessToken(user, jti)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// save to redis
	if err := s.redis.StoreWithTTL(jti, *token, RefreshTokenValidTime); err != nil {
		return nil, errors.Wrap(err)
	}

	return token, nil
}

func (s *service) GenerateAccessToken(user *User, jti string) (string, error) {
	now := time.Now()
	accessTokenExp := now.Add(AccessTokenValidTime).Unix()

	claims := &TokenClaims{
		*user,
		Access,
		jwt.StandardClaims{
			Id:        jti,
			IssuedAt:  now.Unix(),
			ExpiresAt: accessTokenExp,
		},
	}
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessJwt.SignedString([]byte(secret))
	return token, err
}

func (s *service) GenerateRefreshToken(user *User, jti string) (string, error) {
	now := time.Now()
	refreshTokenExp := now.Add(RefreshTokenValidTime).Unix()

	claims := &TokenClaims{
		*user,
		Refresh,
		jwt.StandardClaims{
			Id:        jti,
			IssuedAt:  now.Unix(),
			ExpiresAt: refreshTokenExp,
		},
	}
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := refreshJwt.SignedString([]byte(secret))
	return token, err
}

func (s *service) GenerateJti() (string, error) {
	var length = 64
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
