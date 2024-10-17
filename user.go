package main

import "github.com/gorilla/websocket"

var globaluserindex = 1

type User struct {
	userId int
	wsConn *websocket.Conn
}
type userManager struct {
	usermap map[int]*User
}

var UM userManager

func newUM() userManager {
	return userManager{
		usermap: make(map[int]*User),
	}
}
func (um *userManager) newUser(c *websocket.Conn) User {
	var uid = globaluserindex
	globaluserindex++
	rt := User{
		userId: uid,
		wsConn: c,
	}
	um.usermap[uid] = &rt
	return rt
}
func (um *userManager) getUser(uid int) *User {
	return um.usermap[uid]
}
func (u *User) login(rid int) {
	RM.Enter(u.userId, rid)
}
func (u *User) quit() {
	RM.QuitAnywhere(u.userId)
}
func (u *User) Send(msg []byte) {
	RM.SendByUser(u.userId, msg)
}
