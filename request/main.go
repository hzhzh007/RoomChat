//
//TODO 1. rpc
package main

import (
	. "github.com/hzhzh007/RoomChat/common"
	log "github.com/hzhzh007/RoomChat/common/log"
	rpc_pool "github.com/hzhzh007/RoomChat/common/rpc"
	"math/rand"
	"net/http"
	"net/rpc"
	"time"
)

//random select
func selectRouter() string {
	config := GetConfig()
	index := rand.Intn(len(config.RouterAddr))
	return config.RouterAddr[index]
}

func selectRouterClient() (client *rpc.Client, err error) {
	return rpc_pool.GetRpcClient(selectRouter())
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
	msg.Msg = []byte(r.Form.Get("msg"))
	msg.Msgtype = r.Form.Get("type")
	status, m := msg.Validate()
	if status != OK_CODE {
		WriteJsonResponse(w, status, m, nil)
	}

	client, err := selectRouterClient()
	if err != nil {
		log.Error("get rpc client error", err)
		WriteJsonResponse(w, ERR_INTERNAL, "internal err", nil)
	}
	log.Println("start rpc call")
	var reply string
	err = client.Call("Router.DealMsg", msg, &reply)
	log.Println("rpc call end")
	if err != nil {
		log.Println("deal msg error:", err)
		WriteJsonResponse(w, ERR_INTERNAL, "internal dial err", nil)
	}
	WriteJsonResponse(w, OK_CODE, reply, nil)
}

func main() {
	loadConfig(false)
	http.HandleFunc("/postMsg", postMsgHandler)
	if err := http.ListenAndServe(GetConfig().ServeAddr, nil); err != nil {
		log.Fatal("ListenAndServe err:", err)
	}
}
func init() {
	rand.Seed(time.Now().UnixNano())
}
