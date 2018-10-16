const cryptoSoccer = artifacts.require("Testing");
var k = require('../jsCommons/constants.js');
var f = require('../jsCommons/functions.js');

const skillNames = ["Age","Defense","Speed","Pass","Shoot","Endurance","Role"];

contract('Teams', function(accounts) {

  var instance;
  console.log('Funds in the source account:');
  console.log(web3.eth.getBalance(web3.eth.accounts[0]).toNumber()/web3.toWei(1, "ether"));

  // deploy the contract before each test
  //beforeEach(async () => {
  //  instance = await cryptoSoccer.deployed();
  //});

  it("creates a single contract and computes the gas cost of deploying GameEngine", async () => {
    instance = await cryptoSoccer.new();
    var receipt = await web3.eth.getTransactionReceipt(instance.transactionHash);
    console.log("GameEngine\n\tdeployment cost: ", receipt.gasUsed, "\n\tcontract address:", receipt.contractAddress)
    assert.isTrue(receipt.gasUsed > 2000000);
  });


  it("creates an entire team, an checks that we have 11 players at the end", async () => {
    nCreatedPlayers = await instance.test_getNCreatedPlayers.call();
    assert.equal(nCreatedPlayers,1);
    teamName = "Mataro";
    playerBasename = "Bogarde";
    var newTeamIdx = await f.createTeam(instance, teamName, playerBasename, k.MaxPlayersInTeam, f.createAlineacion(4,3,3));
    await printTeamPlayers(newTeamIdx, instance);
    assert.equal(nCreatedPlayers,k.MaxPlayersInTeam+1);
  });

  it("creates another team and plays a game. With this seed, it checks that the result is 1-5", async () => {
    await createTestTeam(
      instance,
      "Sevilla",
      "Navas",
      k.MaxPlayersInTeam,
      1,
      [220, 50,50,50,50,50], // age, defense, speed, pass, shoot, endurance
      f.createAlineacion(4,3,3)
    );
    await printTeamPlayers(1, instance);
    seed = 232;
    var goals = await instance.test_playGame(0, 1, seed);
    console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    assert.isTrue(goals[0].toNumber()==1);
    assert.isTrue(goals[1].toNumber()==5);
  });

  it("creates an empty team, shows crazy stats, checks name is correct", async () => {
    await instance.test_createTeam("Los Cojos");
    var name = await instance.test_getTeamName(2);
    assert.isTrue(name == "Los Cojos");
    await printTeamPlayers(2, instance);
  });
  
  it("checks that we cannot add 2 teams with same name", async () => {
    hasFailed = false;
    try{ 
        await f.createTeam(instance, "Los Cojos", "Reiziger", k.MaxPlayersInTeam, f.createAlineacion(4,3,3));
    } catch (err) {
      // Great, the transaction failed
      hasFailed = true;
    }
    assert.isTrue(hasFailed);
  });
  
  it("plays a game using a transation, not a call, to compute gas cost", async () => {
    seed = 232;
    var goals = await instance.test_playGame(0, 1, seed);
  });

  it("plays lots of games and checks total goals", async () => {
    var goalsTeam1 = 0;
    var goalsTeam2 = 0;
    nGames = 5;
    console.log("Playing " + nGames + " games");
    for (var game=0; game<nGames; game++) {
      seed = game + 1;
      var goals = await instance.test_playGame(0, 1, seed);
      goalsTeam1 += goals[0].toNumber();
      goalsTeam2 += goals[1].toNumber();
      console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    }
    console.log("Total Goals: " + goalsTeam1 + " - " + goalsTeam2);
    assert.isTrue(goalsTeam1==6);
    assert.isTrue(goalsTeam2==5);
  });
});

async function printTeamPlayers(teamIdx, instance) {
  var totals = Array(k.NumStates).fill(0);
  console.log("Players in team " + teamIdx + ":");
  for (var p=0;p<k.MaxPlayersInTeam;p++) {
    info = "Player " + p + ": ";
    serialized = await instance.test_getStatePlayerInTeam(p, teamIdx);
    decoded = await instance.test_decode(k.NumStates, serialized, k.BitsPerState);
    for (var sk=0;sk<k.NumStates;sk++) {
      if (sk==0) totals[0] += f.unixMonthToAge(decoded[0]);
      else totals[sk] += parseInt(decoded[sk]);
      info += skillNames[sk] + "= " + decoded[sk] + "  ";
    }
    console.log(info);
  }
  console.log("Totals: " + totals);
}





async function createTestTeam(
  instance,
  teamName,
  playerBasename,
  maxPlayersPerTeam,
  teamIdx,
  skills,
  playerRoles
  )
{
  console.log("creating team: " + teamName);
  await instance.test_createTeam(teamName);

  for (var p=0; p<maxPlayersPerTeam; p++) {
      thisName = playerBasename + p.toString();
      var tx = await instance.test_createUnbalancedPlayer(
          thisName,
          teamIdx,
          p,
          skills[k.StatBirth], // monthOfBirthAfterUnixEpoch
          skills[k.StatDef],
          skills[k.StatSpeed],
          skills[k.StatPass],
          skills[k.StatShoot],
          skills[k.StatEndur],
          playerRoles[p]
        );
      var playerIdx = f.catchPlayerIdxFromEvent(tx.logs);
      assert( playerIdx >= 0 );
  }
  nCreatedPlayers = await instance.test_getNCreatedPlayers.call();
  console.log('Final nPlayers in the entire game = ' + nCreatedPlayers);
}

