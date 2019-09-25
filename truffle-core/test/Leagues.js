const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');

contract('Leagues', (accounts) => {
    const tactic442 = 0;
    const tactic433 = 1;

    const createTeamStateFromSinglePlayer = async (skills, engine) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1,
            potential = 3,
            forwardness = 3,
            leftishness = 2,
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
        leagues = await Leagues.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        await leagues.init().should.be.fulfilled;
        await leagues.setEngineAdress(engine.address).should.be.fulfilled;
        TEAMS_PER_LEAGUE = await leagues.TEAMS_PER_LEAGUE().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await leagues.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        MATCHDAYS = await leagues.MATCHDAYS().should.be.fulfilled;
        MATCHES_PER_DAY = await leagues.MATCHES_PER_DAY().should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine);
    });

    it('check initial constants', async () =>  {
        engine = 0;
        MATCHDAYS.toNumber().should.be.equal(14);
        MATCHES_PER_DAY.toNumber().should.be.equal(4);
        TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    });

    it('get teams for match in wrong day', async () => {
        matchIdxInDay = 0; 
        day = MATCHDAYS-1; 
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        day = MATCHDAYS; 
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    });

    it('get teams for match in wrong match in day', async () => {
        day = 0;
        matchIdxInDay = MATCHES_PER_DAY-1;
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        matchIdxInDay = MATCHES_PER_DAY;
        await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        day = 0;
        matchIdxInDay = 0;
        teams = await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        day = Math.floor(MATCHDAYS/2);
        teams = await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });
    
    it('calculate a day in a league', async () => {
        day = 0;
        verseSeed = 0;
        leagueAll50 = await createLeagueStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        leagueTacticsIds = Array(TEAMS_PER_LEAGUE.toNumber()).fill(tactic442);
        result = await leagues.computeMatchday(day, leagueAll50, leagueTacticsIds, verseSeed).should.be.fulfilled;
        result.scores.length.should.be.equal(MATCHES_PER_DAY * 2);
        expectedScores      = [ 0, 1, 0, 0, 1, 5, 3, 1 ]
        actualScores    = Array.from(new Array(result.scores.length), (x,i) => result.scores[i].toNumber());
        // console.log(actualScores);
        for (idx = 0; idx < 2 * MATCHES_PER_DAY; idx++){
            result.scores[idx].toNumber().should.be.equal(expectedScores[idx]);
        }
        day = 3;
        verseSeed = 432;
        leagueAll50 = await createLeagueStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        leagueTacticsIds = Array(TEAMS_PER_LEAGUE.toNumber()).fill(tactic442);
        result = await leagues.computeMatchday(day, leagueAll50, leagueTacticsIds, verseSeed).should.be.fulfilled;
        result.scores.length.should.be.equal(MATCHES_PER_DAY * 2);
        expectedScores      = [ 0, 3, 1, 3, 1, 0, 1, 1 ]
        actualScores    = Array.from(new Array(result.scores.length), (x,i) => result.scores[i].toNumber());
        // console.log(actualScores);
        for (idx = 0; idx < 2 * MATCHES_PER_DAY; idx++){
            result.scores[idx].toNumber().should.be.equal(expectedScores[idx]);
        }
    });

});