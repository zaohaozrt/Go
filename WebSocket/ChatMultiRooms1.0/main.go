// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")
var house = make(map[string]*Hub) //管理所有hub

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// go run ./  运行程序   未加锁的版本
func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/{room}", serveHome)
	r.HandleFunc("/ws/{room}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) //提取这个客户端的动态参数room
		roomId := vars["room"]
		room, ok := house[roomId]
		var hub *Hub
		if ok { //房间之前已经创建
			hub = room
		} else {
			hub = newHub(roomId)
			house[roomId] = hub
			go hub.run()
		}
		serveWs(hub, w, r)

	})
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
