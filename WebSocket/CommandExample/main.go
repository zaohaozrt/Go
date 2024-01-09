// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr    = flag.String("addr", "127.0.0.1:8080", "http service address")
	cmdPath string
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func pumpStdin(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil }) //坚持到下一次ping
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		//s.Bytes() 每次读取刚扫描的那一段
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done) //关闭后，如果已取出原有的值，再取就是0值；一个信号：所有程序停止工作

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	//发送关闭信号
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	//关闭
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done: //收到程序终止信号
			return
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

var upgrader = websocket.Upgrader{}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()
	//创建用于标准输出的Pipe(由新建的进程输出)
	outr, outw, err := os.Pipe()
	if err != nil {
		internalError(ws, "stdout:", err)
		return
	}
	defer outr.Close()
	defer outw.Close()
	//创建用于标准输入的Pipe(向新建的进程输入)
	inr, inw, err := os.Pipe()
	//  w<- ==>  <-r
	if err != nil {
		internalError(ws, "stdin:", err)
		return
	}
	defer inr.Close()
	defer inw.Close()

	proc, err := os.StartProcess(cmdPath, flag.Args(), &os.ProcAttr{
		//standard input, standard output, and standard error
		Files: []*os.File{inr, outw, outw},
	})
	if err != nil {
		internalError(ws, "start:", err)
		return
	}

	inr.Close()
	outw.Close() //先关闭，新建进程会自动打开并读入输出

	stdoutDone := make(chan struct{}) //进程结束的标志
	//处理进程的输出
	go pumpStdout(ws, outr, stdoutDone)
	//跟客户端保持持久的连接
	go ping(ws, stdoutDone)
	//处理进程的输入
	pumpStdin(ws, inw)

	// Some commands will exit when stdin is closed.
	inw.Close()

	//连接运行中出现错误，进程存在也没有意义了，去死吧(程序终止)
	// Other commands need a bonk on the head.send the process an interrupt!
	if err := proc.Signal(os.Interrupt); err != nil {
		log.Println("inter:", err)
	}

	select {
	case <-stdoutDone: //已经终止
	case <-time.After(time.Second): //还没有终止，再来一棒，给我强制退出
		// A bigger bonk on the head. force the process to exit!
		if err := proc.Signal(os.Kill); err != nil {
			log.Println("term:", err)
		}
		<-stdoutDone //终于退出了
	}
	//waits for the Process to exit
	if _, err := proc.Wait(); err != nil {
		log.Println("wait:", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("must specify at least one argument")
	}
	var err error
	cmdPath, err = exec.LookPath(flag.Args()[0]) //在环境变量或者当前目录中寻找PATH，返回一个字符串（可能是绝对路径或相对路径）。
	fmt.Println("cmd:", cmdPath)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
