package engineio

import (
	"time"
)

type ConnParameters struct {
	PingInterval time.Duration
	PingTimeout  time.Duration
	SID          string
	Upgrades     []string
	MaxPayload   int
}

type jsonParameters struct {
	SID          string   `json:"sid"`
	Upgrades     []string `json:"upgrades,omitempty"`
	PingInterval int      `json:"pingInterval,omitempty"`
	PingTimeout  int      `json:"pingTimeout,omitempty"`
	MaxPayload   int      `json:"maxPayload,omitempty"`
}

// func ReadConnParameters(r io.Reader) (ConnParameters, error) {
// 	var param jsonParameters
// 	if err := json.NewDecoder(r).Decode(&param); err != nil {
// 		return ConnParameters{}, err
// 	}

// 	return ConnParameters{
// 		SID:          param.SID,
// 		Upgrades:     param.Upgrades,
// 		PingInterval: time.Duration(param.PingInterval) * time.Millisecond,
// 		PingTimeout:  time.Duration(param.PingTimeout) * time.Millisecond,
// 		MaxPayload:   param.MaxPayload,
// 	}, nil
// }

func (p ConnParameters) ToJson() jsonParameters {
	arg := jsonParameters{
		SID:          p.SID,
		Upgrades:     p.Upgrades,
		PingInterval: int(p.PingInterval / time.Millisecond),
		PingTimeout:  int(p.PingTimeout / time.Millisecond),
		MaxPayload:   int(p.MaxPayload),
	}
	return arg
}
