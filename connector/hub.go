package main

import ()

type hub struct {
	roomid string
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

func newHub(roomid string) *hub {
	return &hub{
		roomid:      roomid,
		broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}
}

//TODO: clean hub
func (h *hub) run() {
	ZkAddRoomConn(g_zk, h.roomid)
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			if h.remove(&g_RoomCenter) {
				break
			}

		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}
func (h *hub) remove(rc *RoomCenter) bool {
	if len(h.connections) != 0 {
		return false
	}
	rc.lock.Lock()
	defer rc.lock.Unlock()
	if len(h.connections) != 0 {
		return false
	}
	ZkDeleteRoomConn(g_zk, h.roomid)
	delete(g_RoomCenter.rooms, h.roomid)
	return true
}
