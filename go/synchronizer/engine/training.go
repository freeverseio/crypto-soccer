package engine

import "github.com/freeverseio/crypto-soccer/go/storage"

type Training struct {
	storage.Training
}

func NewTraining() *Training {
	training := Training{}
	training.Training = *storage.NewTraining()
	return &training
}
