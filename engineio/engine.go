package engineio

import (
	"encoding/json"
	"io"
	"strconv"
)

type PacketType int

const (
	OPEN PacketType = iota
	CLOSE
	PING
	PONG
	MESSAGE
	UPGRADE
	NOOP
)

// Code	Message
// 0	"Transport unknown"
// 1	"Session ID unknown"
// 2	"Bad handshake method"
// 3	"Bad request"
// 4	"Forbidden"
// 5	"Unsupported protocol version"

func (id PacketType) String() string {
	return strconv.Itoa(int(id))
}

type writer struct {
	t PacketType
	i int64
	w io.Writer
}

func (w *writer) Write(p []byte) (int, error) {
	paserData := append([]byte(w.t.String()), p...)
	n, err := w.w.Write(paserData)
	w.i += int64(n)
	return n, err
}

func WriteTo(w io.Writer, t PacketType, arg ...interface{}) (int64, error) {
	writer := writer{
		t: t,
		w: w,
	}
	if len(arg) > 0 {
		err := json.NewEncoder(&writer).Encode(arg[0])
		return writer.i, err
	} else {
		_, err := writer.Write([]byte{})
		return writer.i, err
	}
}

func WriteByte(w io.Writer, t PacketType, p []byte) (int, error) {
	writer := writer{
		t: t,
		w: w,
	}
	return writer.Write(p)
}
