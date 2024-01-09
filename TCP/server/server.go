package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	//监听端口
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen faild ,err:", err)
		return
	}
	for {
		conn, err := listen.Accept() //等待建立连接
		if err != nil {
			fmt.Println("accept faild,err:", err)
			continue
		}
		go process(conn) //启动一个goroutine处理连接
	}

}
func process(conn net.Conn) {
	defer conn.Close() //关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		fmt.Println("等待读取数据")
		n, err := reader.Read(buf[:]) //等待读取数据
		fmt.Println("读取数据完毕")
		if err != nil { //读取失败
			fmt.Println("read from client failed,err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("接收到client端发来的数据:", recvStr)
		fmt.Println("正在向client发送数据", recvStr)
		conn.Write([]byte(recvStr)) //向client发送数据

	}
}
