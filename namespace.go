package socketio

import (
	"sync"
)

type Namespace struct {
	Name         string
	sockets      *connections
	rooms        *rooms
	onConnection connectionEvent
}

func newNamespace(name string) *Namespace {
	return &Namespace{
		Name: name,
		sockets: &connections{
			conn: make(map[string]*Socket),
		},
		rooms: &rooms{
			list: make(map[string]*Room),
		},
		onConnection: connectionEvent{
			list: make(map[string][]connectionEventCallback),
		},
	}
}

func (nps *Namespace) OnConnection(fn connectionEventCallback) {
	nps.onConnection.set("connection", fn)
}

func (nps *Namespace) Emit(event string, agrs ...interface{}) error {
	for _, socket := range nps.sockets.all() {
		socket.Emit(event, agrs...)
	}
	return nil
}

func (nps *Namespace) socketJoinRoom(room string, socket *Socket) {
	nps.rooms.create(room).sockets.set(socket)
	socket.rooms.set(room)
}

func (nps *Namespace) socketLeaveRoom(room string, socket *Socket) {
	if socket.rooms.delete(room) != -1 {
		nps.rooms.create(room).sockets.delete(socket.Id)
	}
}

func (nps *Namespace) socketLeaveAllRooms(socket *Socket) {
	for _, room := range socket.rooms.all() {
		nps.socketLeaveRoom(room, socket)
	}
}

func (nps *Namespace) To(room string) *Room {
	return nps.rooms.next(room)
}

func (nps *Namespace) Sockets() []*Socket {
	return nps.sockets.all()
}

type namespaces struct {
	sync.RWMutex
	list map[string]*Namespace
}

func (n *namespaces) create(name string) *Namespace {
	n.Lock()
	ret, ok := n.list[name]
	if !ok {
		ret = newNamespace(name)
		n.list[name] = ret
	}
	n.Unlock()
	return ret
}

// func (n *namespaces) all() []*Namespace {
// 	n.RLock()
// 	ret := make([]*Namespace, 0)
// 	for _, nps := range n.list {
// 		ret = append(ret, nps)
// 	}
// 	n.RUnlock()
// 	return ret
// }

func (n *namespaces) get(name string) *Namespace {
	n.RLock()
	defer n.RUnlock()
	ret, ok := n.list[name]
	if !ok {
		return nil
	}
	return ret
}
