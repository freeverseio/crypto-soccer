const cryptoSoccer = artifacts.require("Testing");
var k = require('../jsCommons/constants.js');
var f = require('../jsCommons/functions.js');

contract('Leagues', function(accounts) {

    var instance;
    console.log('Funds in the source account:');
    console.log(web3.eth.getBalance(web3.eth.accounts[0]).toNumber()/web3.toWei(1, "ether"));

    it("creates a single contract and computes the gas cost of deploying GameEngine", async () => {
        instance = await cryptoSoccer.new();
        var receipt = await web3.eth.getTransactionReceipt(instance.transactionHash);
        console.log("League\n\tdeployment cost: ", receipt.gasUsed, "\n\tcontract address:", receipt.contractAddress)
        assert.isTrue(receipt.gasUsed > 2000000);
    });

    it("creates 4 teams and puts them into a league. Checks a number of indicators.", async () => {
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
        var nLeagues = await instance.test_getNLeaguesCreated.call();
        assert.equal(nLeagues.toNumber(),1);
        var leagueIdx = nLeagues-1;
        var nTeamsInLeague = await instance.test_getNTeamsInLeague.call(leagueIdx);
        assert.equal(nTeamsInLeague.toNumber(),4);
    });

    it("plays one round of a league, and checks that written results are as expected", async () => {
        var nLeagues = await instance.test_getNLeaguesCreated.call();
        var leagueIdx = nLeagues-1;
        var round = 0;
        var seed = 1;
        await instance.test_playRound(leagueIdx, round, seed);
        var nTeams = await instance.test_getNTeamsInLeague.call(leagueIdx) ;
        var result;
        var info = "LeagueIdx " + leagueIdx + ", with nTeams=" + nTeams + ". RESULTS for round 0:";
        for (var game = 0; game<nTeams/2; game++){
            info += "Game " + game + ": ";
            result = await instance.test_getWrittenResult.call(leagueIdx, nTeams, round, game);
            result = result.toNumber();
            if (result == k.Undef) { info += " UNDEF";};
            if (result == k.HomeWins) { info += " HomeWins";};
            if (result == k.AwayWins) { info += " AwayWins";};
            if (result == k.Tie) { info += " Tie";};
        } 
        expectedInfo="LeagueIdx 0, with nTeams=4. RESULTS for round 0:Game 0:  TieGame 1:  HomeWins";
        assert.equal(info,expectedInfo);
    });


});