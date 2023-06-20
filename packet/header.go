package packet

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Header struct {
	Type Type   `json:"type"`
	Id   *int   `json:"id"`
	Name string `json:"name"`
}

func (h *Header) Decode(data []byte) error {
	parts := bytes.Split(data, Separator)
	if len(parts) != 2 {
		return errors.New("wrong packet format")
	}
	header := parts[0]

	err := json.Unmarshal(header, h)
	if err != nil {
		return err
	}

	return nil
}
