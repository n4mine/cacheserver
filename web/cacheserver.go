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

func httpGetDataInfoHandler(c *gin.Context) {
	counterName := c.DefaultQuery("name", "")
	if len(counterName) == 0 {
		c.JSON(http.StatusBadRequest, "query string: name(string)")
		return
	}
	c.JSON(http.StatusOK, cache.CacheObj.GetInfoByKey(counterName))
}

func httpGetDataHandler(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"msg": err.Error()})
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

func httpPushHandler(c *gin.Context) {
	counterName := c.DefaultQuery("name", "")
	_ts := c.DefaultQuery("ts", "")
	_value := c.DefaultQuery("value", "")

	if len(counterName) == 0 || len(_ts) == 0 || len(_value) == 0 {
		c.JSON(http.StatusBadRequest, "query string: name(string), ts(timestamp), value(float64)")
		return
	}

	ts, err := strconv.ParseInt(_ts, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ts 不合法")
		return
	}
	value, err := strconv.ParseFloat(_value, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "value 不合法")
		return
	}

	err = cache.CacheObj.Push(counterName, ts, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
}
