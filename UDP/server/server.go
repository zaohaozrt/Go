package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{ //设置监听端口
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("listen failed,err:", err)
		return
	}
	defer listen.Close()
	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:]) //监听消息接受数据
		if err != nil {
			fmt.Println("read udp failed,err:", err)
			continue
		}
		fmt.Printf("来自客户端:data:%v addr:%v count:%v\n", string(data[:n]), addr, n)
		_, err = listen.WriteToUDP(data[:n], addr)
		if err != nil {
			fmt.Println("write to udp failed,err:", err)
			continue
		}
	}

}
