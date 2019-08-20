package web

import (
	"github.com/n4mine/cacheserver/config"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Start(c config.WebConfig) {
	if !c.Enable {
		return
	}
	router := gin.New()
	router.Use(gin.Recovery())
	pprof.Register(router)
	register(router)

	router.Run(c.Port)
}
func register(router *gin.Engine) {
	router.GET("/self/ping", httpPingHandler)
	router.GET("/self/version", httpVersionHandler)
	router.GET("/getcount", httpGetInfoHandler)
	router.GET("/getinfobyname", httpGetInfoByNameHandler)
	router.GET("/getdatabyname", httpGetDataByNameHandler)
}
