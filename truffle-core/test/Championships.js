const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');

const Championships = artifacts.require('Championships');
const Engine = artifacts.require('Engine');

contract('Championships', (accounts) => {
    const now = 1570147200; // this number has the property that 7*nowFake % (SECS_IN_DAY) = 0 and it is basically Oct 3, 2019
    const dayOfBirth21 = secsToDays(now) - 21*365/7; // = exactly 17078, no need to round
    const subLastHalf = false;
    const seed = web3.utils.toBN(web3.utils.keccak256("32123"));

    const it2 = async(text, f) => {};
    
    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }

    const createTeamStateFromSinglePlayer = async (skills, engine, forwardness = 3, leftishness = 2, alignedEndOfLastHalfTwoVec = [false, false]) => {
        teamState = []
        sumSkills = skills.reduce((a, b) => a + b, 0);
        var playerStateTemp = await engine.encodePlayerSkills(
            skills, dayOfBirth21, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[0], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 0; player < 11; player++) {
            teamState.push(playerStateTemp)
        }

        playerStateTemp = await engine.encodePlayerSkills(
            skills, dayOfBirth21, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[1], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 11; player < PLAYERS_PER_TEAM_MAX; player++) {
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
    
    function getRand(seed, min, max) {
        return min + (2**Math.abs(Math.floor(Math.sin(seed + 324212) * 24))) % (max - min + 1)
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

    it('computeTeamRankingPoints with no previous points', async () =>  {
        // teamSkills = 5*25
        // rankingPoints = 5*25*100 + ( (6000*2/10000) - 10 ) * 900 = 5*25*100 - 9*900 = 4400
        // 10W SK + SK0 (I P0 + (10-I)P1 - 100) = 10* 100 * 5 * 25 + 18*50 *(6*2-100) = 45800
        result = await champs.computeTeamRankingPoints(teamStateAll1, leagueRanking = 0, prevPerfPoints = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(45800);
        
    });

    it('computeTeamRankingPoints with previous points', async () =>  {
        // teamSkills = 50*25
        // rankingPoints = 5*25*100 + ( (6000*2/10000) - 10 ) * 900 = 5*25*100 - 9*900 = 4400
        // 10W SK + SK0 (I P0 + (10-I)P1 - 100) = 10* 100 * 50 * 25 + 18*50 *(4*50000+6 * 20 -100) = 45800
        result = await champs.computeTeamRankingPoints(teamStateAll50, leagueRanking = 7, prevPerfPoints = 50000).should.be.fulfilled;
        result.toNumber().should.be.equal(181268000);
        
    });
    return
    it('computeLeagueLeaderBoard almost no clashes', async () =>  {
        MATCHES_PER_LEAGUE = 56;
        matchDay = 13;
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [getRand(2*i, 0, 12), getRand(2*i+1, 0, 12)]);
        result = await champs.computeLeagueLeaderBoard(results, matchDay, seed).should.be.fulfilled;
        expectedPoints =  [26000000000, 24000000000, 23000000000, 21000000000, 15000000000, 14000000000, 12001081423, 12000075754];
        expectedRanking = [ 4, 1, 2, 5, 7, 0, 6, 3 ];
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, verbose = false);
        debug.compareArrays(result.points, expectedPoints, toNum = true, verbose = false);
    });
    
    it('computeLeagueLeaderBoard many clashes', async () =>  {
        MATCHES_PER_LEAGUE = 56;
        matchDay = 13;
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [getRand(2*i+1, 0, 2), getRand(2*i+3, 0, 12)]);
        result = await champs.computeLeagueLeaderBoard(results, matchDay, seed).should.be.fulfilled;
        expectedPoints =  [27000000000, 22000052441, 22000051610, 18000000000, 16002047402, 16001047423, 16000049802, 15000000000];
        expectedRanking = [ 7, 2, 1, 5, 4, 3, 0, 6 ];
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, verbose = false);
        debug.compareArrays(result.points, expectedPoints, toNum = true, verbose = false);
    });

    it('computeLeagueLeaderBoard all clashes', async () =>  {
        MATCHES_PER_LEAGUE = 56;
        matchDay = 13;
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [getRand(2*i+1, 0, 1), getRand(2*i+3, 0, 1)]);
        result = await champs.computeLeagueLeaderBoard(results, matchDay, seed).should.be.fulfilled;
        expectedPoints =  [ 13000000802, 13000000754, 13000000610, 13000000441, 13000000423, 13000000402, 13000000389, 13000000110 ];
        expectedRanking = [ 0, 2, 7, 4, 3, 1, 6, 5 ];
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, verbose = false);
        debug.compareArrays(result.points, expectedPoints, toNum = true, verbose = false);
    });

    it('check initial constants', async () =>  {
        MATCHDAYS.toNumber().should.be.equal(14);
        MATCHES_PER_DAY.toNumber().should.be.equal(4);
        TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    });

    it('getTeamsInCupPlayoffMatch', async () => {
        teamsExpected = [0,7,9,14,4,11,13,18,8,15,17,22,12,19,21,26,16,23,25,30,20,27,29,34,24,31,33,38,28,35,37,42,32,39,41,46,36,43,45,50,40,47,49,54,44,51,53,58,48,55,57,62,52,59,61,2,56,63,1,6,60,3,5,10];
        for (t = 0; t < 32; t++) {
            team = await champs.getTeamsInCupPlayoffMatch(matchIdxInDay = t).should.be.fulfilled;
            team[0].toNumber().should.be.equal(teamsExpected[2*t]);
            team[1].toNumber().should.be.equal(teamsExpected[2*t+1]);
        }
        // check that all teams are included, and only once (e.g. by sorting and requiring monotonic growing series)
        teamsExpected.sort((a, b) => a - b);
        for (t = 1; t < 64; t++) {
            (team[0]*0 + teamsExpected[t] > teamsExpected[t-1]).should.be.equal(true);
        }
    });
    
    it('get all teams for groups', async () => {
        teamsExpected = [ 0, 8, 16, 24, 32, 40, 48, 56 ]
        for (t = 0; t < teamsExpected.length; t++) {
            team = await champs.getTeamIdxInCup(groupIdx = 0, posInGroup = t).should.be.fulfilled;
            team.toNumber().should.be.equal(teamsExpected[t]);
            result = await champs.getGroupAndPosInGroup(team.toNumber()).should.be.fulfilled;
            result[0].toNumber().should.be.equal(groupIdx);
            result[1].toNumber().should.be.equal(posInGroup);
        }
        teamsExpected = [71, 79, 87, 95, 103, 111, 119, 127 ]
        for (t = 0; t < teamsExpected.length; t++) {
            team = await champs.getTeamIdxInCup(groupIdx = 15, posInGroup = t).should.be.fulfilled;
            team.toNumber().should.be.equal(teamsExpected[t])
            result = await champs.getGroupAndPosInGroup(team.toNumber()).should.be.fulfilled;
            result[0].toNumber().should.be.equal(groupIdx);
            result[1].toNumber().should.be.equal(posInGroup);
        }
    });

    it('get all teams for particular matches', async () => {
        teams = await champs.getTeamsInCupLeagueMatch(groupIdx = 0, day = 0, matchIdxInDay = 0).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(8);
        teams = await champs.getTeamsInCupLeagueMatch(groupIdx = 0, day = day = Math.floor(MATCHDAYS/2), matchIdxInDay = 0).should.be.rejected;
        teams = await champs.getTeamsInCupLeagueMatch(groupIdx = 15, day = 0, matchIdxInDay = 0).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(71);
        teams[1].toNumber().should.be.equal(79);
    });

    it('get teams for match in wrong day', async () => {
        await champs.getTeamsInLeagueMatch(day = MATCHDAYS-1, matchIdxInDay = 0).should.be.fulfilled;
        await champs.getTeamsInLeagueMatch(day = MATCHDAYS, matchIdxInDay = 0).should.be.rejected;
    });

    it('get teams for match in wrong match in day', async () => {
        await champs.getTeamsInLeagueMatch(day = 0, matchIdxInDay = MATCHES_PER_DAY-1).should.be.fulfilled;
        await champs.getTeamsInLeagueMatch(day = 0, matchIdxInDay = MATCHES_PER_DAY).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        teams = await champs.getTeamsInLeagueMatch(day = 0, matchIdxInDay = 0).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        teams = await champs.getTeamsInLeagueMatch(day = Math.floor(MATCHDAYS/2), matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });
    

});