package storage

import (
	"testing"

	"gotest.tools/assert"
)

func TestAuctionServiceInterface(t *testing.T, service AuctionService) {
	t.Run("TestAuctionByIDUnexistent", func(t *testing.T) {
		auction, err := service.Auction("4343")
		assert.NilError(t, err)
		assert.Assert(t, auction == nil)
	})

	t.Run("TestAuctionInsert", func(t *testing.T) {
		auction := NewAuction()
		auction.ID = "ciao"
		auction.Rnd = 4
		auction.PlayerID = "3"
		auction.CurrencyID = 3
		auction.Price = 3
		auction.ValidUntil = 3
		auction.Signature = "3"
		auction.State = AuctionStarted
		auction.StateExtra = "3"
		auction.PaymentURL = "3"
		auction.Seller = "3"
		assert.NilError(t, service.Insert(*auction))

		result, err := service.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *auction)
	})

	t.Run("TestPendingAuctions", func(t *testing.T) {
		auction := NewAuction()
		auction.ID = "ciao0"
		auction.State = AuctionStarted
		assert.NilError(t, service.Insert(*auction))
		result, err := service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)

		auction.ID = "ciao1"
		auction.State = AuctionAssetFrozen
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 3)

		auction.ID = "ciao2"
		auction.State = AuctionPaying
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 4)

		auction.ID = "ciao3"
		auction.State = AuctionWithdrableBySeller
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 5)

		auction.ID = "ciao4"
		auction.State = AuctionWithdrableByBuyer
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)

		auction.ID = "ciao5"
		auction.State = AuctionFailed
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)

		auction.ID = "ciao6"
		auction.State = AuctionEnded
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)

		auction.ID = "ciao7"
		auction.State = AuctionCancelled
		assert.NilError(t, service.Insert(*auction))
		result, err = service.PendingAuctions()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)
	})

	t.Run("TestAuctionUpdate", func(t *testing.T) {
		auction := NewAuction()
		auction.ID = "ciao20"
		auction.State = AuctionStarted
		auction.StateExtra = "priva"
		assert.NilError(t, service.Insert(*auction))
		result, err := service.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, AuctionStarted)
		assert.Equal(t, result.StateExtra, "priva")

		auction.State = AuctionCancelled
		auction.StateExtra = "privato"
		auction.PaymentURL = "http"
		assert.NilError(t, service.Update(*auction))

		result, err = service.Auction(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *auction)
	})

	t.Run("TestBid().Insert", func(t *testing.T) {
		auction := NewAuction()
		auction.ID = "0"
		assert.NilError(t, service.Insert(*auction))

		bid := NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, service.Bid().Insert(*bid))
	})

	t.Run("TestBidsByAuctionID", func(t *testing.T) {
		auction := NewAuction()
		auction.ID = "03"
		assert.NilError(t, service.Insert(*auction))

		bid := NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, service.Bid().Insert(*bid))
		bid.ExtraPrice = 10
		assert.NilError(t, service.Bid().Insert(*bid))

		bids, err := service.Bid().Bids(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(bids), 2)

		auction.ID = "1"
		assert.NilError(t, service.Insert(*auction))

		bid = NewBid()
		bid.AuctionID = auction.ID
		assert.NilError(t, service.Bid().Insert(*bid))

		bids, err = service.Bid().Bids(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(bids), 1)
	})

	t.Run("TestBidUpdate", func(t *testing.T) {
		auction := NewAuction()
		auction.ID = "04324"
		assert.NilError(t, service.Insert(*auction))

		bid := NewBid()
		bid.AuctionID = auction.ID
		bid.ExtraPrice = 10
		bid.State = BidAccepted
		assert.NilError(t, service.Bid().Insert(*bid))

		bid.State = BidPaid
		bid.StateExtra = "vciao"
		bid.PaymentID = "3"
		bid.PaymentURL = "http"
		bid.PaymentDeadline = 4
		assert.NilError(t, service.Bid().Update(*bid))

		bids, err := service.Bid().Bids(auction.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(bids), 1)
		assert.Equal(t, bids[0], *bid)
	})

	t.Run("TestBidFindBids", func(t *testing.T) {
		bids := []Bid{}
		bids = append(bids, *NewBid())
		bids = append(bids, *NewBid())
		bids = append(bids, *NewBid())
		bids = append(bids, *NewBid())
		result := FindBids(bids, BidPaid)
		assert.Equal(t, len(result), 0)
		result = FindBids(bids, BidAccepted)
		assert.Equal(t, len(result), 4)
		result[0].State = BidPaid
		result = FindBids(bids, BidPaid)
		assert.Equal(t, len(result), 1)
	})
}

