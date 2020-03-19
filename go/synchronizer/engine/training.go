package engine

import (
	"encoding/json"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Training struct {
	storage.Training
}

func NewTraining() *Training {
	training := Training{}
	training.Training = *(storage.NewTraining())
	return &training
}

func (b Training) Marshal() string {
	var result string
	s, _ := json.Marshal(b)
	result += string(s)

	return result
}
