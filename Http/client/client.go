package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, _ := http.Get("http://127.0.0.1:8000/go")
	defer resp.Body.Close()
	//200 OK
	fmt.Println(resp.Status)
	fmt.Println(resp.Header)

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		} else {
			fmt.Println("读取完毕")
			resp := string(buf[:n])
			fmt.Println(resp)
			break
		}

	}
}
