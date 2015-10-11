//
package main

import (
	"errors"
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	rpc_pool "github.com/hzhzh007/RoomChat/common/rpc"
	"net/http"

	"net/rpc"
)

type Router int

var (
	g_router Router
)

//TODO implement the logic
func (r *Router) DealMsg(m Msg, reply *string) (err error) {
	log.Println("recv post msg", m)
	conns, err := getConnectors(m.Roomid)
	if err != nil {
		log.Println("get connectors addrs error", err)
		return err
	}
	b := &BroadcastInfo{Roomid: m.Roomid, Msg: []byte(m.Msg)}
	log.Println("start rpc call connector broadcast", m)
	for _, connAddr := range conns {
		client, err := rpc_pool.GetRpcClient(connAddr)
		if err != nil {
			log.Println("get rpc client error", err)
			continue
		}
		err = client.Call("RoomCenter.Broadcast", b, &reply)
		if err != nil {
			log.Println("rpc broadcast error:", err)
			*reply = "error"
		}
		log.Println("rpc result", *reply)
	}
	*reply = "OK"
	return
}

//TODO to implement the logic
func getConnectors(roomid string) ([]string, error) {
	config := GetConfig()
	if config == nil {
		log.Error("get config error")
		return []string{}, errors.New("get config error")
	}
	return config.ConnAddr, nil
}

func main() {
	loadConfig(false)
	rpc.Register(&g_router)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(GetConfig().ServeAddr, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}

}
