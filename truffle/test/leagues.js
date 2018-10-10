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

  it("creates 4 teams and puts them into a league", async () => {
    nTeams = 4;
    teamsIdx = [];
    for (var t = 0; t < nTeams; t++) {
        var teamName = "team"+t;
        var playerBasename = "player"+t+"_";
        var newTeamIdx = await f.createTeam(instance, teamName, playerBasename, k.MaxPlayersInTeam, f.createAlineacion(4,3,3));
        teamsIdx.push(newTeamIdx);
        }
    const blockFirstGame = 100;        
    const blocksBetweenGames = 10;
    await instance.test_createLeague(teamsIdx, blockFirstGame, blocksBetweenGames);
    
  });
  
});

