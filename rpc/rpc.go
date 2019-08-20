package rpc

import (
	"log"
	"net"
	"net/rpc"

	"github.com/n4mine/cacheserver/config"

	"github.com/ugorji/go/codec"
)

type CacheServer struct{}

func Start(c config.RpcConfig) {
	if !c.Enable {
		return
	}

	cs := &CacheServer{}

	rpc.Register(cs)

	ln, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("listen tcp port: %v, error: %v\n", c.Port, err)
	}

	log.Printf("rpc server starting on port %v\n", c.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept conn error: %v\n", err)
			continue
		}

		var h codec.MsgpackHandle
		c := codec.MsgpackSpecRpc.ServerCodec(conn, &h)

		go rpc.ServeCodec(c)
	}
}
