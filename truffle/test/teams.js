const cryptoSoccer = artifacts.require("Testing");
var k = require('../jsCommons/constants.js');
var f = require('../jsCommons/functions.js');

const skillNames = ["Age","Defense","Speed","Pass","Shoot","Endurance","Role"];

contract('Teams', function(accounts) {

  var instance;
  var nTotalPlayers=0;
  var sourceBalance;
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
    // TODO: derive from the name and the mapping.
    var newTeamIdx = await f.createTeam(instance, "Mataro", "Bogarde", k.MaxPlayersInTeam, f.createAlineacion(4,3,3));
    await printTeamPlayers(newTeamIdx, instance);
    assert.equal(nCreatedPlayers,k.MaxPlayersInTeam+1);
  });

  it("creates another team and plays a game. With this seed, it checks that the result is 1-3", async () => {
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
    var goals = await playGame(instance, 0, 1, 18, 232);
    console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    assert.isTrue(goals[0].toNumber()==1);
    assert.isTrue(goals[1].toNumber()==3);
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
    var goals = await playGame(instance, 0, 1, 18, 232);
  });

  it("plays lots of games and checks total goals", async () => {
    var goalsTeam1 = 0;
    var goalsTeam2 = 0;
    nGames = 5;
    console.log("Playing " + nGames + " games");
    for (var game=0; game<nGames; game++) {
      var goals = await playGame(instance, 0, 1, 18, game+1);
      goalsTeam1 += goals[0].toNumber();
      goalsTeam2 += goals[1].toNumber();
      console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    }
    console.log("Total Goals: " + goalsTeam1 + " - " + goalsTeam2);
    assert.isTrue(goalsTeam1==6);
    assert.isTrue(goalsTeam2==14);
  });
});


async function playGame(instance, teamIdx1, teamIdx2, nRounds, rndSeed)
{
  var rndNums = await f.getRandomNumbers(instance, nRounds, rndSeed);
  var goals = await instance.test_playGame.call(teamIdx1, teamIdx2, rndNums[0], rndNums[1], rndNums[2], rndNums[3]);
  return goals;
}


async function printTeamPlayers(teamIdx, instance) {
//  var state = await instance.test_getSkillsOfPlayersInTeam.call(teamIdx);
  var totals = Array(k.NumStates).fill(0);
  console.log("Players in team " + teamIdx);
  for (var p=0;p<k.MaxPlayersInTeam;p++) {
    process.stdout.write("Player " + p + ": ");
    serializedSkills = await instance.test_getSkill(teamIdx, p);
    decodedSkills = await instance.test_decode(k.NumSkills, serializedSkills, k.BitsPerState);
    //console.log('skills:' + decodedSkills)
    for (var sk=0;sk<k.NumStates;sk++) {
      if (sk==0) totals[0] += unixMonthToAge(decodedSkills[0]);
      else totals[sk] += parseInt(decodedSkills[sk]);
      process.stdout.write(skillNames[sk] + "= " + decodedSkills[sk] + "  ");
    }
    process.stdout.write("\n");
  }
  console.log("Totals: " + totals);
}



function unixMonthToAge(unixMonthOfBirth) {
  // in July 2018, we are at month 582 after 1970.
  age = (582 - unixMonthOfBirth)/12;
  return parseInt(age*10)/10;
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
  // TODO: derive from the name and the mapping
  console.log("creating team: " + teamName);
  await instance.test_createTeam(teamName);

  for (var p=0; p<maxPlayersPerTeam; p++) {
      thisName = playerBasename + p.toString();
      var tx = await instance.test_createUnbalancedPlayer(
          thisName,
          teamIdx,
          p,
          skills[0], // monthOfBirthAfterUnixEpoch
          skills[1], // defense
          skills[2], // speed
          skills[3], // pass
          skills[4], // shoot
          skills[5], // endurance
          playerRoles[p]
        );
      var playerIdx = f.catchPlayerIdxFromEvent(tx.logs);
      assert( playerIdx >= 0 );
  }
  nCreatedPlayers = await instance.test_getNCreatedPlayers.call();
  console.log('Final nPlayers in the entire game = ' + nCreatedPlayers);
}

