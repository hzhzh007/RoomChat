//receive data from other server
package main

import (
	"errors"
	. "github.com/hzhzh007/RoomChat/common"
	"net/http"
)

// handle websocket request
// TODO 1. map with manual rw lock
//      2. connection limit
//		3. validate connection info (ex. roomid, userid) if it needed
func (rc RoomCenter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roomid := r.Form.Get("roomid")
	if len(roomid) == 0 {
		roomid = "1"
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	hub, ok := rc.rooms[roomid]
	if !ok {
		//TODO rwlock
		hub = newHub(roomid)
		rc.rooms[roomid] = hub
		go hub.run()
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, h: hub}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()
	go c.writer()
	c.reader()
	//c.writer()
}

func (rc *RoomCenter) Broadcast(b *BroadcastInfo, reply *string) error {
	hub, ok := rc.rooms[b.Roomid]
	if !ok {
		*reply = "room not found"
		return errors.New(*reply)
	}
	hub.broadcast <- b.Msg
	*reply = "ok"
	return nil
}

//TODO roomid effects
func (rc *RoomCenter) RoomInfo(roomid string, reply *[]RoomInfo) error {
	rc.lock.RLock()
	defer rc.lock.RUnlock()

	for roomid, hub := range rc.rooms {
		*reply = append(*reply, RoomInfo{Roomid: roomid, Conned: len(hub.connections)})
	}
	return nil
}
