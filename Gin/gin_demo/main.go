package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

// 定义中间
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("中间件开始执行了")
		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "中间件")
		// 执行函数
		c.Next() //放行 执行链式的下一个函数，执行完在next下面继续执行
		// c.Abort()//禁行  链式函数到此结束，不再调用后续的函数处理
		// 中间件执行完后续的一些事情
		status := c.Writer.Status()
		fmt.Println("中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("time:", t2)
	}
}
func main() {
	//创建一个服务
	ginServer := gin.Default()
	ginServer.Use(favicon.New("./滑稽图标.ico")) //通过全局中间件换网页图标

	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*") //全局加载
	// ginServer.LoadHTMLFiles("templates/index.html") //单个加载

	//加载资源文件才能被静态页面使用
	ginServer.Static("./static", "/static")

	//响应一个页面给前端
	ginServer.GET("/index", func(ctx *gin.Context) {
		// ctx.JSON(200, gin.H{"msg": "hello"})  //返回json数据

		ctx.HTML(200, "index.html", gin.H{
			"msg": "这是Go后台传递来的数据",
		})

	})

	//	获取前端传来的参数
	// /user/info?userid=xxx&username=xxx
	ginServer.GET("/user/info", func(ctx *gin.Context) {
		userid := ctx.Query("userid")
		username := ctx.Query("username")
		ctx.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	// /user/info/1/name
	ginServer.GET("/user/info/:userid/:username", func(ctx *gin.Context) {
		userid := ctx.Param("userid")
		username := ctx.Param("username")
		ctx.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	//前端给后端传递json
	ginServer.POST("/json", func(ctx *gin.Context) {
		//request.body
		data, _ := ctx.GetRawData()  //获取raw文本数据
		var m map[string]interface{} //变量m值nil
		_ = json.Unmarshal(data, &m)
		ctx.JSON(http.StatusOK, m)
	})

	//前端传来表单
	ginServer.POST("/user/add", func(ctx *gin.Context) {
		//该方法默认解析的是x-www-form-urlencoded或from-data格式的参数
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		ctx.JSON(200, gin.H{
			"msg":      "ok",
			"username": username,
			"password": password,
		})
	})

	//重定向 301
	ginServer.GET("/redirct", func(ctx *gin.Context) {
		ctx.Redirect(301, "https://www.baidu.com")
	})

	//404
	ginServer.NoRoute(func(ctx *gin.Context) {
		ctx.String(200, "404")
	})

	//路由组
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add", func(ctx *gin.Context) { ctx.String(200, "add") })
		userGroup.GET("/login", func(ctx *gin.Context) { ctx.String(200, "login") })
	}
	//中间件
	ginServer.GET("/middleware", MiddleWare(), func(ctx *gin.Context) {
		request, _ := ctx.Get("request")
		ctx.String(200, request.(string))
	})
	//服务器端口
	ginServer.Run(":8080")
}
