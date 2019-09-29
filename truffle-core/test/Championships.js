const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Championships = artifacts.require('Championships');
const Engine = artifacts.require('Engine');

contract('Championships', (accounts) => {
    const tactic442 = 0;
    const tactic433 = 1;

    const createTeamStateFromSinglePlayer = async (skills, engine) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1321312,
            [potential = 3,
            forwardness = 3,
            leftishness = 2,
            aggressiveness = 0],
            alignedLastHalf = false,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0
        ).should.be.fulfilled;

        teamState = []
        for (player = 0; player < PLAYERS_PER_TEAM_MAX; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };
    
    const createLeagueStateFromSinglePlayer = async (skills, engine) => {
        const teamState = await createTeamStateFromSinglePlayer(skills, engine).should.be.fulfilled;
        leagueState = []
        for (team = 0; team < TEAMS_PER_LEAGUE; team++) {
            leagueState.push(teamState)
        }
        return leagueState;
    };
    
    function transpose(a) {
        return Object.keys(a[0]).map(function(c) {
            return a.map(function(r) { return r[c]; });
        });
    }
    

    beforeEach(async () => {
        champs = await Championships.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        await champs.setEngineAdress(engine.address).should.be.fulfilled;
        TEAMS_PER_LEAGUE = await champs.TEAMS_PER_LEAGUE().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await champs.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        MATCHDAYS = await champs.MATCHDAYS().should.be.fulfilled;
        MATCHES_PER_DAY = await champs.MATCHES_PER_DAY().should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine);
    });

    it('check initial constants', async () =>  {
        engine = 0;
        MATCHDAYS.toNumber().should.be.equal(14);
        MATCHES_PER_DAY.toNumber().should.be.equal(4);
        TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    });

    it('get teams for groups', async () => {
        teams = await champs.getTeamsInGroup(groupIdx = 0).should.be.fulfilled;
        teamsExpected = [ 0, 8, 16, 24, 32, 40, 48, 56 ]
        for (t = 0; t < teams.length; t++) {
            teams[t].toNumber().should.be.equal(teamsExpected[t])
        }
        teams = await champs.getTeamsInGroup(groupIdx = 15).should.be.fulfilled;
        teamsExpected = [71, 79, 87, 95, 103, 111, 119, 127 ]
        for (t = 0; t < teams.length; t++) {
            teams[t].toNumber().should.be.equal(teamsExpected[t])
        }
    });

    it('get teams for match in wrong day', async () => {
        matchIdxInDay = 0; 
        day = MATCHDAYS-1; 
        await champs.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        day = MATCHDAYS; 
        await champs.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    });

    it('get teams for match in wrong match in day', async () => {
        day = 0;
        matchIdxInDay = MATCHES_PER_DAY-1;
        await champs.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        matchIdxInDay = MATCHES_PER_DAY;
        await champs.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        day = 0;
        matchIdxInDay = 0;
        teams = await champs.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        day = Math.floor(MATCHDAYS/2);
        teams = await champs.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });
    
    // it('calculate a day in a league', async () => {
    //     day = 0;
    //     verseSeed = 0;
    //     leagueAll50 = await createLeagueStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
    //     leagueTacticsIds = Array(TEAMS_PER_LEAGUE.toNumber()).fill(tactic442);
    //     result = await champs.computeMatchday(day, leagueAll50, leagueTacticsIds, verseSeed).should.be.fulfilled;
    //     result.length.should.be.equal(MATCHES_PER_DAY * 2);
    //     expectedScores      = [ 0, 1, 0, 0, 1, 5, 3, 1 ]
    //     actualScores    = Array.from(new Array(result.length), (x,i) => result[i].toNumber());
    //     // console.log(actualScores);
    //     for (idx = 0; idx < 2 * MATCHES_PER_DAY; idx++){
    //         result[idx].toNumber().should.be.equal(expectedScores[idx]);
    //     }
    //     day = 3;
    //     verseSeed = 432;
    //     leagueAll50 = await createLeagueStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
    //     leagueTacticsIds = Array(TEAMS_PER_LEAGUE.toNumber()).fill(tactic442);
    //     result = await champs.computeMatchday(day, leagueAll50, leagueTacticsIds, verseSeed).should.be.fulfilled;
    //     result.length.should.be.equal(MATCHES_PER_DAY * 2);
    //     expectedScores      = [ 0, 3, 1, 3, 1, 0, 1, 1 ]
    //     actualScores    = Array.from(new Array(result.length), (x,i) => result[i].toNumber());
    //     // console.log(actualScores);
    //     for (idx = 0; idx < 2 * MATCHES_PER_DAY; idx++){
    //         result[idx].toNumber().should.be.equal(expectedScores[idx]);
    //     }
    // });

});