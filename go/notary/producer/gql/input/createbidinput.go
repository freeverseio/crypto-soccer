package input

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/graph-gophers/graphql-go"
)

type CreateBidInput struct {
	Signature  string
	Auction    graphql.ID
	ExtraPrice int32
	Rnd        int32
	TeamId     string
}

func (b CreateBidInput) ID() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%s%d%d%s", b.Signature, string(b.Auction), b.ExtraPrice, b.Rnd, b.TeamId)))
	return hex.EncodeToString(h.Sum(nil))
}

func (b CreateBidInput) Hash() (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	return common.Hash{}, nil
}
