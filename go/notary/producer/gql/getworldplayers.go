package gql

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) GetWorldPlayers(args struct{ Input input.GetWorldPlayersInput }) ([]*WorldPlayer, error) {
	log.Debugf("GetWorldPlayers %v", args)

	result := []*WorldPlayer{}

	if b.ch == nil {
		return result, errors.New("internal error: no channel")
	}

	hash, err := args.Input.Hash()
	if err != nil {
		return result, err
	}
	sign, err := hex.DecodeString(args.Input.Signature)
	if err != nil {
		return result, err
	}

	isValid, err := input.VerifySignature(hash, sign)
	if err != nil {
		return result, err
	}
	if !isValid {
		return result, errors.New("Invalid signature")
	}

	sender, err := input.AddressFromSignature(hash, sign)
	if err != nil {
		return result, err
	}
	log.Infof("TODO sender is %v", sender)

	// TODO put the 30 in a smarter place
	for i := 0; i < 30; i++ {
		worldPlayer := NewWorldPlayer(
			graphql.ID(i),
			"dummy"+fmt.Sprintf("%d", i),
		)
		result = append(result, worldPlayer)
	}

	return result, nil
}
