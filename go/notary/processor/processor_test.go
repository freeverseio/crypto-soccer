package processor_test

// func TestOutdatedAuction(t *testing.T) {
// 	bc, err := testutils.NewBlockchainNode()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = bc.DeployContracts(bc.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	now := time.Now().Unix()
// 	auction := storage.Auction{
// 		UUID:       uuid.New(),
// 		ValidUntil: now - 10,
// 		State:      storage.AUCTION_STARTED,
// 	}
// 	err = sto.CreateAuction(auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	processor, err := processor.NewProcessor(sto, bc.Contracts, bc.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = processor.Process()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	auctions, err := sto.GetAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if auctions[0].State != storage.AUCTION_NO_BIDS {
// 		t.Fatalf("Expected %v but %v", storage.AUCTION_NO_BIDS, auctions[0].State)
// 	}
// }

// func TestAuctionWithBid(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bc, err := testutils.NewBlockchainNode()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = bc.DeployContracts(bc.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	now := time.Now().Unix()
// 	auction := storage.Auction{
// 		UUID:       uuid.New(),
// 		PlayerID:   big.NewInt(65),
// 		Price:      big.NewInt(7),
// 		Rnd:        big.NewInt(73),
// 		ValidUntil: now + 100,
// 		Signature:  "0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c",
// 		State:      storage.AUCTION_STARTED,
// 	}
// 	err = sto.CreateAuction(auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bid := storage.Bid{
// 		Auction:         auction.UUID,
// 		State:           storage.BIDACCEPTED,
// 		PaymentDeadline: 0,
// 	}
// 	err = sto.CreateBid(bid)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	processor, err := processor.NewProcessor(sto, bc.Contracts, bc.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = processor.Process()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	auctions, err := sto.GetAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if auctions[0].State != storage.AUCTION_FAILED {
// 		t.Fatalf("Expected %v but %v", storage.AUCTION_FAILED, auctions[0].State)
// 	}
// 	bids, err := sto.GetBidsOfAuction(auctions[0].UUID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if bids[0].State != storage.BIDACCEPTED {
// 		t.Fatalf("Expects %v got %v", storage.BIDACCEPTED, bids[0].State)
// 	}
// }
