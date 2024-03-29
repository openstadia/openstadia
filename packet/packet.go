package packet

import (
	"bytes"
	"encoding/json"
	"errors"
)

var Separator = []byte{'|'}

type Packet[T any] struct {
	Header  Header
	Payload T
}

func (p *Packet[T]) Decode(data []byte) error {
	parts := bytes.Split(data, Separator)
	if len(parts) != 2 {
		return errors.New("wrong packet format")
	}
	header, payload := parts[0], parts[1]

	err := json.Unmarshal(header, &p.Header)
	if err != nil {
		return err
	}

	err = json.Unmarshal(payload, &p.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (p *Packet[T]) Encode() ([]byte, error) {
	header, err := json.Marshal(p.Header)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(p.Payload)
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{header, payload}, Separator), nil
}

func MakeAck[T any, R any](pack *Packet[T], payload R) *Packet[R] {
	ack := Packet[R]{
		Header: Header{
			Type: TypeAck,
			Id:   pack.Header.Id,
		},
		Payload: payload,
	}

	return &ack
}
