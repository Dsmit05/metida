package consts

import "time"

const (
	RoleUser  = "User"
	RoleAdmin = "Admin"
)

const (
	AccessTokenTTL  = time.Second * 60 * 5
	RefreshTokenTTL = time.Hour * 24
)
