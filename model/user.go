package model

import (
	"time"
)

type User struct {
	UserId			int64     `json:"user_id"`
	UserName    string    `json:"user_name"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Token       string    `json:"token"`
	LastLogin   time.Time `json:"last_login"`
}
