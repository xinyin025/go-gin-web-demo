package routes

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/controller"
	"go-web-demo/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthenticateJWT())
	{
		authorized.GET("/myInfo", controller.QueryUser)
	}

	return r
}
