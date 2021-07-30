package auth

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"golang.org/x/exp/rand"
	"gopkg.in/errgo.v2/errors"
	"time"
)

const RefreshTokenValidTime = 12 * time.Hour
const AccessTokenValidTime = time.Hour
const Refresh = "refresh"
const Access = "access"

type Auth struct {
	redis  *redis.Client
	secret string
}

func NewAuth(redis *redis.Client) *Auth {
	return &Auth{
		redis:  redis,
		secret: env.GetEnvAsStringOrFallback("SECRET", "some-thing-very-secret"),
	}
}

func (a *Auth) ParseToken(inputToken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, ProcessJwtError(err)
	}

	if a.IsTokenRevoked(token) {
		return nil, cerrors.UnAuthorized("token invalid")
	}

	return token.Claims.(*TokenClaims), nil
}

func (a *Auth) GenerateToken(user *User) (*Token, error) {
	// generate Jwt Token Id
	jti, err := a.GenerateJti()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// generate refresh token
	refreshToken, err := a.GenerateRefreshToken(user, jti)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// generate access token
	accessToken, err := a.GenerateAccessToken(user, jti)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// save to redis
	if err := a.redis.StoreWithTTL(user.Email, *token, RefreshTokenValidTime); err != nil {
		return nil, errors.Wrap(err)
	}

	return token, nil
}

func (a *Auth) GenerateAccessToken(user *User, jti string) (string, error) {
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
	token, err := accessJwt.SignedString([]byte(a.secret))
	return token, err
}

func (a *Auth) GenerateRefreshToken(user *User, jti string) (string, error) {
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
	token, err := refreshJwt.SignedString([]byte(a.secret))
	return token, err
}

func (a *Auth) GenerateJti() (string, error) {
	var length = 64
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *Auth) IsTokenRevoked(token *jwt.Token) bool {
	claims := token.Claims.(*TokenClaims)
	result := a.redis.Exists(claims.User.Email).Val()
	return result == 0
}
