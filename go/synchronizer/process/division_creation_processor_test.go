package process_test

// func TestDivisionCreationProcess(t *testing.T) {
// 	db, err := storage.NewSqlite3("../../../universe.db/00_schema.sql")
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
// 	process, err := process.NewDivisionCreationProcessor(db, bc.Leagues)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	event := leagues.LeaguesDivisionCreation{
// 		Timezone:             uint8(1),
// 		CountryIdxInTZ:       big.NewInt(0),
// 		DivisionIdxInCountry: big.NewInt(0),
// 	}
// 	err = process.Process(event)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
