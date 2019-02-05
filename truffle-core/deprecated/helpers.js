require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Testing = artifacts.require("TeamFactoryMock");

var k = require('../jsCommons/constants.js');

contract('Testing', function(accounts) {

  var instance;

  // deploy the contract before each test
  beforeEach(async () => {
    const Players = await Players.new().should.be.fulfilled;
    const Teams = await Teams.new().should.be.fulfilled;
    instance = await Testing.new(Teams.address, Players.address).should.be.fulfilled;
  });

  it("tests that we compute correctly the teams playing in a certain day", async () => {
    nTeams = 6;
    info = "";
    for (round=0; round<2*(nTeams-1); round++){
        if (round==nTeams-1) info += "\n";
        info += "Round " + round + ": ";
        for (game=0; game<3; game++) {
            result = await instance.test_teamsInGame(round, game, nTeams).should.be.fulfilled;
            info += result[0].toNumber() + "-" + result[1].toNumber() +", ";
        }
    }
    console.log(info);
    correctInfo = "Round 0: 0-1, 5-2, 4-3, Round 1: 2-0, 3-1, 4-5, Round 2: 0-3, 2-4, 1-5, Round 3: 4-0, 5-3, 1-2, Round 4: 0-5, 4-1, 3-2, \nRound 5: 1-0, 2-5, 3-4, Round 6: 0-2, 1-3, 5-4, Round 7: 3-0, 4-2, 5-1, Round 8: 0-4, 3-5, 2-1, Round 9: 5-0, 1-4, 2-3, ";
    assert.equal(info, correctInfo)
  });

  it("tests serialize/decode with bits", async () => {
    input = [1714, 311, 42, 3];
    expected = 206864348850;
    bits = 12;
    len = 4;

    // encoding
    result = await instance.test_serialize(len, input, bits);
    assert.equal(result, expected);

    // decoding
    output = await instance.test_decode(len, expected, bits);
    for (var i=0; i<len; i++)
      assert.equal(output[i], input[i])

    // get num at index
    for (var i=0;i<len;i++) {
      assert.equal(await instance.test_getNumAtIndex(expected,i,bits), input[i])
    }

    // set num at index
    newInput = [3, 3410, 790, 21]
    result = expected
    for (var i=0;i<len;i++) {
      result = await instance.test_setNumAtIndex(newInput[i], result, i, bits)
    }
    expected = 1456376979459
    assert.equal(result, expected)
    for (var i=0;i<len;i++) {
      assert.equal(await instance.test_getNumAtIndex(expected,i,bits), newInput[i])
    }
  });

  it("serializes 11 indices into uint256", async () =>{
    values = [12, 13, 14, 15, 16, 17, 18, 19 ,20, 21, 22]
    len = 11
    bits = 20
    expected = 35352669156133374304817968464537810297241735541930614246604812 // 20 bits
    //expected = 29249832381825824921796662605174067212 // 12 bits
    //expected = 27914334146814444649904307450892 // 10 bits
    result = 0
    for (var i=0;i<len;i++) {
      result = await instance.test_setNumAtIndex(values[i], result, i, bits)
    // console.log('setting value ' + values[i] + ' at idx: ' + i, ' result: '+ result)
    }
    assert.equal(result, expected)
    mask = ((1 << bits)-1)
    for (var i=0;i<len;i++) {
      value = await instance.test_getNumAtIndex(result, i, bits)
    // console.log('index at ' + i + ' is ' + value)
    }
  });
  it("tests encoding/decoding player states with bits", async () => {
    states = [
      [262,67,64,24,36,57,0],
      [335,34,64,38,49,65,1],
      [384,44,58,44,52,52,1],
      [311,60,54,33,44,59,1],
      [177,42,63,53,47,45,1],
      [323,39,38,72,47,51,2],
      [201,56,47,62,51,34,2],
      [238,51,56,39,43,59,2],
      [201,33,52,62,58,42,3],
      [177,56,34,30,63,65,3],
      [250,46,69,32,35,64,3]
    ]
    nPlayers = 11
    nStates = 7
    bits = 14
    var serializedStates = Array(nPlayers).fill(0)
    for (player=0; player<nPlayers; player++) {
      serializedStates[player] = await instance.test_serialize(nStates, states[player], bits);
    }
    // console.log("serialized states:\n", serializedStates)

    for (player=0; player<nPlayers; player++) {
      decodedState = await instance.test_decode(nStates, serializedStates[player], bits);
      for (i=0; i<nStates; i++) {
        assert.equal(decodedState[i], states[player][i])
      }
    }

  });

  it("tests throwDice and throwDiceArray", async () =>{
    maxRnd = 1e10;
    var winsTeam2=0;
    var winsTeam3=0;
    for (var i=0; i<100; i++) {
      rnd = Math.floor(Math.random() * maxRnd);
      var winner = await instance.test_throwDice.call(1,9,rnd, maxRnd);
      var winner2 = await instance.test_throwDiceArray.call([1,9],rnd, maxRnd);
      assert.isTrue(winner.toNumber()==winner2.toNumber());
      var winner3 = await instance.test_throwDiceArray.call([1,4,5],rnd, maxRnd);
      winsTeam2 += winner.toNumber();
      if (winner3.toNumber()==2) { winsTeam3++};
    }
    console.log("For dice to be OK, this number should be close to 90: " + winsTeam2);
    console.log("For dice to be OK, this number should be close to 50: " + winsTeam3);
  });

  it("tests the function getRndNumArrays used to get the randoms used in a game", async () =>{
    roundsPerGame = 18;
    seed = 0;
    info = "";
    var rndArray = await instance.test_getRndNumArrays.call(seed, roundsPerGame, k.BitsPerRndNum);
    for (var round=0; round < roundsPerGame; round++) { info += ", " + rndArray[round].toNumber(); }
    info += "\n";
    rndArray = await instance.test_getRndNumArrays.call(seed+1, roundsPerGame, k.BitsPerRndNum);
    for (var round=0; round < roundsPerGame; round++) { info += ", " + rndArray[round].toNumber(); }
    expectedInfo = ", 9571, 15311, 12640, 3044, 13878, 35, 5252, 12069, 2982, 16161, 902, 10850, 837, 9048, 13866, 5410, 11481, 9271\n, 3318, 8168, 15915, 10226, 13101, 2753, 12204, 4818, 3316, 10440, 6118, 16220, 11981, 11419, 8307, 7556, 11602, 1080";
    assert.equal(info,expectedInfo);
});

});



