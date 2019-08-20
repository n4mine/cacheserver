package web

import (
	"net/http"
	"strconv"

	"github.com/n4mine/cacheserver/cache"

	"github.com/gin-gonic/gin"
)

func httpGetInfoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]int64{"TotalCounterNumber": cache.CacheObj.Count()})
}

func httpGetInfoByNameHandler(c *gin.Context) {
	counterName := c.DefaultQuery("name", "")
	if len(counterName) == 0 {
		c.JSON(http.StatusBadRequest, "query string: name(string)")
		return
	}
	c.JSON(http.StatusOK, cache.CacheObj.GetInfoByKey(counterName))
}

func httpGetDataByNameHandler(c *gin.Context) {
	counterName := c.DefaultQuery("name", "")
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	if len(counterName) == 0 || len(from) == 0 || len(to) == 0 {
		c.JSON(http.StatusBadRequest, "query string: name(string), from(timestamp), to(timestamp)")
		return
	}

	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "from 不合法")
		return
	}
	toInt, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "to 不合法")
		return
	}

	res := make(map[uint32]float64)
	data, err := cache.CacheObj.Get(counterName, fromInt, toInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"msg": "stay away form me, go kiss your graph", "err": err.Error()})
		return
	}
	for _, iter := range data {
		for iter.Next() {
			t, v := iter.Values()
			res[t] = v
		}
	}

	c.JSON(http.StatusOK, res)
}
