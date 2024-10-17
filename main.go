package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

func server() {
	http.HandleFunc("/ws", handleWebSocket)
	http.ListenAndServe(":8080", nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	user := UM.newUser(conn)
	user.login(dfid)
	for {
		// 接收用户发的信息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("2222Error reading message:", err)
			break
		}
		user.Send(message)
	}
}

var preonce sync.Once

func pre() {

	UM = newUM()
	RM = newRM()

}
func main() {
	pre()
	server()
}
