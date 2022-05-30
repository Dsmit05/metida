package response

import (
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/gin-gonic/gin"
)

type Success struct {
	Code        int         `json:"code" example:"0"`
	Description string      `json:"desc" example:"ok"`
	Data        interface{} `json:"data" swaggertype:"object,string" example:"key1:value,key2:value2"`
}

func GinSuccess(c *gin.Context, status, code int, data interface{}, description string) {
	logger.RequestsInfo(c, status, code, description, data)
	newResponse := Success{
		Code:        code,
		Description: description,
		Data:        data,
	}
	c.JSON(status, newResponse)
}
