package models

import (
	"time"
)

// Blog блог приложения.
type Blog struct {
	ID          int32
	Name        string
	Description string
}

// Content объект, которым владеет пользователь.
type Content struct {
	ID          int32
	UserEmail   string
	Name        string
	Description string
}

// Session хранить информацию о сессиях пользователя.
type Session struct {
	ID           int32
	UserEmail    string
	RefreshToken string
	AccessToken  string // Можно выдавать сервисам токен с увеличенным сроком жизни, сейчас не используется.
	UserAgent    string
	IP           string
	ExpiresIn    int64
	CreatedAt    time.Time
}

// User хранить информацию о пользователе.
type User struct {
	ID        int32
	Name      string
	Password  string
	Email     string
	Role      string
	IsDeleted bool
}