func TestPlaystoreOrderServiceInterface(t *testing.T, service PlaystoreOrderService) {
	t.Run("insert", func(t *testing.T) {
		order := NewPlaystoreOrder()
		order.PurchaseToken = "ciao"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"

		assert.NilError(t, service.Insert(*order))
		result, err := service.Order(order.OrderId)
		assert.NilError(t, err)
		assert.Assert(t, result != nil)
		assert.Equal(t, *result, *order)
	})

	t.Run("pending orders", func(t *testing.T) {
		order := NewPlaystoreOrder()
		order.PurchaseToken = "ciao1"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"
		assert.NilError(t, service.Insert(*order))

		orders, err := service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "43d"
		order.State = PlaystoreOrderOpen
		assert.NilError(t, service.Insert(*order))

		orders, err = service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)

		order.PurchaseToken = "43d1"
		order.State = PlaystoreOrderAcknowledged
		assert.NilError(t, service.Insert(*order))

		orders, err = service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)

		order.PurchaseToken = "43d2"
		order.State = PlaystoreOrderComplete
		assert.NilError(t, service.Insert(*order))

		orders, err = service.PendingOrders()
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 2)
	})

	t.Run("update state", func(t *testing.T) {
		order := NewPlaystoreOrder()
		order.PurchaseToken = "ciao"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4"
		order.TeamId = "pippo"
		order.State = PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"
		assert.NilError(t, service.UpdateState(*order))

		order.State = PlaystoreOrderOpen
		order.StateExtra = "recdia"
		assert.NilError(t, service.UpdateState(*order))

		result, err := service.Order(order.OrderId)
		assert.NilError(t, err)
		assert.Equal(t, result.State, order.State)
		assert.Equal(t, result.StateExtra, order.StateExtra)

	})

	t.Run("pending order by playerId", func(t *testing.T) {
		order := NewPlaystoreOrder()
		order.PurchaseToken = "ciao12"
		order.PackageName = "dsd"
		order.ProductId = "444"
		order.OrderId = "fdrd"
		order.PlayerId = "4343534"
		order.TeamId = "pippo"
		order.State = PlaystoreOrderFailed
		order.StateExtra = "prova"
		order.Signature = "erere"

		orders, err := service.PendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		assert.NilError(t, service.Insert(*order))
		orders, err = service.PendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "ciao432423"
		order.State = PlaystoreOrderComplete
		assert.NilError(t, service.Insert(*order))
		orders, err = service.PendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 0)

		order.PurchaseToken = "ciao4324233"
		order.State = PlaystoreOrderAcknowledged
		assert.NilError(t, service.Insert(*order))
		orders, err = service.PendingOrdersByPlayerId(order.PlayerId)
		assert.NilError(t, err)
		assert.Equal(t, len(orders), 1)
	})
}

func TestOfferServiceInterface(t *testing.T, service OfferService) {
	t.Run("TestOfferByIDUnexistent", func(t *testing.T) {
		auction, err := service.Offer("4343")
		assert.NilError(t, err)
		assert.Assert(t, auction == nil)
	})

	t.Run("TestOfferInsert", func(t *testing.T) {
		offer := NewOffer()
		offer.ID = "ciao"
		offer.Rnd = 4
		offer.PlayerID = "3"
		offer.CurrencyID = 3
		offer.Price = 3
		offer.ValidUntil = 3
		offer.Signature = "3"
		offer.State = OfferStarted
		offer.StateExtra = "3"
		offer.Seller = "3"
		assert.NilError(t, service.Insert(*offer))

		result, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

	t.Run("TestPendingOffers", func(t *testing.T) {
		offer := NewOffer()
		offer.ID = "ciao0"
		offer.State = OfferStarted
		assert.NilError(t, service.Insert(*offer))
		result, err := service.PendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)

		offer.ID = "ciao5"
		offer.State = OfferFailed
		assert.NilError(t, service.Insert(*offer))
		result, err = service.PendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)

		offer.ID = "ciao6"
		offer.State = OfferEnded
		assert.NilError(t, service.Insert(*offer))
		result, err = service.PendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)

		offer.ID = "ciao7"
		offer.State = OfferCancelled
		assert.NilError(t, service.Insert(*offer))
		result, err = service.PendingOffers()
		assert.NilError(t, err)
		assert.Equal(t, len(result), 6)
	})

	t.Run("TestOfferUpdate", func(t *testing.T) {
		offer := NewOffer()
		offer.ID = "ciao20"
		offer.State = OfferStarted
		offer.StateExtra = "priva"
		offer.Seller = "yo"
		offer.AuctionID = "123456"
		assert.NilError(t, service.Insert(*offer))
		result, err := service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.State, OfferStarted)
		assert.Equal(t, result.StateExtra, "priva")
		assert.Equal(t, result.Seller, "yo")
		assert.Equal(t, result.AuctionID, "123456")

		offer.State = OfferCancelled
		offer.StateExtra = "privato"
		offer.Seller = "yo2"
		offer.AuctionID = "1234562"
		assert.NilError(t, service.Update(*offer))

		result, err = service.Offer(offer.ID)
		assert.NilError(t, err)
		assert.Equal(t, *result, *offer)
	})

}
