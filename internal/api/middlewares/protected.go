package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Dsmit05/metida/internal/api/response"

	"github.com/Dsmit05/metida/internal/cryptography"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserRole = errors.New("the user has a negative role")
)

type ProtectedMidleware struct {
	auth cryptography.ManagerToken
}

func NewProtectedMidleware(token cryptography.ManagerToken) *ProtectedMidleware {
	return &ProtectedMidleware{auth: token}
}

func (o *ProtectedMidleware) AuthMidleware(c *gin.Context) {
	email, role, err := o.parseAuthHeader(c)
	if err != nil {
		response.GinError(c, http.StatusExpectationFailed, response.CodeUnknownUser, "Please log in", err)
		return
	}

	c.Set("email", email)
	c.Set("role", role)
}

func (o *ProtectedMidleware) parseAuthHeader(c *gin.Context) (email, role string, err error) {
	header := c.GetHeader("Authorizations")
	if header == "" || len(header) > 250 {
		return "", "", fmt.Errorf("bad header")
	}

	return o.auth.ParseToken(header)
}

func CheckAccessRights(c *gin.Context, roles ...string) error {
	roleFromContext, ok := c.Get("role")
	if !ok {
		return ErrUserRole
	}

	for _, role := range roles {
		if roleFromContext == role {
			return nil
		}
	}

	return ErrUserRole
}
