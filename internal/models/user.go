package models

import "time"

type User struct {
	ID            int       `json:"id" db:"id"`
	Auth0ID       string    `json:"auth0_id" db:"auth0_id"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	Picture       *string   `json:"picture,omitempty" db:"picture"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type CreateUserRequest struct {
	Auth0ID       string  `json:"auth0_id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	Picture       *string `json:"picture,omitempty"`
	EmailVerified bool    `json:"email_verified"`
}

type UpdateUserRequest struct {
	Name          *string    `json:"name,omitempty"`
	Picture       *string    `json:"picture,omitempty"`
	EmailVerified *bool      `json:"email_verified,omitempty"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

type APIKey struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	KeyHash     string     `json:"-" db:"key_hash"`
	KeyPrefix   string     `json:"key_prefix" db:"key_prefix"`
	Permissions []string   `json:"permissions" db:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateAPIKeyRequest struct {
	Name        string     `json:"name"`
	Permissions []string   `json:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type CreateAPIKeyResponse struct {
	APIKey *APIKey `json:"api_key"`
	Key    string  `json:"key"`
}

type UpdateAPIKeyRequest struct {
	Name        *string    `json:"name,omitempty"`
	Permissions *[]string  `json:"permissions,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}