package models

import "time"

type User struct {
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"autoCreateTime:false"`
	UpdatedAt  time.Time `json:"autoUpdateTime:false"`
	AdminToken string    `json:"admin_token"`
}
