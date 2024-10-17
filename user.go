package main

import "github.com/gorilla/websocket"

var globaluserindex = 1

type User struct {
	userId int
	wsConn *websocket.Conn
}

func newUser(c *websocket.Conn) User {
	var uid = globaluserindex
	globaluserindex++
	return User{
		userId: uid,
		wsConn: c,
	}
}
func (u *User) login() {
	agent.Enter(0, u.userId)
}
func (u *User) quit() {
	agent.Quit(u.userId)
}
