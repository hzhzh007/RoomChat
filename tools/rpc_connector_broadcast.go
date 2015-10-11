package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type BroadcastInfo struct {
	Roomid string
	Msg    []byte
}

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8082")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var args = "hello rpc"
	var reply string
	b := &BroadcastInfo{Roomid: "0", Msg: []byte(args)}
	err = client.Call("RoomCenter.Broadcast", b, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("rpc result", reply)
}
