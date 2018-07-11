const cryptoSoccer = artifacts.require("GameEngine");
const maxPlayersPerTeam = 11;
const playerRoles433 = [0,1,1,1,1,2,2,2,3,3,3];
const playerRoles442 = [0,1,1,1,1,2,2,2,2,3,3];
const playerRoles541 = [0,1,1,1,1,1,2,2,2,2,3];
const playerRoles631 = [0,1,1,1,1,1,1,2,2,2,3];
const playerRoles640 = [0,1,1,1,1,1,1,2,2,2,2];
const playerRoles451 = [0,1,1,1,1,2,2,2,2,2,3];

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
    nCreatedPlayers = await instance.getNCreatedPlayers.call();
    assert.equal(nCreatedPlayers,0);
    // TODO: derive from the name and the mapping.
    await createTeam(instance, "Mataro", "Bogarde", maxPlayersPerTeam, 0, playerRoles433);
    await printTeamPlayers(0, instance);
    assert.equal(nCreatedPlayers,maxPlayersPerTeam);
  });


  // TODO: add test that you cannot create 2 teams with same name
  it("creates another team and plays a game. With this seed, it checks that the result is 2-2", async () => {
    // await createTeam(instance, "Sevilla", "Navas", maxPlayersPerTeam, 1, playerRoles433);
    await createTestTeam(
      instance, 
      "Sevilla", 
      "Navas", 
      maxPlayersPerTeam, 
      1, 
      [220, 50,50,50,50,50], // age, defense, speed, pass, shoot, endurance
      playerRoles433
    );
    await printTeamPlayers(1, instance);
    var goals = await instance.playGame.call(0,1,232);
    console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    assert.isTrue(goals[0].toNumber()==2);
    assert.isTrue(goals[1].toNumber()==2);
  });

  it("plays a game using a transation, not a call, to compute gas cost", async () => {
    var goals = await instance.playGame(0,1,232);
  });

  it("plays lots of games and checks total goals", async () => {
    var goalsTeam1 = 0;
    var goalsTeam2 = 0;
    nGames = 5;
    console.log("Playing " + nGames + " games");
    for (var game=0; game<nGames; game++) {
      var goals = await instance.playGame.call(0,1,game+1);
      goalsTeam1 += goals[0].toNumber();
      goalsTeam2 += goals[1].toNumber();
      console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    }
    console.log("Total Goals: " + goalsTeam1 + " - " + goalsTeam2);
    assert.isTrue(goalsTeam1==6);
    assert.isTrue(goalsTeam2==13);
  });
});

function catchPlayerIdxFromEvent(logs) {
    var playerIdx = -1;
    for (var i = 0; i < logs.length; i++) {
        var log = logs[i];
        if (log.event == "PlayerCreation") {
            //console.log("...created player with idx = " + log.args.playerIdx.toNumber());
            playerIdx = log.args.playerIdx.toNumber();
        }
    }
    return playerIdx;
}

async function createTeam(instance, teamName, playerBasename, maxPlayersPerTeam, teamIdx, playerRoles ) {
  // TODO: derive from the name and the mapping
  console.log("creating team: " + teamName);
  await instance.createTeam(teamName);
  const userChoice=1;

  for (var p=0; p<maxPlayersPerTeam; p++) {
      thisName = playerBasename + p.toString();
      var tx = await instance.createRandomPlayer(thisName,teamIdx,userChoice,p,playerRoles[p]);
      var playerIdx = catchPlayerIdxFromEvent(tx.logs);
      assert( playerIdx >= 0 );
  }
  nCreatedPlayers = await instance.getNCreatedPlayers.call();
  console.log('Final nPlayers in the entire game = ' + nCreatedPlayers);
}

async function printTeamPlayers(teamIdx, instance) {
  var state = await instance.getSkillsOfPlayersInTeam.call(teamIdx);
  var totals = Array(7).fill(0);
  console.log("Players in team " + teamIdx);
  for (var p=0;p<maxPlayersPerTeam;p++) {
    process.stdout.write("Player " + p + ": ");
    for (var sk=0;sk<7;sk++) {
      if (sk==0) state[p][sk] = unixMonthToAge(state[p][sk]);
      totals[sk] += parseInt(state[p][sk]);
      process.stdout.write(skillNames[sk] + "= " + state[p][sk] + "  ");
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
  await instance.createTeam(teamName);

  for (var p=0; p<maxPlayersPerTeam; p++) {
      thisName = playerBasename + p.toString();
      var tx = await instance.createTestPlayer(
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
      var playerIdx = catchPlayerIdxFromEvent(tx.logs);
      assert( playerIdx >= 0 );
  }
  nCreatedPlayers = await instance.getNCreatedPlayers.call();
  console.log('Final nPlayers in the entire game = ' + nCreatedPlayers);
}
