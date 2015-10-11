package common

import (
	"encoding/json"
	. "io"
	"net/http"
)

type Msg struct {
	From    string
	Roomid  string
	Userid  string
	Msg     []byte
	Msgtype string
}

type RoomInfo struct {
	Roomid string
	Conned int
}

type BroadcastInfo struct {
	Roomid string
	Msg    []byte
}

type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

const (
	//authorization error
	ERR_AUTH = 401
	// forbid
	ERR_FORBID = 403

	// internal err
	ERR_INTERNAL = 500

	//Ok
	OK_CODE = 0

	OK_STR = "OK"
)

func WriteJsonResponse(w http.ResponseWriter, status int, msg string, data interface{}) (int, error) {
	var r = Response{Status: status, Msg: msg, Data: data}
	b, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	return w.Write(b)
}

func (m Msg) Validate() (int, string) {
	if len(m.Msg) == 0 || len(m.From) == 0 {
		return ERR_AUTH, "bad param"
	}
	return 0, ""
}
