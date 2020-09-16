package gql

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) AcceptOffer(args struct{ Input input.AcceptOfferInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] create auction %+v", args.Input)

	id, err := args.Input.AuctionID()
	if err != nil {
		return graphql.ID(""), err
	}

	isOwner, err := args.Input.IsSignerOwnerOfPlayer(b.contracts)
	if err != nil {
		return id, err
	}
	if !isOwner {
		return id, fmt.Errorf("signer is not the owner of playerId %v", args.Input.PlayerId)
	}

	playerIdString, ok := new(big.Int).SetString(args.Input.PlayerId, 10)
	if !ok {
		return id, fmt.Errorf("error converting playerId to bignum")
	}

	isValidForBlockchain, err := args.Input.IsValidForBlockchainFreeze(b.contracts)
	if err != nil {
		return id, err
	}
	if !isValidForBlockchain {
		return id, fmt.Errorf("blockchain failed trying to freeze the asset")
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}

	offerID, err := args.Input.AuctionID()
	if err != nil {
		return id, err
	}

	offer, err := tx.Offer(string(offerID))
	if err != nil {
		return id, err
	}

	currentTeamId, err := b.contracts.Market.GetCurrentTeamIdFromPlayerId(&bind.CallOpts{}, playerIdString)
	if err != nil {
		return id, errors.New("internal error: no currentTeamIdFromPlayerId")
	}
	if currentTeamId.String() == offer.BuyerTeamID {
		return id, errors.New("the buyerTeam already owns the player it is making an offer for")
	}

	if offer.State != storage.OfferStarted {
		return id, errors.New("Auctions can only be created for offers in Started state")
	}

	// TODO: Consider the need of this next check when DB does not allow it anyway
	highestOffer, err := getHigestOffer(tx, args.Input.PlayerId)
	if err != nil {
		return id, err
	}
	if highestOffer.AuctionID != string(offer.AuctionID) {
		return id, errors.New("You can only accept the highest offer")
	}

	if err := CreateAuctionFromOffer(tx, args.Input, offer); err != nil {
		tx.Rollback()
		return id, err
	}
	if err := b.CreateBidFromOffer(tx, args.Input, offer); err != nil {
		tx.Rollback()
		return id, err
	}
	return id, tx.Commit()
}

func CreateAuctionFromOffer(tx storage.Tx, in input.AcceptOfferInput, highestOffer *storage.Offer) error {
	auction := storage.NewAuction()
	id, err := in.AuctionID()
	if err != nil {
		return err
	}
	auction.ID = string(id)
	auction.Rnd = int64(in.Rnd)
	auction.PlayerID = in.PlayerId
	auction.CurrencyID = int(in.CurrencyId)
	auction.Price = int64(in.Price)
	if auction.ValidUntil, err = strconv.ParseInt(in.ValidUntil, 10, 64); err != nil {
		return fmt.Errorf("invalid validUntil %v", in.ValidUntil)
	}
	if auction.OfferValidUntil, err = strconv.ParseInt(in.OfferValidUntil, 10, 64); err != nil {
		return fmt.Errorf("invalid OfferValidUntil %v", in.OfferValidUntil)
	}
	auction.Signature = in.Signature
	auction.State = storage.AuctionStarted
	auction.StateExtra = ""
	auction.PaymentURL = ""
	signerAddress, err := in.SignerAddress()
	if err != nil {
		return err
	}
	auction.Seller = signerAddress.Hex()
	if err = tx.AuctionInsert(*auction); err != nil {
		return err
	}
	return nil
}

func getHigestOffer(tx storage.Tx, playerId string) (storage.Offer, error) {
	existingOffers, err := tx.OffersByPlayerId(playerId)
	if err != nil {
		return storage.Offer{}, errors.New("could not find existing offers")
	}
	if existingOffers == nil {
		return storage.Offer{}, errors.New("existingOffers is nil")
	}
	highestOffer, err := highestOfferFromExistingOffers(existingOffers)
	if err != nil {
		return storage.Offer{}, err
	}
	return *highestOffer, nil
}

func highestOfferFromExistingOffers(offers []storage.Offer) (*storage.Offer, error) {
	length := len(offers)
	if length == 0 {
		return nil, errors.New("There are no offers for this playerId")
	}
	if length == 1 {
		return &offers[0], nil
	}
	idx := -1
	price := int64(-1)
	for i, offer := range offers {
		if offer.State == storage.OfferStarted {
			if idx == -1 {
				idx = i
				price = offer.Price
			} else {
				if offer.Price > price {
					idx = i
					price = offer.Price
				}
			}
		}
	}
	if idx == -1 {
		return nil, errors.New("There are no acceptable offers")
	}
	return &offers[idx], nil
}

func (b *Resolver) CreateBidFromOffer(tx storage.Tx, acceptOfferIn input.AcceptOfferInput, highestOffer *storage.Offer) error {

	bidInput := input.CreateBidInput{}
	bidInput.Signature = highestOffer.Signature
	bidInput.AuctionId = graphql.ID(highestOffer.AuctionID)
	bidInput.ExtraPrice = int32(0)
	bidInput.Rnd = int32(0)
	bidInput.TeamId = highestOffer.BuyerTeamID

	_, err := b.CreateBid(struct{ Input input.CreateBidInput }{bidInput})

	if err != nil {
		return err
	}
	return nil
}
