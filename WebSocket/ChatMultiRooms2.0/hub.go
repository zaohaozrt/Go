// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	//room  path
	roomId string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(roomId string) *Hub {
	return &Hub{
		roomId:     roomId,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	defer func() {
		close(h.register)
		close(h.unregister)
		close(h.broadcast)
	}()
	for {
		select {
		case client := <-h.register: //加锁防止注销过程中注册进来一个client导致连接失败
			h.clients[client] = true
			mutex.Unlock()

		case client := <-h.unregister:

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			if len(h.clients) == 0 {
				delete(house, h.roomId) //只是删除字典kv值,hub依然存在

				mutex.Unlock()

				return //结束hub
			}

			mutex.Unlock()

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default: //发送不出去消息就删除客户端
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
