package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// 一个域名+端口为一个cookie储存对象，例如http://localhost:8080,储存网址下所有的cookie
var store = sessions.NewCookieStore([]byte("something-very-secret"))

// session内容保存在服务端，客户端只有sessionID，服务端通过ID找到对应的session
func main() {
	http.HandleFunc("/save", SaveSession)
	http.HandleFunc("/get", GetSession)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
}
func SaveSession(w http.ResponseWriter, r *http.Request) {
	//通过cookie的name找到对应的值，该值是sessionID
	//在通过sessionID找到服务端储存的session内容
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	session.Save(r, w) //r:将session保存的内容写到store w:将session-name：sessionID以cookie形式写给服务端
}
func GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name") //读取Cookie中session-name对应的值为每个用户的唯一标识
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	foo := session.Values["foo"]
	fmt.Println(foo)
}
