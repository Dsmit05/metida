package api

import (
	"net/http"
	"time"

	"github.com/Dsmit05/metida/internal/models"
)

type repositoryI interface {
	CreateUser(name string, password string, email string, role string) error
	ReadUser(email string) (*models.User, error)
	UpdateUser(email string, name string, password string, role string, isDeleted bool) error
	DeleteUser(email string) error
	CreateSession(email string, refreshToken string, userAgent string, ip string, expiresIn int64) error
	ReadSession(email string, userAgent string, ip string) (*models.Session, error)
	UpdateSession(email string, refreshToken string, newRefreshToken string, expiresIn int64) error
	UpdateSessionTokenOnly(refreshToken string, newRefreshToken string, expiresIn int64) error
	ReadEmailRoleWithRefreshToken(refreshToken string) (*models.UserEmailRole, error)
	DeleteSession(email string, ip string, userAgent string) error
	CreatContent(email string, name string, description string) error
	ReadContent(email string, id int32) (*models.Content, error)
	CreatBlog(name string, description string) error
	ReadBlog(id int32) (*models.Blog, error)
}

type cryptographyI interface {
	CreateToken(email string, role string, ttl time.Duration) (string, error)
	ParseToken(inputToken string) (email string, role string, err error)
	CreateRefreshToken() (string, error)
}

type configApiI interface {
	IfDebagOn() bool
	GetDebagAddr() string
	GetApiAddr() string
	GetApiReadTimeout() time.Duration
	GetApiWriteTimeout() time.Duration
}

type configGinBuilderI interface {
	IfDebagOn() bool
}

type metricI interface {
	MetricsMiddleware(next http.Handler) http.Handler
}
