package model

import (
	"time"
)

const (
	AccountHeaderUserInfo = "User-Info"
	DefaultExpirationUser = 3600
)

type AuthUser struct {
	UserID   int64
	Username string
	Email    string
	Fullname string
}

type Account struct {
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
