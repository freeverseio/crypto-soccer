package engine

import (
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
