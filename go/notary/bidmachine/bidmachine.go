package bidmachine

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

// type BidMachine struct {
// 	market          marketpay.IMarketPay
// 	auction         *storage.Auction
// 	bid             *storage.Bid
// 	contracts       *contracts.Contracts
// 	freeverse       *ecdsa.PrivateKey
// 	signer          *signer.Signer
// 	postAuctionTime int64
// }

// func New(
// 	market marketpay.IMarketPay,
// 	auction *storage.Auction,
// 	bid *storage.Bid,
// 	contracts *contracts.Contracts,
// 	freeverse *ecdsa.PrivateKey,
// ) (*BidMachine, error) {
// 	if market == nil {
// 		return nil, errors.New("No market instance given")
// 	}
// 	if auction.State != storage.AUCTION_PAYING {
// 		return nil, errors.New("Auction is not in PAYING state")
// 	}
// 	if auction.ID != bid.Auction {
// 		return nil, errors.New("Bid of wrong auction")
// 	}
// 	postAuctionTime, err := contracts.ConstantsGetters.GetPOSTAUCTIONTIME(&bind.CallOpts{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &BidMachine{
// 		market,
// 		auction,
// 		bid,
// 		contracts,
// 		freeverse,
// 		signer.NewSigner(contracts, freeverse),
// 		postAuctionTime.Int64(),
// 	}, nil
// }

func FirstAlive(bids []storage.Bid) *storage.Bid {
	// first searching for PAYING bid
	for i := range bids {
		if bids[i].State == storage.BidPaying {
			return &bids[i]
		}
	}
	// then search for the highest ACCEPTED bid
	idx := -1
	extraPrice := int64(-1)
	for i, bid := range bids {
		if bid.State == storage.BidAccepted {
			if idx == -1 {
				idx = i
				extraPrice = bid.ExtraPrice
			} else {
				if bid.ExtraPrice > extraPrice {
					idx = i
					extraPrice = bid.ExtraPrice
				}
			}
		}
	}
	if idx == -1 {
		return nil
	}
	return &bids[idx]
}

// func (b *BidMachine) Process() error {
// 	switch b.bid.State {
// 	case storage.BIDPAYING:
// 		return b.processPaying()
// 	case storage.BIDACCEPTED:
// 		return b.processAccepted()
// 	case storage.BIDFAILED:
// 		return nil
// 	default:
// 		return errors.New("Unknown bid state")
// 	}
// }

// func (b *BidMachine) processPaying() error {
// 	now := time.Now().Unix()
// 	if now > b.bid.PaymentDeadline {
// 		b.bid.State = storage.BIDFAILED
// 		b.bid.StateExtra = "Expired"
// 		return nil
// 	}
// 	if b.bid.PaymentID == "" { // create order
// 		log.Infof("[bid] Auction %v, extra_price %v | create MarketPay order", b.bid.Auction, b.bid.ExtraPrice)
// 		price := fmt.Sprintf("%.2f", float64(b.auction.Price.Int64()+b.bid.ExtraPrice)/100.0)
// 		name := "Freeverse Player transaction"
// 		order, err := b.market.CreateOrder(name, price)
// 		if err != nil {
// 			return err
// 		}
// 		b.bid.PaymentID = order.TrusteeShortlink.Hash
// 		b.bid.PaymentURL = order.TrusteeShortlink.ShortURL
// 	} else { // check if order is paid
// 		log.Warningf("[bid] Auction %v, extra_price %v | waiting for order %v to be processed", b.bid.Auction, b.bid.ExtraPrice, b.bid.PaymentID)
// 		order, err := b.market.GetOrder(b.bid.PaymentID)
// 		if err != nil {
// 			return err
// 		}
// 		paid := b.market.IsPaid(*order)
// 		if paid {
// 			isOffer2StartAuction := false
// 			bidHiddenPrice, err := b.signer.BidHiddenPrice(big.NewInt(b.bid.ExtraPrice), big.NewInt(b.bid.Rnd))
// 			if err != nil {
// 				return err
// 			}
// 			auctionHiddenPrice, err := signer.HashPrivateMsg(b.auction.CurrencyID, b.auction.Price, b.auction.Rnd)
// 			if err != nil {
// 				return err
// 			}
// 			var sig [2][32]byte
// 			var sigV uint8
// 			_, err = b.signer.HashBidMessage(
// 				b.contracts.Market,
// 				b.auction.CurrencyID,
// 				b.auction.Price,
// 				b.auction.Rnd,
// 				b.auction.ValidUntil,
// 				b.auction.PlayerID,
// 				big.NewInt(b.bid.ExtraPrice),
// 				big.NewInt(b.bid.Rnd),
// 				b.bid.TeamID,
// 				isOffer2StartAuction,
// 			)
// 			if err != nil {
// 				return err
// 			}
// 			sig[0], sig[1], sigV, err = signer.RSV(b.bid.Signature)
// 			if err != nil {
// 				return err
// 			}
// 			tx, err := b.contracts.Market.CompletePlayerAuction(
// 				bind.NewKeyedTransactor(b.freeverse),
// 				auctionHiddenPrice,
// 				big.NewInt(b.auction.ValidUntil),
// 				b.auction.PlayerID,
// 				bidHiddenPrice,
// 				b.bid.TeamID,
// 				sig,
// 				sigV,
// 				isOffer2StartAuction,
// 			)
// 			if err != nil {
// 				b.bid.State = storage.BIDFAILED
// 				b.bid.StateExtra = err.Error()
// 				return err
// 			}
// 			receipt, err := helper.WaitReceipt(b.contracts.Client, tx, 60)
// 			if err != nil {
// 				b.bid.State = storage.BIDFAILED
// 				b.bid.StateExtra = "Timeout waiting for the receipt"
// 				return err
// 			}
// 			if receipt.Status == 0 {
// 				b.bid.State = storage.BIDFAILED
// 				b.bid.StateExtra = "Mined but receipt.Status == 0"
// 				return err
// 			}
// 			b.auction.PaymentURL = order.SettlorShortlink.ShortURL
// 			b.bid.State = storage.BIDPAID
// 		}
// 	}
// 	return nil
// }

// func (b *BidMachine) processAccepted() error {
// 	b.bid.State = storage.BIDPAYING
// 	b.bid.PaymentDeadline = b.auction.ValidUntil + b.postAuctionTime
// 	return nil
// }
