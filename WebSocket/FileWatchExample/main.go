// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	filename  string
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		//文档更改时间在上次更新时间之前，不需要更新内容
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}

// 设置ReadWriteDeadline是为了防止读写时卡死
func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))

	//ReadMessage 方法将调用参数函数h,函数h将不断延长ReadDeadline，能够等到下一次服务端ping客户端
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {

		_, _, err := ws.ReadMessage() //有消息才会读，否则一直等待

		if err != nil { //只要一刷新就关闭了通道，ReadMessage就会出错，否则一直正常运行
			break

		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod) //发送ping帧
	fileTicker := time.NewTicker(filePeriod) //10s检查一次文件修改
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for { //循环向客户端传送消息
		select {
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readFileIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil { //文件已修改
				ws.SetWriteDeadline(time.Now().Add(writeWait)) //需要writeWait时间之内完成发送请求，单次有效
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		case <-pingTicker.C: //服务器发ping给浏览器，浏览器返回pong消息，表明自己仍然在线
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	//传输不是用的http协议，用的ws协议,握手阶段采用 http 协议
	ws, err := upgrader.Upgrade(w, r, nil) //与客户端建立ws协议连接，只建立一次
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n) //获去上次文件修改时间
	}
	go writer(ws, lastMod)
	reader(ws)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //其他非法路径的请求也会在这处理
		http.Error(w, "No found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //响应头

	//p:数据data  lastMod:最后修改时间 err:错误
	p, lastMod, err := readFileIfModified(time.Time{}) //对应时间0001/01/01 00:00

	if err != nil {
		p = []byte(err.Error()) //如果读取文件有误，将错误信息展示
		lastMod = time.Unix(0, 0)
	}
	var v = struct {
		Host    string
		Data    string
		LastMod string
	}{
		r.Host,
		string(p),
		strconv.FormatInt(lastMod.UnixNano(), 16), //将纳秒转化为16进制字符串，用作下次websocket请求参数
	}
	homeTempl.Execute(w, &v)
}

func main() {
	flag.Parse() //解析命令行
	if flag.NArg() != 1 {
		log.Fatal("filename not specified")
	}
	filename = flag.Args()[0]
	http.HandleFunc("/", serveHome) //只是第一次运行，后面就会请求/ws
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}

}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
