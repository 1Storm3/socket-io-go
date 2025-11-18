package socketio

import (
	"errors"
	"sync"
)

var (
	ErrorInvalidConnection = errors.New("invalid connection")
	ErrorUUIDDuplication   = errors.New("UUID already exists")
)

type connections struct {
	sync.RWMutex
	conn map[string]*Socket
}

func (l *connections) set(socket *Socket) {
	l.Lock()
	l.conn[socket.Id] = socket
	l.Unlock()
}

func (l *connections) get(key string) (*Socket, error) {
	l.RLock()
	ret, ok := l.conn[key]
	l.RUnlock()
	if !ok {
		return nil, ErrorInvalidConnection
	}
	return ret, nil
}

func (l *connections) all() []*Socket {
	l.RLock()
	ret := make([]*Socket, 0)
	for _, socket := range l.conn {
		ret = append(ret, socket)
	}
	l.RUnlock()
	return ret
}

func (l *connections) delete(key string) {
	l.Lock()
	delete(l.conn, key)
	l.Unlock()
}
