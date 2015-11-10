//
package main

import (
	"encoding/json"
	"errors"
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	rpc_pool "github.com/hzhzh007/RoomChat/common/rpc"
	"github.com/samuel/go-zookeeper/zk"
	"net/http"

	"net/rpc"
)

type Router int

var (
	g_router Router
	g_zk     *zk.Conn
	g_room   = make(map[string][]string)
)

//TODO implement the logic
func (r *Router) DealMsg(m Msg, reply *string) (err error) {
	log.Println("recv post msg", m)
	conns, err := getConnectors(m.Roomid)
	if err != nil {
		log.Println("get connectors addrs error", err)
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(" json marshal error:%v", err)
		return err
	}
	b := &BroadcastInfo{Roomid: m.Roomid, Msg: data}
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
	l, ok := g_room[roomid]
	if ok {
		return l, nil
	}
	GetRoomConnectorFromZk(roomid)
	l, ok = g_room[roomid]
	if ok {
		return l, nil
	}
	return []string{}, errors.New("get room :" + roomid + "errror")
}
func GetRoomConnectorFromZk(roomid string) error {
	list := GetRoomConnector(g_zk, roomid)
	_, old := g_room[roomid]
	if len(list) > 0 {
		g_room[roomid] = list
	}
	if old == false {
		cancel := make(chan int, 0)
		go WatchRoomConnChange(g_zk, roomid, cancel)
	}
	return nil
}

func main() {
	loadConfig(false)
	rpc.Register(&g_router)
	rpc.HandleHTTP()
	g_zk, _ = InitZK()
	if err := http.ListenAndServe(GetConfig().ServeAddr, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}

}
