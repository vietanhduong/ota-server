package user

import "time"

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}
