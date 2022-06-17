package consts

import "time"

const (
	RoleUser  = "User"
	RoleAdmin = "Admin"
)

const (
	AccessTokenTTL  = time.Second * 60 * 30
	RefreshTokenTTL = time.Hour * 48
)
