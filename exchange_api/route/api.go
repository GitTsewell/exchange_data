package route

import (
	"exchange_api/app"
	"exchange_api/tool"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.Use(Cors()) // 跨域

	r.POST("/login",app.LoginPost)

	// auth
	auth := r.Group("/")
	auth.Use(tool.JWTAuth())

	auth.GET("/depth",app.DepthIndex)
	auth.GET("/depth/edit",app.DepthEdit)
	auth.PUT("/depth/:platform",app.DepthUpdate)
	auth.GET("/depth/check/:platform",app.DepthCheck)
	auth.GET("/depth/commit",app.DepthCommit)

	return r
}
