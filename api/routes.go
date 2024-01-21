package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *API) Run() {
	r := a.NewRouter()

	r.Run()
}

func (a *API) NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/signup", a.CreateAccountHandler)
		v1.POST("/login", a.VerifyCredentialHandler)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
