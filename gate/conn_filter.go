// select the connector and sort

package main

import (
	. "github.com/hzhzh007/RoomChat/common"
	rpc_pool "github.com/hzhzh007/RoomChat/common/rpc"
	"log"
	"time"
)

var (
	connectorInfo *map[string][]RoomInfo
)

//get server list by room id
//TODO  implement it(return by room load ...)
func roomServerList(roomid string) ([]string, error) {
	conns := []string{"127.0.0.1:8888", "127.0.0.1:8889"}
	return conns, nil
}

//TODO  lock the map ? or pay attention to the gc
func updateConn() {
	temp := new(map[string][]RoomInfo)
	conf := GetConfig()
	for _, conn_addr := range conf.Connectors {
		client, err := rpc_pool.GetRpcClient(conn_addr.RpcAddr)
		if err != nil {
			log.Println("get rpc client error", err)
			continue
		}
		var reply []RoomInfo
		err = client.Call("RoomCenter.RoomInfo", "", &reply)
		if err != nil {
			log.Println("rpc broadcast error:", err)
			continue
		}
		temp[conn_addr.ServeAddr] = reply
		log.Println("rpc result", reply)
	}
	connectorInfo = temp
}

func updateConnectorInfo() {
	for {
		updateConn()
		time.Sleep(5 * time.Second)
	}
}
