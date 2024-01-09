package main

import (
	"ginEssential/controller"
	"ginEssential/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.GET("/info", middleware.AuthMiddleware(), controller.Info)
	r.GET("/ws/:room/:name", func(ctx *gin.Context) {
		roomId := ctx.Param("room")
		name := ctx.Param("name")
		room, ok := house[roomId]
		var hub *Hub
		if ok { //房间之前已经创建
			hub = room
		} else {
			hub = newHub(roomId)
			house[roomId] = hub
			go hub.run()
		}
		serveWs(hub, ctx, name)
	})
	return r
}
