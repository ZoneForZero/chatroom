package main

type Space interface {
	Enter(roomId int, userId int)
	Quit(userId int)
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

type SpecialRoom struct {
	Room
	roomId int
}

type DefaultRoom struct {
	Room
}

func newdfr() *DefaultRoom {
	d := new(DefaultRoom)
	d.userIds = make(map[int]bool)
	return d
}
