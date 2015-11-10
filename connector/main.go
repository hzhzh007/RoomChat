package main

import (
	"flag"
	log "github.com/hzhzh007/RoomChat/common/log"
	"github.com/samuel/go-zookeeper/zk"
	"net"
	"net/http"
	"net/rpc"
)

var (
	g_RoomCenter = RoomCenter{rooms: make(map[string]*hub)}
	g_zk         *zk.Conn
)

func rpcServer() {
	rpc.Register(&g_RoomCenter)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", GetConfig().RpcAddr)
	if e != nil {
		log.Fatal("Listen error:", e)
	}
	http.Serve(l, nil)
}

func main() {
	flag.Parse()
	loadConfig(false)
	go rpcServer()
	g_zk, _ = InitZK()
	http.Handle("/ws", g_RoomCenter)
	if err := http.ListenAndServe(GetConfig().ServeAddr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
