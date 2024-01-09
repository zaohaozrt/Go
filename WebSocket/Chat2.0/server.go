package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// 直接打开html文件，go run ./ 运行
func main() {
	router := mux.NewRouter()
	go h.run()
	router.HandleFunc("/ws", myws)
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}
