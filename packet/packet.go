package packet

import "encoding/json"

type Type string

const (
	TypeEvent Type = "EVENT"
	TypeAck   Type = "ACK"
)

type Packet[T any] struct {
	Type Type `json:"type"`
	Data T    `json:"data"`
	Id   *int `json:"id"`
}

func (p *Packet[T]) Decode(data []byte) {
	err := json.Unmarshal(data, p)
	if err != nil {
		panic(err)
	}
}

func (p *Packet[T]) Encode() []byte {
	marshal, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return marshal
}
