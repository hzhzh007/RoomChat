//
//TODO 1. rpc
package main

import (
	"errors"
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	rpc_pool "github.com/hzhzh007/RoomChat/common/rpc"
	"github.com/samuel/go-zookeeper/zk"
	"math/rand"
	"net/http"
	"net/rpc"
	"time"
)

var (
	g_zk                 *zk.Conn
	g_routers            *[]string
	ROUTER_RPC_INTERFACE = "Router.DealMsg"
)

//random select
func selectRouter() string {
	list := g_routers
	if list == nil || len(*list) == 0 {
		return ""
	}
	index := rand.Intn(len(*list))
	log.Debug("rand index :%d, value", index, (*list)[index])
	return (*list)[index]
}

func UpdateRouters() {
	list := GetRouters(g_zk)
	g_routers = &list
}

func timeUpdaterRouter() {
	UpdateRouters()
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for t := range ticker.C {
		_ = t
		UpdateRouters()
	}
}

func selectRouterClient() (client *rpc.Client, err error) {
	return rpc_pool.GetRpcClient(selectRouter())
}

//todo
func selectRouterClientAdnCall(rpcName string, arg interface{}, ret interface{}) error {
	routerAddr := selectRouter()
	if len(routerAddr) == 0 {
		return errors.New("no router")
	}
	client, err := rpc_pool.GetRpcClient(routerAddr)
	if client == nil {
		return errors.New("no router")
	}
	err = client.Call(rpcName, arg, ret)
	if err == rpc.ErrShutdown {
		rpc_pool.RemoveRpcClient(routerAddr)
	}
	return err
}

// post param
// uid:
// msg:
// type:
// roomid:
func postMsgHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("receive request")
	var msg Msg
	r.ParseForm()
	msg.From = r.Form.Get("uid")
	msg.Roomid = r.Form.Get("roomid")
	msg.Msg = r.Form.Get("msg")
	msg.Msgtype = r.Form.Get("type")
	status, m := msg.Validate()
	if status != OK_CODE {
		WriteJsonResponse(w, status, m, nil)
	}
	reply := ""
	err := selectRouterClientAdnCall(ROUTER_RPC_INTERFACE, msg, &reply)
	if err != nil {
		log.Println("deal msg error:", err)
		WriteJsonResponse(w, ERR_INTERNAL, "internal dial err", nil)
	}
	WriteJsonResponse(w, OK_CODE, reply, nil)
}

func main() {
	loadConfig(false)
	http.HandleFunc("/postMsg", postMsgHandler)
	g_zk, _ = InitZK()
	go timeUpdaterRouter()
	if err := http.ListenAndServe(GetConfig().ServeAddr, nil); err != nil {
		log.Fatal("ListenAndServe err:", err)
	}
}
func init() {
	rand.Seed(time.Now().UnixNano())
}
