package route

import (
	"exchange_api/app"
	"exchange_api/middleware"
	"exchange_api/tool"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.Use(Cors()) // 跨域
	r.Use(middleware.LoggerToFile()) // 日志

	r.POST("/login",app.LoginPost)

	// auth
	auth := r.Group("/")
	auth.Use(tool.JWTAuth())

	auth.GET("/depth",app.DepthIndex)
	auth.GET("/depth/edit",app.DepthEdit)
	auth.PUT("/depth/:platform",app.DepthUpdate)
	auth.GET("/depth/check/:platform",app.DepthCheck)
	auth.GET("/depth/commit",app.DepthCommit)

	auth.GET("/system",app.SystemIndex)
	auth.PUT("/system",app.SystemUpdate)
	auth.GET("/system/:key",app.SystemExec)

	auth.GET("/exchange",app.ExchangeEdit)
	auth.PUT("/exchange",app.ExchangeUpdate)

	return r
}
