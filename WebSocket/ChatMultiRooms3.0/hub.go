// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"
)

// Hub： 每个房间的对于其中所有客户端的中央控制器
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Identity of room.
	roomId string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(roomId string) *Hub {
	return &Hub{
		roomId:     roomId,
		broadcast:  make(chan []byte),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	defer func() {
		close(h.unregister)
		close(h.broadcast)
	}()
	for {
		select {
		case client := <-h.unregister:
			roomMutex := roomMutexes[h.roomId]
			roomMutex.Lock() //加锁防止注销后注册的发生

			if _, ok := h.clients[client]; ok {
				time.Sleep(5 * time.Second)

				delete(h.clients, client)
				close(client.send)
				if len(h.clients) == 0 {
					roomMutex.Unlock() //把删除过程中新来的放出来，他卡在lock了，同时放到mutexForRoomMutexes.lock之前避免注册和注销过程中AB BA死锁发生

					//<--在这期间可能进来新用户-->

					//房间的锁也考虑要不要删了
					//如果已经有新来的了，就让他直接复用roomMutex和hub,不需要新建
					mutexForRoomMutexes.Lock() //新进来的用户一定已经获得lock锁
					if roomMutex.TryLock() {   //true：已经注册完毕/没人来 false:新来的还在注册过程中，占领锁
						if len(h.clients) == 0 { //0:没人来   !0:已经注册完毕，让他继续使用房间锁，不删了
							roomMutex.Unlock()
							delete(roomMutexes, h.roomId)
							house.Delete(h.roomId)
							return //确定没人来，直接摧毁
						} else {
							roomMutex.Unlock()
						}
					}
					mutexForRoomMutexes.Unlock()
				}
			} else {
				roomMutex.Unlock()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
