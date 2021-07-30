package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"golang.org/x/exp/rand"
	"gopkg.in/errgo.v2/errors"
	"regexp"
	"strings"
	"time"
)

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

	// generate exchange code
	exchangeCode, err := a.GenerateExchangeCode()
	if err != nil {
		return nil, err
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExchangeCode: exchangeCode,
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

func (a *Auth) RevokeToken(email string) error {
	// if email does not exist in redis
	// stop revoke token
	if a.redis.Exists(email).Val() != 1 {
		return nil
	}

	// retrieve token payload in redis
	var token *Token
	if err := a.redis.GetValue(email, &token); err != nil {
		return err
	}

	// revoke token
	if err := a.redis.Del(email).Err(); err != nil {
		return err
	}

	// revoke exchange code
	if err := a.RevokeExchangeCode(token.ExchangeCode); err != nil {
		return err
	}
	return nil
}

func (a *Auth) ExtractToken(authorization string) string {
	if authorization == "" {
		return ""
	}
	var validToken = regexp.MustCompile(`^((?i)bearer|(?i)token|(?i)jwt)\s`)
	if validToken.MatchString(authorization) {
		token := validToken.ReplaceAllString(authorization, "")
		return strings.Trim(token, "")
	}
	return ""
}

func (a *Auth) GetClaimsInContext(ctx echo.Context) *TokenClaims {
	jwtToken := ctx.Get(CtxKey)
	if jwtToken != nil {
		return ctx.Get(CtxKey).(*TokenClaims)
	}
	return nil
}

func (a *Auth) GenerateExchangeCode() (string, error) {
	var base = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var defaultLength = 28
	code, err := gonanoid.Generate(base, defaultLength)
	if err != nil {
		return "", err
	}

	// store into redis with timeout equal access token
	if err := a.redis.StoreWithTTL(fmt.Sprint("exchange_", code), true, AccessTokenValidTime); err != nil {
		return "", err
	}
	return code, nil
}

func (a *Auth) IsExchangeCodeExist(code string) bool {
	return a.redis.Exists(fmt.Sprint("exchange_", code)).Val() == 1
}

func (a *Auth) RevokeExchangeCode(code string) error {
	if code == "" {
		return nil
	}
	if !a.IsExchangeCodeExist(code) {
		return nil
	}

	return a.redis.Del(fmt.Sprint("exchange_", code)).Err()
}

func (a *Auth) RequiredLogin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// extract access token from request header
			token := a.ExtractToken(ctx.Request().Header.Get("Authorization"))
			if token == "" {
				return cerrors.UnAuthorized("unauthorized")
			}
			// parse to claims
			claims, err := a.ParseToken(token)
			if err != nil {
				return errors.Wrap(err)
			}
			// just access with token has type is `access`
			if claims.TokenType != Access {
				return cerrors.UnAuthorized("token invalid")
			}
			// set claims into context
			ctx.Set(CtxKey, claims)
			return next(ctx)
		}
	}
}

func (a *Auth) RequiredExchangeCode() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			code := ctx.QueryParam("code")
			if code == "" {
				return cerrors.UnAuthorized("missing exchange code")
			}

			if !a.IsExchangeCodeExist(code) {
				return cerrors.UnAuthorized("exchange code invalid")
			}

			return next(ctx)
		}
	}
}
