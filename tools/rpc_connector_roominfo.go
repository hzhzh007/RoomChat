package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type RoomInfo struct {
	Roomid string
	Conned int
}

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8082")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply []RoomInfo
	err = client.Call("RoomCenter.RoomInfo", "0", &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("rpc result", reply)
}
