package useractions

import (
	"encoding/json"
)

type Stack []interface{}

func (b *Stack) Push(el interface{}) {
	*b = append(*b, el)
}

func (b Stack) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Stack) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}
