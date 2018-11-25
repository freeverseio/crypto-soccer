require('chai')
    .use(require('chai-as-promised'))
    .should();

const TeamFactory = artifacts.require("TeamFactoryMock");
const Testing = artifacts.require("Testing");

var k = require('../jsCommons/constants.js');
var f = require('../jsCommons/functions.js');

contract('Leagues', function(accounts) {
    let teamFactory;
    let instance;
    console.log('Funds in the source account:');
    console.log(web3.eth.getBalance(web3.eth.accounts[0]).toNumber() / web3.toWei(1, "ether"));

    it('deployment', async () => {
        teamFactory = await TeamFactory.new().should.be.fulfilled;
        instance = await Testing.new(teamFactory.address).should.be.fulfilled;
    });

    // it("creates a single contract and computes the gas cost of deploying GameEngine", async () => {
    //     var receipt = await web3.eth.getTransactionReceipt(instance.transactionHash);
    //     console.log("League\n\tdeployment cost: ", receipt.gasUsed, "\n\tcontract address:", receipt.contractAddress)
    //     assert.isTrue(receipt.gasUsed > 2000000);
    // });

    it("creates 4 teams and puts them into a league. Checks a number of indicators.", async () => {
        nTeams = 4;
        teamsIdx = [];
        for (var t = 0; t < nTeams; t++) {
            var teamName = "team"+t;
            var playerBasename = "player"+t+"_";
            var newTeamIdx = await f.createTeam(teamFactory, teamName, playerBasename, k.MaxPlayersInTeam, f.createAlineacion(4,3,3));
            teamsIdx.push(newTeamIdx);
        }
        const blockFirstGame = 100;        
        const blocksBetweenGames = 10;
        await instance.test_createLeague(teamsIdx, blockFirstGame, blocksBetweenGames);   
        var nLeagues = await instance.test_getNLeaguesCreated.call();
        assert.equal(nLeagues.toNumber(),1);
        var leagueIdx = nLeagues-1;
        var nTeamsInLeague = await instance.test_getNTeamsInLeague.call(leagueIdx);
        assert.equal(nTeamsInLeague.toNumber(),4);
    });

    it("plays one round of a league, and checks that written results are as expected", async () => {
        var nLeagues = await instance.test_getNLeaguesCreated.call().should.be.fulfilled;
        var leagueIdx = nLeagues-1;
        var round = 0;
        var seed = 1;
        await instance.test_playRound(leagueIdx, round, seed).should.be.fulfilled;
        var nTeams = await instance.test_getNTeamsInLeague.call(leagueIdx).should.be.fulfilled;
        var result;
        var info = "LeagueIdx " + leagueIdx + ", with nTeams=" + nTeams + ". RESULTS for round 0:";
        for (var game = 0; game<nTeams/2; game++){
            info += "Game " + game + ": ";
            result = await instance.test_getWrittenResult.call(leagueIdx, nTeams, round, game).should.be.fulfilled;
            result = result.toNumber();
            if (result == k.Undef) { info += " UNDEF";};
            if (result == k.HomeWins) { info += " HomeWins";};
            if (result == k.AwayWins) { info += " AwayWins";};
            if (result == k.Tie) { info += " Tie";};
        } 
        expectedInfo="LeagueIdx 0, with nTeams=4. RESULTS for round 0:Game 0:  TieGame 1:  HomeWins";
        assert.equal(info,expectedInfo);
    });

    it("plays all remaining rounds of a league", async () => {
        var nLeagues = await instance.test_getNLeaguesCreated.call();
        var leagueIdx = nLeagues-1;
        var info = "LeagueIdx " + leagueIdx;
        for (var round = 0; round < 2 * (nTeams-1); round++) {
            seed = round;
            await instance.test_playRound(leagueIdx, round, seed);
            var info = "Round " + round + ": ";
            for (var game = 0; game<nTeams/2; game++){
                info += "\n  Game " + game + ":  ";
                teams = await instance.test_teamsInGame.call(round, game, nTeams);
                // console.log(teams)
                info += teams[0].toNumber() + " vs " + teams[1].toNumber() + " => ";
                result = await instance.test_getWrittenResult.call(leagueIdx, nTeams, round, game);
                result = result.toNumber();
                if (result == k.Undef) { info += " UNDEF";};
                if (result == k.HomeWins) { info += " HomeWins";};
                if (result == k.AwayWins) { info += " AwayWins";};
                if (result == k.Tie) { info += " Tie";};
                info += "    ";
            } 
            console.log(info);
        }
    });

        it("creates another team and plays a game. With this seed, it checks that the result is 1-5", async () => {
    await createTestTeam(
      teamFactory,
      "Sevilla",
      "Navas",
      k.MaxPlayersInTeam,
      1,
      [220, 50,50,50,50,50], // age, defense, speed, pass, shoot, endurance
      f.createAlineacion(4,3,3)
    );
    seed = 232;
    var goals = await instance.test_playGame.call(0, 1, seed);
    
    console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    assert.isTrue(goals[0].toNumber()==1);
    assert.isTrue(goals[1].toNumber()==5);
  });


  it("reads the game events of the previous game", async () => {
    // plays same game as before, but now as a transaction, 
    // so that events are generated (and stored in the BChain)
    var tx = await instance.test_playGame(0, 1, seed);
    // reads the gameID, which is basically the hash(teamIdx1, teamIdx2,seed)
    var gameId = await instance.test_getGameId(0, 1, seed);
    // catches events and prints them out
    var gameEvents = f.catchGameResults(tx.logs,gameId) ;
    printGameEvents(gameEvents);
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
      var goals = await instance.test_playGame.call(0, 1, seed);
      goalsTeam1 += goals[0].toNumber();
      goalsTeam2 += goals[1].toNumber();
      console.log("Goals: " + goals[0].toNumber() + " - " + goals[1].toNumber());
    }
    console.log("Total Goals: " + goalsTeam1 + " - " + goalsTeam2);
    assert.isTrue(goalsTeam1==6);
    assert.isTrue(goalsTeam2==5);
  });
});

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

function printGameEvents(gameEvents) {
    console.log("EVENTS: ");
    for (var r = 0; r < k.RoundsPerGame; r++) {
        // we add a bit of noise so that events are not always at minute 5,10,15...
        var rndNoise = Math.round(-2+Math.floor(Math.random() * 4)); 
        var thisMinute = (r+1)*5 + rndNoise;
        var t = f.getEntryForAGivenRound(gameEvents.teamThatAttacks,r);
        console.log("Min " +thisMinute + ": Opportunity for team " + t[1] + "...");
        var result = f.getEntryForAGivenRound(gameEvents.shootResult,r);
        if (result.length==0) { console.log("  ... well tackled by defenders, did not prosper!");}
        else {
            console.log("  ... that leads to a shoot by attacker " + result[2]);
            if (result[1]) { console.log("  ... and GOAAAAL!")} 
            else {console.log("  ... blocked by the goalkeeper!!");} 
        }
    }
}