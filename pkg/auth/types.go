package auth

import "github.com/dgrijalva/jwt-go"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type TokenClaims struct {
	User      User   `json:"user"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}