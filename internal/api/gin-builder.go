package api

import (
	"net/http"

	_ "github.com/Dsmit05/metida/docs"
	"github.com/Dsmit05/metida/internal/api/controllers"
	"github.com/Dsmit05/metida/internal/api/middlewares"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/gin-gonic/gin"
)

type GinBuilder struct {
	userAuth    *controllers.UserAuth
	userContent *controllers.UserContent
	siteBlog    *controllers.SiteBlog
	*middlewares.ProtectedMidleware
	r *gin.Engine
}

func NewGinBuilder(
	db repositoryI,
	managerToken cryptographyI,
	cfg configGinBuilderI,
) *GinBuilder {

	userAuth := controllers.NewUserAuth(db, managerToken)
	wallEditorialsHandler := controllers.NewWallEditorials(db)
	siteBlog := controllers.NewSiteBlog(db)
	protectedMidleware := middlewares.NewProtectedMidleware(managerToken)

	var r *gin.Engine

	if cfg.IfDebagOn() {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	return &GinBuilder{
		userAuth,
		wallEditorialsHandler,
		siteBlog,
		protectedMidleware,
		r,
	}
}

func (o *GinBuilder) AddV1(basePath string) *GinBuilder {
	logger.Info("GinBuilder.AddV1()", "add api v1")
	v1 := o.r.Group(basePath)

	control := v1.Group("/auth")
	{
		control.POST("/sign-up", o.userAuth.CreateUser)
		control.POST("/sign-in", o.userAuth.AuthenticationUser)
		control.POST("/refresh", o.userAuth.RefreshTokenUser)
	}

	lk := v1.Group("/lk")
	lk.Use(o.AuthMidleware)
	{
		lk.GET("/content/:id", o.userContent.ShowContent)
		lk.POST("/content", o.userContent.CreateContent)
		lk.POST("/blog", o.siteBlog.CreateBlog)
	}
	v1.GET("/blog/:id", o.siteBlog.ShowBlog)

	return o
}

func (o *GinBuilder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	o.r.ServeHTTP(w, req)
}
