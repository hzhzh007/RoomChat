package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var (
	g_RoomCenter = RoomCenter{rooms: make(map[string]*hub)}
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
	http.Handle("/ws", g_RoomCenter)
	if err := http.ListenAndServe(GetConfig().ServeAddr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
