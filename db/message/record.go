package message

import "encoding/json"

type Record struct {
	Rev  int
	Key  string
	Data []byte
}

func (r *Record) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Record) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
