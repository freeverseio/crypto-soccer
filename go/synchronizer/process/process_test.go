package process_test

// TODO reactive
// func TestSyncTeams(t *testing.T) {
// 	tx, err := universedb.Begin()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer tx.Rollback()
// 	relaytx, err := relaydb.Begin()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer relaytx.Rollback()
// 	// storage, err := storage.NewPostgres("postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bc, err := testutils.NewBlockchainNodeDeployAndInit()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	namesdb, err := names.New("../../names/sql/names.db")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	p, err := process.NewEventProcessor(
// 		bc.Contracts,
// 		namesdb,
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	count, err := p.Process(tx, relaytx, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if count == 0 {
// 		t.Fatal("processed 0 blocks")
// 	}

// 	// the null timezone (0) is only used by the Academy Team
// 	if count, err := storage.TimezoneCount(tx); err != nil {
// 		t.Fatal(err)
// 	} else if count != 2 {
// 		t.Fatalf("Expected 2 time zones at time of creation,  actual %v", count)
// 	}

// 	// one country belongs to timezone = 0
// 	if count, err := storage.CountryCount(tx); err != nil {
// 		t.Fatal(err)
// 	} else if count != 2 {
// 		t.Fatalf("Expected 2 countries at time of creation,  actual %v", count)
// 	}

// 	// one team (the Academy) belongs to timezone = 0
// 	if count, err := storage.TeamCount(tx); err != nil {
// 		t.Fatal(err)
// 	} else if count != (128 + 1) {
// 		t.Fatalf("Expected 128 actual %v", count)
// 	}
// 	if count, err := storage.PlayerCount(tx); err != nil {
// 		t.Fatal(err)
// 		t.Fatalf("Expected 128*18=2304 actual %v", count)
// 	} else if count != 128*18 {
// 	}

// 	timezoneIdx := uint8(1)
// 	countryIdx := big.NewInt(0)

// 	tx0, err := bc.Contracts.Assets.TransferFirstBotToAddr(
// 		bind.NewKeyedTransactor(bc.Owner),
// 		timezoneIdx,
// 		countryIdx,
// 		crypto.PubkeyToAddress(bc.Owner.PublicKey),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = helper.WaitReceipt(bc.Client, tx0, 3)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var txs []*types.Transaction
// 	for i := 0; i < 24*4; i++ {
// 		var root [32]byte
// 		tx, err := bc.Contracts.Updates.SubmitActionsRoot(
// 			bind.NewKeyedTransactor(bc.Owner),
// 			root,
// 			"cid",
// 		)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		txs = append(txs, tx)
// 	}
// 	err = helper.WaitReceipts(bc.Client, txs, 3)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = p.Process(tx, relaytx, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	playerIdx := big.NewInt(30)
// 	playerID, err := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, playerIdx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	owner, err := bc.Contracts.Assets.GetOwnerPlayer(&bind.CallOpts{}, playerID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if owner.String() != storage.BotOwner {
// 		t.Fatalf("Owner is wrong %v", owner.String())
// 	}

// 	tx0, err = bc.Contracts.Assets.TransferFirstBotToAddr(
// 		bind.NewKeyedTransactor(bc.Owner),
// 		timezoneIdx,
// 		countryIdx,
// 		crypto.PubkeyToAddress(bc.Owner.PublicKey),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = helper.WaitReceipt(bc.Client, tx0, 3)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = p.Process(tx, relaytx, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	owner, err = bc.Contracts.Assets.GetOwnerPlayer(&bind.CallOpts{}, playerID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if owner != crypto.PubkeyToAddress(bc.Owner.PublicKey) {
// 		t.Fatalf("Owner is wrong %v", owner.String())
// 	}

// 	matchCount, err := storage.MatchesByTimezoneIdxCountryIdxLeagueIdx(tx, 1, 0, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(matchCount) != 56 {
// 		t.Fatalf("Wrong number of matches %v", len(matchCount))
// 	}

// 	matchEventsCount, err := storage.MatchEventCountByTimezoneCountryLeague(tx, 1, 0, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if matchEventsCount > 54*34 {
// 		t.Fatalf("Wrong numnber of match events > 54*34 %v", matchEventsCount)
// 	}
// }
