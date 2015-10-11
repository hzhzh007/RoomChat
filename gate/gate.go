// chat gate
//TODO here just a demo
package main

import (
	. "github.com/hzhzh007/RoomChat/common"
	"net/http"
)

//get server list by room id
//TODO  implement it(return by room load ...)
func roomServerList(roomid string) ([]string, error) {
	conns := []string{"127.0.0.1:8888", "127.0.0.1:8889"}
	return conns, nil
}

func ConnectorRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomid := r.Form.Get("roomid")
	conns, _ := roomServerList(roomid)
	WriteJsonResponse(w, OK_CODE, "ok", conns)
}

func main() {
	loadConfig()
	go updateConnectorInfo()
	http.HandleFunc("/connectors.json", ConnectorRequest)
	http.ListenAndServe(":8082", nil)
}
