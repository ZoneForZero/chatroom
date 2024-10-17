package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Agent struct {
	sr      map[int]*SpecialRoom
	df      *DefaultRoom
	usermap map[int]int
	ucmap   map[int]*websocket.Conn
}

func (a *Agent) Enter(rid int, uid int) {
	if orid, isExist := a.usermap[uid]; isExist {
		a.sr[orid].Quit(uid)
		delete(a.usermap, uid)
	}
	if rid <= 0 {
		a.df.Enter(uid)
	} else {
		if tempsr, isExist := a.sr[rid]; isExist {
			tempsr.Enter(uid)
		} else {
			a.sr[rid] = &SpecialRoom{
				roomId: rid,
			}
			a.sr[rid].Enter(uid)
		}
		a.usermap[uid] = rid
	}
}

func (a *Agent) Quit(uid int) {
	if rid, b := a.usermap[uid]; b {
		a.sr[rid].Quit(uid)
	} else {
		a.df.Quit(uid)
	}
}
func newagent() *Agent {
	a := new(Agent)
	a.sr = make(map[int]*SpecialRoom)
	a.df = newdfr()
	a.ucmap = make(map[int]*websocket.Conn)
	return a
}

var agent *Agent

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
	user := newUser(conn)
	user.login()
	agent.ucmap[user.userId] = conn
	for {
		fmt.Println('a')
		// 读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("2222Error reading message:", err)
			break
		}
		// todo 发给他的房间
		// 打印接收到的消息
		fmt.Printf("Received: %s\n", string(message))
	}
}

func main() {
	agent = newagent()
	go func() {
		// 连接到Websocket服务器
		time.Sleep(10 * time.Second)
		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
		log.Printf("Connecting to %s", u.String())
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("Dial:", err)
		}
		for {
			time.Sleep(5 * time.Second)
			msg := fmt.Sprintf("%d:%s", 0, "每隔5秒我就吵你们")
			c := conn
			//for _, c := range agent.ucmap {
			for {
				// 回显消息
				//fmt.Println("发到", c)
				time.Sleep(5 * time.Second)

				err := c.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Println("Error wri88ting message:", err)
					break
				}
			}
		}
		defer conn.Close()
	}()
	//// s
	//go func() {
	//	time.Sleep(15 * time.Second)
	//
	//	for {
	//		msg := fmt.Sprintf("%d:%s", 0, "每隔5秒我就吵你们")
	//		for _, c := range agent.ucmap {
	//			// 回显消息
	//			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
	//			if err != nil {
	//				log.Println("Error wri88ting message:", err)
	//				break
	//			}
	//		}
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
	server()
	//c

}
