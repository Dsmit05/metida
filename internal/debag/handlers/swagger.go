package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ProfilingHandler .
type SwaggerHandler struct {
	r *gin.Engine
}

func NewSwaggerHandler() *gin.Engine {
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func (o *SwaggerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	o.r.ServeHTTP(w, req)
}
