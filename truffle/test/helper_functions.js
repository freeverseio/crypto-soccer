const cryptoSoccer = artifacts.require("HelperFunctions");

contract('Helpers', function(accounts) {

  var instance;

  // deploy the contract before each test
  beforeEach(async () => {
    instance = await cryptoSoccer.deployed();
  });

  it("tests encode/decode with bits", async () => {
    input = [1714, 311, 42, 3];
    expected = 206864348850;
    bits = 12;
    len = 4;

    // encoding
    result = await instance.encode(len, input, bits);
    assert.equal(result, expected);

    // decoding
    output = await instance.decode(len, expected, bits);
    for (var i=0; i<len; i++)
      assert.equal(output[i], input[i])

    // get num at index
    for (var i=0;i<len;i++) {
      assert.equal(await instance.getNumAtIndex(expected,i,bits), input[i])
    }

    // set num at index
    newInput = [3, 3410, 790, 21]
    result = expected
    for (var i=0;i<len;i++) {
      result = await instance.setNumAtIndex(result, newInput[i], i, bits)
    }
    expected = 1456376979459
    assert.equal(result, expected)
    for (var i=0;i<len;i++) {
      assert.equal(await instance.getNumAtIndex(expected,i,bits), newInput[i])
    }

  });

  it("tests readNumbersFromUint", async () =>{
    var rnds = await instance.readNumbersFromUint.call(3,1234567890,10000);
    assert.isTrue(rnds[0]==7890);
    assert.isTrue(rnds[1]==3456);
    assert.isTrue(rnds[2]==12);
  });
 
  it("tests encodeIntoLongIntArray and readNumbersFromUint with longer states", async () =>{
    var original = [90,78,67,45,23,01,55];
    var longState = await instance.encodeIntoLongIntArray.call(original.length,original,100);
    var states = await instance.readNumbersFromUint.call(original.length,longState,100);
    for (var s=0; s<original.length;s++) {
      assert.isTrue(states[s] == original[s]);
    }
  });

  it("tests readNumbersFromHash, and shows that if factor is 1e4, you can safely read 19 nums. For a factor 1e6, these are 11 numbers.", async () =>{
    var numToRead = 25;
    var factor = 10000;
    var seed = 2;
    hash = await instance.computeKeccak256(seed)
    rnds = await instance.readNumbersFromUint.call(numToRead, hash, factor);
    // for (var n=0; n<numToRead; n++) {
      // console.log(n + " - " + rnds[n]);
    // } 
    assert.isTrue(rnds[18] != 0);
    assert.isTrue(rnds[20] == 0);
  });

  it("tests throwDice and throwDiceArray", async () =>{
    maxRnd = 1e10;
    var winsTeam2=0;
    var winsTeam3=0;
    for (var i=0; i<100; i++) {
      rnd = Math.floor(Math.random() * maxRnd);
      var winner = await instance.throwDice.call(1,9,rnd, maxRnd);
      var winner2 = await instance.throwDiceArray.call([1,9],rnd, maxRnd);
      assert.isTrue(winner.toNumber()==winner2.toNumber());
      var winner3 = await instance.throwDiceArray.call([1,4,5],rnd, maxRnd);
      winsTeam2 += winner.toNumber();
      if (winner3.toNumber()==2) { winsTeam3++};
    }
    console.log("For dice to be OK, this number should be close to 90: " + winsTeam2);
    console.log("For dice to be OK, this number should be close to 50: " + winsTeam3);
  });

  it("tests setNumAtPos", async () =>{
    var factor = 1000;
    var longState = 999999999999999; // this has 15 digits
    // interface: setNumAtPos.call(numToWrite, longState, pos, factor);
    longState = await instance.setNumAtPos.call(789, longState, 0, factor);
    longState = await instance.setNumAtPos.call(456, longState, 1, factor);
    longState = await instance.setNumAtPos.call(123, longState, 2, factor);
    longState = await instance.setNumAtPos.call(111, longState, 5, factor);
    assert.isTrue(longState == 111999999123456789);
  });

})


