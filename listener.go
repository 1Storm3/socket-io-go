package socketio

import "sync"

type AckCallback func(data ...interface{})

type EventPayload struct {
	Name   string //event name
	SID    string //socket id
	Socket *Socket
	Error  error
	Data   []interface{}
	Ack    AckCallback
}

type eventCallback func(data *EventPayload)

type listeners struct {
	sync.RWMutex
	list map[string][]eventCallback
}

func (l *listeners) set(event string, callback eventCallback) {
	l.Lock()
	l.list[event] = append(l.list[event], callback)
	l.Unlock()
}

func (l *listeners) get(event string) []eventCallback {
	l.RLock()
	defer l.RUnlock()
	if _, ok := l.list[event]; !ok {
		return make([]eventCallback, 0)
	}
	ret := make([]eventCallback, 0)
	ret = append(ret, l.list[event]...)
	return ret
}
