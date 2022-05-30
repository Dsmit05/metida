package response

import (
	"fmt"
	"github.com/Dsmit05/metida/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code        int    `json:"code" example:"3"`
	Description string `json:"desc" example:"invalid input params"`
	Error       string `json:"error" example:"status bad request"`
}

func ginError(c *gin.Context, status, code int, description string, err error) {
	if err == nil {
		err = fmt.Errorf(http.StatusText(status))
	}

	if len(description) == 0 {
		description = "please try again later"
	}

	newErr := Error{
		Code:        code,
		Description: description,
		Error:       err.Error(),
	}

	c.JSON(status, newErr)
}

func GinError(c *gin.Context, status, code int, description string, err error) {
	logger.RequestsError(c, status, code, description, err)
	c.Abort()
	ginError(c, status, code, description, err)
}
