// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var house = make(map[string]*Hub) //管理所有hub

func serveHome(c *gin.Context) {
	log.Println(c.Request.URL)

	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed", "status": http.StatusMethodNotAllowed})
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{})
}

// go run ./  运行程序   未加锁的版本
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("*")
	r.GET("/:room", serveHome)
	r.GET("/ws/:room", func(ctx *gin.Context) {
		roomId := ctx.Param("room")
		room, ok := house[roomId]
		var hub *Hub
		if ok { //房间之前已经创建
			hub = room
		} else {
			hub = newHub(roomId)
			house[roomId] = hub
			go hub.run()
		}
		serveWs(hub, ctx)
	})
	r.Run(":8888")

}
