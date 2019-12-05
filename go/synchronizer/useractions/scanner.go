package useractions

import relay "github.com/freeverseio/crypto-soccer/go/relay/storage"

type Scanner struct {
	db *relay.Storage
}

func NewScanner() (*Scanner, error) {
	scanner := Scanner{}
	return &scanner, nil
}
