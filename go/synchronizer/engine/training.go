package engine

import (
	"encoding/json"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Training struct {
	storage.Training
}

func NewTraining(sto storage.Training) *Training {
	training := Training{}
	training.Training = sto
	return &training
}

func (b Training) Marshal() string {
	var result string
	s, _ := json.Marshal(b)
	result += string(s)

	return result
}
