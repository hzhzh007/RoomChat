// chat gate
//TODO here just a demo
package main

import (
	. "github.com/hzhzh007/RoomChat/common"
	"net/http"
)

func ConnectorRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomid := r.Form.Get("roomid")
	conns, _ := roomServerList(roomid)
	WriteJsonResponse(w, OK_CODE, "ok", conns)
}

func main() {
	loadConfig(false)
	go updateConnectorInfo()
	InitZK()
	http.HandleFunc("/connectors.json", ConnectorRequest)
	http.ListenAndServe(":8082", nil)
}
