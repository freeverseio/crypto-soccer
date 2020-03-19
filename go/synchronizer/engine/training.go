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

func (b Training) ToString() string {
	var result string
	s, _ := json.MarshalIndent(b, "", "\t")
	result += string(s)
	return result
}
