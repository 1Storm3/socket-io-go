package socketio

import (
	"sync"
)

type Room struct {
	Name           string
	To             func(room string) *Room
	sockets        connections
	connectSockets []*Socket
}

func newRoom(name string) *Room {
	return &Room{
		Name: name,
		sockets: connections{
			conn: make(map[string]*Socket),
		},
	}
}

func (room *Room) Emit(event string, agrs ...interface{}) error {
	if len(room.connectSockets) > 0 {
		for _, socket := range room.connectSockets {
			socket.Emit(event, agrs...)
		}
	} else {
		for _, socket := range room.sockets.all() {
			socket.Emit(event, agrs...)
		}
	}
	return nil
}

func (room *Room) Sockets() []*Socket {
	return room.sockets.all()
}

type roomNames struct {
	sync.RWMutex
	list []string
}

func (l *roomNames) set(name string) {
	l.Lock()
	l.list = append(l.list, name)
	l.Unlock()
}

func (l *roomNames) delete(name string) int {
	l.Lock()
	defer l.Unlock()
	for index, n := range l.list {
		if n == name {
			l.list = append(l.list[:index], l.list[index+1:]...)
			return index
		}
	}
	return -1
}

func (l *roomNames) all() []string {
	l.RLock()
	defer l.RUnlock()
	return append([]string{}, l.list...)
}

// func (l *roomNames) get(event string) []eventCallback {
// 	l.RLock()
// 	defer l.RUnlock()
// 	if _, ok := l.list[event]; !ok {
// 		return make([]eventCallback, 0)
// 	}
// 	ret := make([]eventCallback, 0)
// 	ret = append(ret, l.list[event]...)
// 	return ret
// }

type rooms struct {
	sync.RWMutex
	list map[string]*Room
}

func (n *rooms) create(name string) *Room {
	n.Lock()
	ret, ok := n.list[name]
	if !ok {
		ret = newRoom(name)
		n.list[name] = ret
	}
	n.Unlock()
	return ret
}

func (n *rooms) next(name string, preRoom ...*Room) *Room {
	n.Lock()
	ret, ok := n.list[name]
	if !ok {
		ret = newRoom(name)
		n.list[name] = ret
	}
	if len(preRoom) == 0 {
		newRoom := newRoom(name)
		newRoom.connectSockets = append(newRoom.connectSockets, ret.sockets.all()...)
		ret.To = func(room string) *Room {
			nextRoom := n.next(room, newRoom)
			newRoom.Name += "_" + nextRoom.Name
			newRoom.connectSockets = append(newRoom.connectSockets, nextRoom.sockets.all()...)
			return newRoom
		}
	} else {
		ret.To = func(room string) *Room {
			nextRoom := n.next(room, preRoom[0])
			preRoom[0].Name += "_" + nextRoom.Name
			preRoom[0].connectSockets = append(preRoom[0].connectSockets, nextRoom.sockets.all()...)
			return preRoom[0]
		}
	}
	n.Unlock()
	return ret
}

// func (n *rooms) get(name string) *Room {
// 	n.RLock()
// 	defer n.RUnlock()
// 	ret, ok := n.list[name]
// 	if !ok {
// 		return nil
// 	}
// 	return ret
// }
