package web

import (
	"net/http"

	"github.com/n4mine/cacheserver/models"

	"github.com/gin-gonic/gin"
)

func httpPingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func httpVersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"ver": models.GitVer, "buildTime": models.BuildTime})
}
