package match

// func CreateDummyPlayer(
// 	t *testing.T,
// 	contracts *contracts.Contracts,
// 	skills [5]uint8,
// 	forwardness uint8,
// 	leftishness uint8,
// ) *Player {
// 	dayOfBirth := big.NewInt(0)
// 	gen := uint8(0)
// 	playerID := big.NewInt(2132321)
// 	potential := uint8(3)
// 	aggr := uint8(0)
// 	birthTraits := [4]uint8{potential, forwardness, leftishness, aggr}

// 	contracts.Engine.EncodePlayerSkills(
// 		&bind.CallOpts{},
// 		skills,
// 		dayOfBirth,
// 		gen,
// 		playerID,
// 		birthTraits,
// 		alignedEndOfLastHalf,
// 		redCardLastGame,
// 		gameNonStopping,
// 	)
// 	// var playerStateTemp = await engine.encodePlayerSkills(
// 	//         skills, dayOfBirth21, gen = 0, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
// 	//         alignedEndOfLastHalfTwoVec[0], redCardLastGame = false, gamesNonStopping = 0,
// 	//         injuryWeeksLeft = 0, subLastHalf, sumSkills
// 	//     ).should.be.fulfilled;
// }
