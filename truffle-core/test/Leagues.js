/*
 Tests for all functions in Leauges.sol
*/
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const deployUtils = require('../utils/deployUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Leagues = artifacts.require('Leagues');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Challenges = artifacts.require('Challenges');
const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');
const EngineApplyBoosters = artifacts.require('EngineApplyBoosters');

const UniverseInfo = artifacts.require('UniverseInfo');
const EncodingSkills = artifacts.require('EncodingSkills');
const EncodingState = artifacts.require('EncodingState');
const EncodingSkillsSetters = artifacts.require('EncodingSkillsSetters');
const UpdatesBase = artifacts.require('UpdatesBase');

contract('Leagues', (accounts) => {
    const inheritedArtfcts = [UniverseInfo, EncodingSkills, EncodingState, EncodingSkillsSetters, UpdatesBase];
    const now = 1570147200; // this number has the property that 7*nowFake % (SECS_IN_DAY) = 0 and it is basically Oct 3, 2019
    const dayOfBirth21 = secsToDays(now) - 21*365/7; // = exactly 17078, no need to round
    const subLastHalf = false;
    const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
    const INIT_TZ = 4;
    const TWO_TO_28 = 2**28;
    const MAX_TEAMIDX_IN_COUNTRY = TWO_TO_28 - 1;
    const MATCHES_PER_LEAGUE = 56;
    const TEN_TO_13 = web3.utils.toBN(10**13);
         
    const it2 = async(text, f) => {};
    
    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }

    const createTeamStateFromSinglePlayer = async (skills, engine, forwardness = 3, leftishness = 2, alignedEndOfLastHalfTwoVec = [false, false]) => {
        teamState = []
        sumSkills = skills.reduce((a, b) => a + b, 0);
        var playerStateTemp = await assets.encodePlayerSkills(
            skills, dayOfBirth21, gen = 0, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[0], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 0; player < 11; player++) {
            teamState.push(playerStateTemp)
        }

        playerStateTemp = await assets.encodePlayerSkills(
            skills, dayOfBirth21, gen = 0, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[1], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 11; player < PLAYERS_PER_TEAM_MAX; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };
    
    function getRand(seed, min, max) {
        return min + (2**Math.abs(Math.floor(Math.sin(seed + 324212) * 24))) % (max - min + 1)
    }
    
    beforeEach(async () => {
        defaultSetup = deployUtils.getDefaultSetup(accounts);
        owners = defaultSetup.owners;
        depl = await deployUtils.deploy(owners, Proxy, Assets, Market, Updates, Challenges, inheritedArtfcts);
        [proxy, assets, market, updates, challenges] = depl;
        await deployUtils.setProxyContractOwners(proxy, assets, owners, owners.company).should.be.fulfilled;
        await assets.initSingleTZ(INIT_TZ, {from: owners.COO}).should.be.fulfilled;

        constants = await ConstantsGetters.new().should.be.fulfilled;
        leagues = await Leagues.new(assets.address).should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        applyBoosters = await EngineApplyBoosters.new().should.be.fulfilled;
        engine = await Engine.new(precomp.address, applyBoosters.address).should.be.fulfilled;

        TEAMS_PER_LEAGUE = await constants.get_TEAMS_PER_LEAGUE().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await constants.get_PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        MATCHDAYS = await leagues.MATCHDAYS().should.be.fulfilled;
        MATCHES_PER_DAY = await leagues.MATCHES_PER_DAY().should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine);
    });

    it('computeTeamRankingPoints with no previous points', async () =>  {
        // teamSkills = 5*25
        // rankingPoints = 5*25*100 + ( (6000*2/10000) - 10 ) * 900 = 5*25*100 - 9*900 = 4400
        // 10W SK + SK0 (I P0 + (10-I)P1 - 100) = 10* 20 * 5 * 25 + 18*1000 *(6*20-100) = 43000
        result = await leagues.computeTeamRankingPointsPure(teamStateAll1, leagueRanking = 0, prevPerfPoints = 0).should.be.fulfilled;
        result[0].toNumber().should.be.equal(385000);
        // prevPerfPoints = 0.6 * 20 = 12
        result[1].toNumber().should.be.equal(12);
    });

    it('computeTeamRankingPoints with previous points', async () =>  {
        // teamSkills = 5*50*25
        // rankingPoints = 5*25*100 + ( (6000*2/10000) - 10 ) * 900 = 5*25*100 - 9*900 = 4400
        // 10W SK + SK0 (I P0 + (10-I)P1 - 100) = 10* 20 * 5*50 * 25 + 18*1000 *(4*10+ 6 * 2 -100) = 6206800
        result = await leagues.computeTeamRankingPointsPure(teamStateAll50, leagueRanking = 7, prevPerfPoints = 10).should.be.fulfilled;
        result[0].toNumber().should.be.equal(386000);
        // prevPerfPoints = 0.6 * 2 + 0.4 * 10 = 5.2
        result[1].toNumber().should.be.equal(5);
    });

    it('computeTeamRankingPoints with previous points and non-null teamId', async () =>  {
        // teamSkills = 5*50*25
        // rankingPoints = 5*25*100 + ( (6000*2/10000) - 10 ) * 900 = 5*25*100 - 9*900 = 4400
        // 10W SK + SK0 (I P0 + (10-I)P1 - 100) = 10* 20 * 5*50 * 25 + 18*1000 *(4*10+ 6 * 2 -100) = 386000
        teamId = await leagues.encodeTZCountryAndVal(tz = INIT_TZ, countryIdxInTZ = 0, teamIdxInCountry = 0)
        // the team is Still a Bot:
        result = await leagues.computeTeamRankingPoints(teamStateAll50, leagueRanking = 7, prevPerfPoints = 10, teamId).should.be.fulfilled;
        result[0].toNumber().should.be.equal(0 * TWO_TO_28 + MAX_TEAMIDX_IN_COUNTRY - teamIdxInCountry);
        result[1].toNumber().should.be.equal(0);
        // make it human:
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, accounts[0], {from: owners.relay}).should.be.fulfilled;
        result = await leagues.computeTeamRankingPoints(teamStateAll50, leagueRanking = 7, prevPerfPoints = 10, teamId).should.be.fulfilled;
        result[0].toNumber().should.be.equal(386000*TWO_TO_28 + MAX_TEAMIDX_IN_COUNTRY - teamIdxInCountry);
        // prevPerfPoints = 0.6 * 2 + 0.4 * 10 = 5.2
        result[1].toNumber().should.be.equal(5);
    });

    it('computeLeagueLeaderBoard almost no clashes', async () =>  {
        matchDay = 12;
        teamIds = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => web3.utils.toBN(MAX_TEAMIDX_IN_COUNTRY - 10 - i));
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [getRand(2*i, 0, 12), getRand(2*i+1, 0, 12)]);
        expectedPoints = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => 0);
        for (m = 0; m < ((matchDay + 1) * 4); m++) {
            teams = await leagues.getTeamsInLeagueMatch(Math.floor(m / 4), m % 4); 
            if (results[m][0] == results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 1;
                expectedPoints[teams[1].toNumber()] += 1;
            } else if (results[m][0] > results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 3;
            } else {
                expectedPoints[teams[1].toNumber()] += 3;
            }
            // the next printout shows that team 6 won against team 3, the oponent with the same ranking points
            // if ((teams[0].toNumber() == 6) || (teams[1].toNumber() == 6)) {
            //     console.log(teams[0].toNumber(), teams[1].toNumber(), results[m][0], results[m][1]);
            // }
        }
        debug.compareArrays(expectedPoints, [ 14, 24, 23, 12, 26, 21, 12, 15 ], toNum = false, isBN = false);
        expectedPointsSorted = expectedPoints.slice().sort((a,b)=>b-a);

        result = await leagues.computeLeagueLeaderBoard(teamIds, results, matchDay).should.be.fulfilled;
        expectedPoints9Digits =  [ 26000000000, 24000000000, 23000000000, 21000000000, 15000000000, 14000000000, 12001081000, 12000075000 ];
        expectedPointsTeamIdPart = [ 0, 0, 0, 0, 0, 0, 9999731564561, 9999731564558 ];
        expectedRanking = [ 4, 1, 2, 5, 7, 0, 6, 3 ];
        reportedPoints2Digits = [];
        reportedPoints9Digits = [];
        reportedPointsTeamIdPart = [];
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPoints2Digits.push(result.points[t].div(TEN_TO_13).div(web3.utils.toBN(10**9)).toNumber());
            reportedPoints9Digits.push(result.points[t].div(TEN_TO_13).toNumber());
            reportedPointsTeamIdPart.push(result.points[t].mod(TEN_TO_13).toNumber());
        }
        debug.compareArrays(reportedPoints2Digits, expectedPointsSorted, toNum = false);
        reportedPointsTeamIdPart[6].should.be.equal(TEN_TO_13.toNumber()-teamIds[expectedRanking[6]].toNumber());
        reportedPointsTeamIdPart[7].should.be.equal(TEN_TO_13.toNumber()-teamIds[expectedRanking[7]].toNumber());

        debug.compareArrays(result.ranking, expectedRanking, toNum = true, isBN = false);
        debug.compareArrays(reportedPoints9Digits, expectedPoints9Digits, toNum = false, isBN = false);
        debug.compareArrays(reportedPointsTeamIdPart, expectedPointsTeamIdPart, toNum = false, isBN = false);
    });

    it('computeLeagueLeaderBoard at start', async () =>  {
        matchDay = 0;
        teamIds = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => web3.utils.toBN(MAX_TEAMIDX_IN_COUNTRY - 10 - i));
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [0,0]);
        expectedPoints = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => 0);
        for (m = 0; m < ((matchDay + 1) * 4); m++) {
            teams = await leagues.getTeamsInLeagueMatch(Math.floor(m / 4), m % 4); 
            if (results[m][0] == results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 1;
                expectedPoints[teams[1].toNumber()] += 1;
            } else if (results[m][0] > results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 3;
            } else {
                expectedPoints[teams[1].toNumber()] += 3;
            }
            // the next printout shows that team 6 won against team 3, the oponent with the same ranking points
            // if ((teams[0].toNumber() == 6) || (teams[1].toNumber() == 6)) {
            //     console.log(teams[0].toNumber(), teams[1].toNumber(), results[m][0], results[m][1]);
            // }
        }
        debug.compareArrays(expectedPoints, [ 1, 1, 1, 1, 1, 1, 1, 1 ], toNum = false, isBN = false);
        expectedPointsSorted = expectedPoints.slice().sort((a,b)=>b-a);

        result = await leagues.computeLeagueLeaderBoard(teamIds, results, matchDay).should.be.fulfilled;
        expectedPoints9Digits =  [ 1000000000, 1000000000, 1000000000, 1000000000, 1000000000, 1000000000, 1000000000, 1000000000 ];
        expectedPointsTeamIdPart = [ 9999731564562, 9999731564561, 9999731564560, 9999731564559, 9999731564558, 9999731564557, 9999731564556, 9999731564555 ];
        expectedRanking = [ 7, 6, 5, 4, 3, 2, 1, 0 ];
        reportedPoints2Digits = [];
        reportedPoints9Digits = [];
        reportedPointsTeamIdPart = [];
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPoints2Digits.push(result.points[t].div(TEN_TO_13).div(web3.utils.toBN(10**9)).toNumber());
            reportedPoints9Digits.push(result.points[t].div(TEN_TO_13).toNumber());
            reportedPointsTeamIdPart.push(result.points[t].mod(TEN_TO_13).toNumber());
        }
        debug.compareArrays(reportedPoints2Digits, expectedPointsSorted, toNum = false);
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, isBN = false);
        debug.compareArrays(reportedPoints9Digits, expectedPoints9Digits, toNum = false, isBN = false);
        debug.compareArrays(reportedPointsTeamIdPart, expectedPointsTeamIdPart, toNum = false, isBN = false);
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPointsTeamIdPart[t].should.be.equal(TEN_TO_13.toNumber()-teamIds[expectedRanking[t]].toNumber());
        }
    });

    it('computeLeagueLeaderBoard at end of league', async () =>  {
        matchDay = 13;
        teamIds = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => web3.utils.toBN(MAX_TEAMIDX_IN_COUNTRY - 10 - i));
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [5,5]);
        expectedPoints = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => 0);
        for (m = 0; m < ((matchDay + 1) * 4); m++) {
            teams = await leagues.getTeamsInLeagueMatch(Math.floor(m / 4), m % 4); 
            if (results[m][0] == results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 1;
                expectedPoints[teams[1].toNumber()] += 1;
            } else if (results[m][0] > results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 3;
            } else {
                expectedPoints[teams[1].toNumber()] += 3;
            }
            // the next printout shows that team 6 won against team 3, the oponent with the same ranking points
            // if ((teams[0].toNumber() == 6) || (teams[1].toNumber() == 6)) {
            //     console.log(teams[0].toNumber(), teams[1].toNumber(), results[m][0], results[m][1]);
            // }
        }
        debug.compareArrays(expectedPoints, [ 14, 14, 14, 14, 14, 14, 14, 14 ], toNum = false, isBN = false);
        expectedPointsSorted = expectedPoints.slice().sort((a,b)=>b-a);

        result = await leagues.computeLeagueLeaderBoard(teamIds, results, matchDay).should.be.fulfilled;
        expectedPoints9Digits =  [ 14000070000, 14000070000, 14000070000, 14000070000, 14000070000, 14000070000, 14000070000, 14000070000 ];
        expectedPointsTeamIdPart = [ 9999731564562, 9999731564561, 9999731564560, 9999731564559, 9999731564558, 9999731564557, 9999731564556, 9999731564555 ];
        expectedRanking = [ 7, 6, 5, 4, 3, 2, 1, 0 ];
        reportedPoints2Digits = [];
        reportedPoints9Digits = [];
        reportedPointsTeamIdPart = [];
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPoints2Digits.push(result.points[t].div(TEN_TO_13).div(web3.utils.toBN(10**9)).toNumber());
            reportedPoints9Digits.push(result.points[t].div(TEN_TO_13).toNumber());
            reportedPointsTeamIdPart.push(result.points[t].mod(TEN_TO_13).toNumber());
        }
        debug.compareArrays(reportedPoints2Digits, expectedPointsSorted, toNum = false);
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, isBN = false);
        debug.compareArrays(reportedPoints9Digits, expectedPoints9Digits, toNum = false, isBN = false);
        debug.compareArrays(reportedPointsTeamIdPart, expectedPointsTeamIdPart, toNum = false, isBN = false);
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPointsTeamIdPart[t].should.be.equal(TEN_TO_13.toNumber()-teamIds[expectedRanking[t]].toNumber());
        }
    });

    
    it('computeLeagueLeaderBoard many clashes', async () =>  {
        matchDay = 12;
        teamIds = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => web3.utils.toBN(MAX_TEAMIDX_IN_COUNTRY - 10 - i));
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [getRand(2*i+1, 0, 2), getRand(2*i+3, 0, 12)]);
        expectedPoints = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => 0);
        for (m = 0; m < ((matchDay + 1) * 4); m++) {
            teams = await leagues.getTeamsInLeagueMatch(Math.floor(m / 4), m % 4); 
            if (results[m][0] == results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 1;
                expectedPoints[teams[1].toNumber()] += 1;
            } else if (results[m][0] > results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 3;
            } else {
                expectedPoints[teams[1].toNumber()] += 3;
            }
            // the next printout shows that team 6 won against team 3, the oponent with the same ranking points
            // if ((teams[0].toNumber() == 6) || (teams[1].toNumber() == 6)) {
            //     console.log(teams[0].toNumber(), teams[1].toNumber(), results[m][0], results[m][1]);
            // }
        }
        debug.compareArrays(expectedPoints, [ 16, 22, 22, 16, 16, 18, 15, 27 ], toNum = false, isBN = false);
        expectedPointsSorted = expectedPoints.slice().sort((a,b)=>b-a);

        result = await leagues.computeLeagueLeaderBoard(teamIds, results, matchDay).should.be.fulfilled;
        expectedPoints9Digits =  [ 27000000000, 22000052000, 22000051000, 18000000000, 16002047000, 16001047000, 16000049000, 15000000000 ];
        expectedPointsTeamIdPart = [ 0, 9999731564557, 9999731564556, 0, 9999731564559, 9999731564558, 9999731564555, 0 ];
        expectedRanking = [ 7, 2, 1, 5, 4, 3, 0, 6 ];
        reportedPoints2Digits = [];
        reportedPoints9Digits = [];
        reportedPointsTeamIdPart = [];
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPoints2Digits.push(result.points[t].div(TEN_TO_13).div(web3.utils.toBN(10**9)).toNumber());
            reportedPoints9Digits.push(result.points[t].div(TEN_TO_13).toNumber());
            reportedPointsTeamIdPart.push(result.points[t].mod(TEN_TO_13).toNumber());
        }
        debug.compareArrays(reportedPoints2Digits, expectedPointsSorted, toNum = false);
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, isBN = false);
        debug.compareArrays(reportedPoints9Digits, expectedPoints9Digits, toNum = false, isBN = false);
        debug.compareArrays(reportedPointsTeamIdPart, expectedPointsTeamIdPart, toNum = false, isBN = false);
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            if (expectedPointsTeamIdPart[t] != 0) {
                reportedPointsTeamIdPart[t].should.be.equal(TEN_TO_13.toNumber()-teamIds[expectedRanking[t]].toNumber());
            }
        }
    });


    it('computeLeagueLeaderBoard all clashes', async () =>  {
        // all results in this league are [0,0]
        // since all results are [0, 0], there is no extra points due to goals, nor due to winning a team you played with
        // so the only breaking number is teamId
        matchDay = 12;
        teamIds = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => web3.utils.toBN(MAX_TEAMIDX_IN_COUNTRY - 10 - i));
        results = Array.from(new Array(MATCHES_PER_LEAGUE), (x,i) => [getRand(2*i+1, 0, 1), getRand(2*i+3, 0, 1)]);
        expectedPoints = Array.from(new Array(TEAMS_PER_LEAGUE.toNumber()), (x,i) => 0);
        for (m = 0; m < ((matchDay + 1) * 4); m++) {
            teams = await leagues.getTeamsInLeagueMatch(Math.floor(m / 4), m % 4); 
            if (results[m][0] == results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 1;
                expectedPoints[teams[1].toNumber()] += 1;
            } else if (results[m][0] > results[m][1]) {
                expectedPoints[teams[0].toNumber()] += 3;
            } else {
                expectedPoints[teams[1].toNumber()] += 3;
            }
            // the next printout shows that team 6 won against team 3, the oponent with the same ranking points
            // if ((teams[0].toNumber() == 6) || (teams[1].toNumber() == 6)) {
                // console.log(teams[0].toNumber(), teams[1].toNumber(), results[m][0], results[m][1]);
            // }
        }
        debug.compareArrays(expectedPoints, [ 13, 13, 13, 13, 13, 13, 13, 13 ], toNum = false, isBN = false);
        expectedPointsSorted = expectedPoints.slice().sort((a,b)=>b-a);

        result = await leagues.computeLeagueLeaderBoard(teamIds, results, matchDay).should.be.fulfilled;
        expectedPoints9Digits =  [ 13000000000, 13000000000, 13000000000, 13000000000, 13000000000, 13000000000, 13000000000, 13000000000 ];
        expectedPointsTeamIdPart = [ 9999731564562, 9999731564561, 9999731564560, 9999731564559, 9999731564558, 9999731564557, 9999731564556, 9999731564555 ];
        expectedRanking = [ 7, 6, 5, 4, 3, 2, 1, 0 ];
        reportedPoints2Digits = [];
        reportedPoints9Digits = [];
        reportedPointsTeamIdPart = [];
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            reportedPoints2Digits.push(result.points[t].div(TEN_TO_13).div(web3.utils.toBN(10**9)).toNumber());
            reportedPoints9Digits.push(result.points[t].div(TEN_TO_13).toNumber());
            reportedPointsTeamIdPart.push(result.points[t].mod(TEN_TO_13).toNumber());
        }
        debug.compareArrays(reportedPoints2Digits, expectedPointsSorted, toNum = false);
        debug.compareArrays(result.ranking, expectedRanking, toNum = true, isBN = false);
        debug.compareArrays(reportedPoints9Digits, expectedPoints9Digits, toNum = false, isBN = false);
        debug.compareArrays(reportedPointsTeamIdPart, expectedPointsTeamIdPart, toNum = false, isBN = false);
        for (t = 0; t < TEAMS_PER_LEAGUE.toNumber(); t++) {
            if (expectedPointsTeamIdPart[t] != 0) {
                reportedPointsTeamIdPart[t].should.be.equal(TEN_TO_13.toNumber()-teamIds[expectedRanking[t]].toNumber());
            }
        }
    });

    it('check initial constants', async () =>  {
        MATCHDAYS.toNumber().should.be.equal(14);
        MATCHES_PER_DAY.toNumber().should.be.equal(4);
        TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    });

    it('getTeamsInCupPlayoffMatch', async () => {
        teamsExpected = [0,7,9,14,4,11,13,18,8,15,17,22,12,19,21,26,16,23,25,30,20,27,29,34,24,31,33,38,28,35,37,42,32,39,41,46,36,43,45,50,40,47,49,54,44,51,53,58,48,55,57,62,52,59,61,2,56,63,1,6,60,3,5,10];
        for (t = 0; t < 32; t++) {
            team = await leagues.getTeamsInCupPlayoffMatch(matchIdxInDay = t).should.be.fulfilled;
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
            team = await leagues.getTeamIdxInCup(groupIdx = 0, posInGroup = t).should.be.fulfilled;
            team.toNumber().should.be.equal(teamsExpected[t]);
            result = await leagues.getGroupAndPosInGroup(team.toNumber()).should.be.fulfilled;
            result[0].toNumber().should.be.equal(groupIdx);
            result[1].toNumber().should.be.equal(posInGroup);
        }
        teamsExpected = [71, 79, 87, 95, 103, 111, 119, 127 ]
        for (t = 0; t < teamsExpected.length; t++) {
            team = await leagues.getTeamIdxInCup(groupIdx = 15, posInGroup = t).should.be.fulfilled;
            team.toNumber().should.be.equal(teamsExpected[t])
            result = await leagues.getGroupAndPosInGroup(team.toNumber()).should.be.fulfilled;
            result[0].toNumber().should.be.equal(groupIdx);
            result[1].toNumber().should.be.equal(posInGroup);
        }
    });

    it('get all teams for particular matches', async () => {
        teams = await leagues.getTeamsInCupLeagueMatch(groupIdx = 0, day = 0, matchIdxInDay = 0).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(8);
        teams = await leagues.getTeamsInCupLeagueMatch(groupIdx = 0, day = day = Math.floor(MATCHDAYS/2), matchIdxInDay = 0).should.be.rejected;
        teams = await leagues.getTeamsInCupLeagueMatch(groupIdx = 15, day = 0, matchIdxInDay = 0).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(71);
        teams[1].toNumber().should.be.equal(79);
    });

    it('get teams for match in wrong day', async () => {
        await leagues.getTeamsInLeagueMatch(day = MATCHDAYS-1, matchIdxInDay = 0).should.be.fulfilled;
        await leagues.getTeamsInLeagueMatch(day = MATCHDAYS, matchIdxInDay = 0).should.be.rejected;
    });

    it('get teams for match in wrong match in day', async () => {
        await leagues.getTeamsInLeagueMatch(day = 0, matchIdxInDay = MATCHES_PER_DAY-1).should.be.fulfilled;
        await leagues.getTeamsInLeagueMatch(day = 0, matchIdxInDay = MATCHES_PER_DAY).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        teams = await leagues.getTeamsInLeagueMatch(day = 0, matchIdxInDay = 0).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        teams = await leagues.getTeamsInLeagueMatch(day = Math.floor(MATCHDAYS/2), matchIdxInDay).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });
    

});