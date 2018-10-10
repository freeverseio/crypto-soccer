/*
const cryptoSoccer = artifacts.require("Testing");
var k = require('../jsCommons/constants.js');
var f = require('../jsCommons/functions.js');

contract('Leagues', function(accounts) {

  var instance;
  var nTotalPlayers=0;
  var sourceBalance;
  console.log('Funds in the source account:');
  console.log(web3.eth.getBalance(web3.eth.accounts[0]).toNumber()/web3.toWei(1, "ether"));

  it("creates a single contract and computes the gas cost of deploying GameEngine", async () => {
    instance = await cryptoSoccer.new();
    var receipt = await web3.eth.getTransactionReceipt(instance.transactionHash);
    console.log("League\n\tdeployment cost: ", receipt.gasUsed, "\n\tcontract address:", receipt.contractAddress)
    assert.isTrue(receipt.gasUsed > 2000000);
  });

  it("creates 6 teams", async () => {
    nTeams = 2;
    for (var t = 0; t < nTeams; t++) {
        await f.createTeam(instance, "Masnou", "Loop", k.MaxPlayersInTeam, f.createAlineacion(4,3,3));
    }

  });
  
});


*/