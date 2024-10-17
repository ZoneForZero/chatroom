package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

const dfid = 0

func isdfr(rid int) bool {
	return rid == dfid
}

var RM roomManager

type roomManager struct {
	sr           map[int]*SpecialRoom
	df           *DefaultRoom
	defaultIndex int
	urmap        map[int]int
}

func newRM() roomManager {
	rt := roomManager{
		sr: make(map[int]*SpecialRoom),
		df: &DefaultRoom{
			Room{userIds: make(map[int]bool)},
		},
		defaultIndex: 1,
		urmap:        make(map[int]int),
	}
	return rt
}
func (rm *roomManager) createSR(rid int) *SpecialRoom {
	RM.sr[rid] = &SpecialRoom{
		Room{
			userIds: make(map[int]bool),
		},
		rid,
	}
	return RM.sr[rid]
}
func (rm *roomManager) QuitAnywhere(uid int) {
	if oldrid, isexist := rm.urmap[uid]; isexist {
		if isdfr(oldrid) {
			rm.df.Quit(uid)
		} else if sproom, isexist1 := rm.sr[oldrid]; isexist1 {
			sproom.Quit(uid)
		}
	}
}
func (rm *roomManager) SendByUser(uid int, msg []byte) {
	if oldrid, isexist := rm.urmap[uid]; isexist {
		if isdfr(oldrid) {
			rm.df.Send(uid, msg)
		} else if sproom, isexist1 := rm.sr[oldrid]; isexist1 {
			sproom.Send(uid, msg)
		}
	}
}
func (rm *roomManager) Enter(uid, rid int) {
	rm.QuitAnywhere(uid)
	if isdfr(rid) {
		RM.df.Enter(uid)
	} else if sr, isexist := rm.sr[rid]; isexist {
		sr.Enter(uid)
	}
	rm.urmap[uid] = rid
}

type Space interface {
	Enter(userId int)
	Quit(userId int)
	Send(sendUserId int, msg []byte)
}
type Room struct {
	userIds map[int]bool
}

func (r *Room) Enter(uid int) {
	r.userIds[uid] = true
}
func (r *Room) Quit(uid int) {
	delete(r.userIds, uid)
}
func (r *Room) Send(sendUserId int, msg []byte) {
	for uid, _ := range r.userIds {
		if uid != sendUserId {
			err := UM.getUser(uid).wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error wri88ting message:", err)
				break
			}
		}
	}
}

type SpecialRoom struct {
	Room
	roomId int
}

type DefaultRoom struct {
	Room
}

func getSR(rid int) (*SpecialRoom, error) {
	if isdfr(rid) {
		return &SpecialRoom{}, fmt.Errorf("cuo")
	}
	if r, isexist := RM.sr[rid]; isexist {
		return r, nil
	} else {
		return RM.createSR(rid), nil
	}
}
