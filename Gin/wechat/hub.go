package main

import (
	"encoding/json"
)

type Hub struct {
	roomId string

	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}
type People struct {
	T         string   `json:"type"`
	RoomMates []string `json:"persons"`
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
		case client := <-h.register:
			h.clients[client] = true

			//计算现有人数
			var person []string

			for client := range h.clients {
				person = append(person, client.name)
			}

			people := People{T: "roommates", RoomMates: person}
			peopleBytes, _ := json.Marshal(people)
			//发送现有人数
			for client := range h.clients {
				select {
				case client.send <- peopleBytes:
				default: //发送不出去消息就删除客户端
					close(client.send)
					delete(h.clients, client)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			if len(h.clients) == 0 {
				delete(house, h.roomId) //只是删除字典kv值,hub依然存在
				return                  //结束hub
			}

			//计算现有人数
			var person []string

			for client := range h.clients {
				person = append(person, client.name)
			}

			people := People{T: "roommates", RoomMates: person}
			peopleBytes, _ := json.Marshal(people)
			//发送现有人数
			for client := range h.clients {
				select {
				case client.send <- peopleBytes:
				default: //发送不出去消息就删除客户端
					close(client.send)
					delete(h.clients, client)
				}
			}
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
