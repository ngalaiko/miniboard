package sockets

import (
	"golang.org/x/net/websocket"
)

func (s *Sockets) itemsLoad(ws *websocket.Conn, userID string, req *request) {
	s.respond(ws, &response{
		ID:     req.ID,
		Target: "#items-list",
		Reset:  true,
	})
	s.itemsLoadmore(ws, userID, req)
}
