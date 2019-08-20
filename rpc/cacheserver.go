package rpc

import (
	"log"
	"sync"

	"github.com/n4mine/cacheserver/cache"
	"github.com/n4mine/cacheserver/chunks"
	"github.com/n4mine/cacheserver/models"
)

func (cs *CacheServer) Push(ps []models.Point, r *SimpleRpcResponse) error {
	var cnt, fail int64
	var err error
	for _, p := range ps {
		if err = cache.CacheObj.Push(p.Key, p.Timestamp, p.Value); err != nil {
			log.Printf("push obj error, obj: %v, error: %v\n", p, err)
			fail++
		}
		cnt++
	}

	r.Code = 0

	return nil
}

func (cs *CacheServer) Get(req []DataReq, resp *[]*DataResp) error {
	ch := make(chan *DataResp, 1e3)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(req)) // 因为query限制了查询曲线的条数, 所以这里的并发不会太大

		for _, r := range req {
			go func(singleReq DataReq) {
				defer wg.Done()

				key := singleReq.Key
				from := singleReq.From
				to := singleReq.To

				var singleResp DataResp

				// 请求参数有误
				if len(key) == 0 || from == 0 || to == 0 || from >= to {
					singleResp.Code = models.CodeUserErr
					singleResp.Msg = "请求参数有错误"
					singleResp.Key = singleReq.Key
					singleResp.From = singleReq.From
					singleResp.To = singleReq.To
					singleResp.Step = singleReq.Step
					singleResp.RRA = singleReq.RRA
					singleResp.Data = nil

					ch <- &singleResp
					return
				}

				data, err := cache.CacheObj.Get(key, from, to)
				// 获取数据错误
				if err != nil {
					var code int
					switch err {
					case models.ErrNonExistSeries:
						code = models.CodeNonExistSeries
					case models.ErrNonEnoughData:
						code = models.CodeNonEnoughErr
					case models.ErrInternalError:
						code = models.CodeInternalErr
					}

					singleResp.Code = code
					singleResp.Msg = err.Error()
					singleResp.Key = singleReq.Key
					singleResp.From = singleReq.From
					singleResp.To = singleReq.To
					singleResp.Step = singleReq.Step
					singleResp.RRA = singleReq.RRA
					singleResp.Data = nil

					ch <- &singleResp
					return
				}
				// 正常返回
				singleResp.Code = models.CodeSucc
				singleResp.Key = singleReq.Key
				singleResp.From = singleReq.From
				singleResp.To = singleReq.To
				singleResp.Step = singleReq.Step
				singleResp.RRA = singleReq.RRA
				singleResp.Data = data

				ch <- &singleResp
			}(r)
		}

		wg.Wait()
		close(ch)
	}()

	for r := range ch {
		*resp = append(*resp, r)
	}

	return nil
}

// SimpleRpcResponse 同transfer定义
type SimpleRpcResponse struct {
	// code == 0, normal
	// code >  0, exception
	Code int `msg:"code"`
}

// DataReq, 同query定义
type DataReq struct {
	Key  string `msg:"key"`
	From int64  `msg:"from"`
	To   int64  `msg:"to"`
	Step int    `msg:"step"`
	RRA  int    `msg:"rra"`
}

// DataResp, 同query定义
type DataResp struct {
	// code == 0, normal
	// code >  0, exception
	Code int           `msg:"code"`
	Msg  string        `msg:"msg"`
	Key  string        `msg:"key"`
	From int64         `msg:"from"`
	To   int64         `msg:"to"`
	Step int           `msg:"step"`
	RRA  int           `msg:"rra"`
	Data []chunks.Iter `msg:"data"`
}
