package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/logger"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
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

// ParseToken parse and validate input token
func (a *Auth) ParseToken(inputToken string) (*TokenClaims, error) {
	// try to parse input token
	token, err := jwt.ParseWithClaims(inputToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, ProcessJwtError(err)
	}

	// verify token is revoked or not
	if a.IsTokenRevoked(token) {
		return nil, cerrors.UnAuthorized("token invalid")
	}

	return token.Claims.(*TokenClaims), nil
}

// GetToken get token from redis by input jti
func (a *Auth) GetToken(jti string) (*Token, error) {
	if a.redis.Exists(jti).Val() != 1 {
		logger.Logger.Warnf("JTI: %s does not exist", jti)
		return nil, cerrors.BadRequest("token invalid")
	}

	var token *Token
	if err := a.redis.GetValue(jti, &token); err != nil {
		return nil, errors.Wrap(err)
	}

	return token, nil
}

// GenerateToken generate token
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
	if err := a.redis.StoreWithTTL(jti, *token, RefreshTokenValidTime); err != nil {
		return nil, errors.Wrap(err)
	}

	return token, nil
}

// GenerateAccessToken generate access token
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

// GenerateRefreshToken generate refresh token
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

// GenerateJti generate jwt token id
func (a *Auth) GenerateJti() (string, error) {
	var base = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var defaultLength = 28
	code, err := gonanoid.Generate(base, defaultLength)
	if err != nil {
		return "", err
	}
	return code, nil
}

// IsTokenRevoked check token is revoked or not
func (a *Auth) IsTokenRevoked(token *jwt.Token) bool {
	claims := token.Claims.(*TokenClaims)
	result := a.redis.Exists(claims.Id).Val()
	return result == 0
}

// RevokeToken remove token from redis
func (a *Auth) RevokeToken(accessToken string) error {
	// if email does not exist in redis
	// stop revoke token
	claims, err := a.ParseToken(accessToken)
	if err != nil {
		return errors.Wrap(err)
	}

	if claims.TokenType != Access {
		return cerrors.UnAuthorized("token invalid")
	}

	// stop revoke token if not found jti
	if a.redis.Exists(claims.Id).Val() != 1 {
		return nil
	}

	var token *Token
	if err := a.redis.GetValue(claims.Id, &token); err != nil {
		return errors.Wrap(err)
	}

	// revoke token
	if err := a.redis.Del(claims.Id).Err(); err != nil {
		return err
	}

	// revoke exchange code
	if err := a.RevokeExchangeCode(token.ExchangeCode); err != nil {
		return err
	}
	return nil
}

// ExtractToken extract token request header
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

// GetClaimsInContext get claim was stored in context
// when user logged in
func (a *Auth) GetClaimsInContext(ctx echo.Context) *TokenClaims {
	jwtToken := ctx.Get(CtxKey)
	if jwtToken != nil {
		return ctx.Get(CtxKey).(*TokenClaims)
	}
	return nil
}

// GenerateExchangeCode generate exchange code by using nanoid
// after generated exchange code will be stored in redis
// format exchange code should be: exchange_<code>
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

// IsExchangeCodeExist verify file exchange code is existed in redis
func (a *Auth) IsExchangeCodeExist(code string) bool {
	return a.redis.Exists(fmt.Sprint("exchange_", code)).Val() == 1
}

// RevokeExchangeCode revoke exchange code stored in redis
// by remove key in redis
// format key: exchange_<code>
func (a *Auth) RevokeExchangeCode(code string) error {
	if code == "" {
		return nil
	}
	if !a.IsExchangeCodeExist(code) {
		return nil
	}

	return a.redis.Del(fmt.Sprint("exchange_", code)).Err()
}

// RequiredLogin middleware function
// this function required token in header
// accept prefix jwt, bearer, token
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
				return err
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

// RequiredExchangeCode middleware function
// this function required code in query param
func (a *Auth) RequiredExchangeCode() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// get code form query param
			code := ctx.QueryParam("code")
			if code == "" {
				return cerrors.UnAuthorized("missing exchange code")
			}
			// verify exist in redis
			if !a.IsExchangeCodeExist(code) {
				return cerrors.UnAuthorized("exchange code invalid")
			}

			return next(ctx)
		}
	}
}
