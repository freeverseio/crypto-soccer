const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const fs = require('fs');
const truffleAssert = require('truffle-assertions');
const logUtils = require('../utils/matchLogUtils.js');
const debug = require('../utils/debugUtils.js');

const Utils = artifacts.require('Utils');
const TrainingPoints = artifacts.require('TrainingPoints');
const Evolution = artifacts.require('Evolution');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');
const EngineApplyBoosters = artifacts.require('EngineApplyBoosters');
const PlayAndEvolve = artifacts.require('PlayAndEvolve');
const Shop = artifacts.require('Shop');
const Championships = artifacts.require('Championships');


contract('FullLeague', (accounts) => {
    const JUST_CHECK_AGAINST_EXPECTED_RESULTS = 0;
    const WRITE_NEW_EXPECTED_RESULTS = 1;
    const nLeafs = 1024;
    const nMatchdays = 14;
    const nMatchesPerDay = 4;
    const nTeamsInLeague = 8;
    const nMatchesPerLeague = nMatchesPerDay * nMatchdays;
    const nPlayersInTeam = 25;
    const substitutions = [6, 10, 0];
    const subsRounds = [3, 7, 1];
    const noSubstitutions = [11, 11, 11];
    const lineup0 = [0, 3, 4, 5, 6, 9, 10, 11, 12, 15, 16, 7, 13, 17];
    const lineup1 = [0, 3, 4, 5, 6, 9, 10, 11, 16, 17, 18, 7, 13, 17];
    const lineupConsecutive =  Array.from(new Array(14), (x,i) => i);
    const extraAttackNull =  Array.from(new Array(10), (x,i) => 0);
    const tacticId442 = 0; // 442
    const tacticId433 = 2; // 433
    const playersPerZone442 = [1,2,1,1,2,1,0,2,0];
    const playersPerZone433 = [1,2,1,1,1,1,1,1,1];
    const PLAYERS_PER_TEAM_MAX = 25;
    const firstHalfLog = [0, 0];
    const subLastHalf = false;
    const is2ndHalf = false;
    const isHomeStadium = true;
    const isPlayoff = false;
    const matchBools = [is2ndHalf, isHomeStadium, isPlayoff]
    const IDX_R = 1;
    const IDX_C = 2;
    const IDX_CR = 3;
    const IDX_L = 4;
    const IDX_LR = 5;
    const IDX_LC = 6;
    const IDX_LCR = 7;
    const fwd442 =  [0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3];
    const left442 = [0, IDX_L, IDX_C, IDX_C, IDX_R, IDX_L, IDX_C, IDX_C, IDX_R, IDX_C, IDX_C];
    // const now = Math.floor(new Date()/1000);
    // const dayOfBirth21 = Math.round(secsToDays(now) - 21/7);
    const now = 1570147200; // this number has the property that 7*nowFake % (SECS_IN_DAY) = 0 and it is basically Oct 3, 2019
    const dayOfBirth21 = secsToDays(now) - 21*365/7; // = exactly 17078, no need to round
    const MAX_PENALTY = 10000;
    const MAX_GOALS = 12;
    const RED_CARD = 3;

    const assistersIdx = Array.from(new Array(MAX_GOALS), (x,i) => i);
    const shootersIdx  = Array.from(new Array(MAX_GOALS), (x,i) => 1);
    const shooterForwardPos  = Array.from(new Array(MAX_GOALS), (x,i) => 1);
    const penalties  = Array.from(new Array(7), (x,i) => false);
    const typesOutOfGames = [0, 0];
    const outOfGameRounds = [0, 0];
    const yellowCardedDidNotFinish1stHalf = [false, false];
    const ingameSubs1 = [0, 0, 0]
    const ingameSubs2 = [0, 0, 0]
    const outOfGames = [14, 14]
    const yellowCards1 = [14, 14]
    const yellowCards2 = [14, 14]
    const halfTimeSubstitutions = [14, 14, 14]
    const nDefs1 = 4; 
    const nDefs2 = 4; 
    const nTot = 11; 
    const winner = 2; // DRAW = 2
    const isHomeSt = false;
    const teamSumSkillsDefault = 0;
    const trainingPointsInit = 0;
    
    // Skills: shoot, speed, pass, defence, endurance
    const SK_SHO = 0;
    const SK_SPE = 1;
    const SK_PAS = 2;
    const SK_DEF = 3;
    const SK_END = 4;
    
    const it2 = async(text, f) => {};

    function setNoSubstInLineUp(lineup, substitutions) {
        modifiedLineup = [...lineup];
        NO_SUBST = 11;
        NO_LINEUP = 25;
        for (s = 0; s < 3; s++) {
            if (substitutions[s] == NO_SUBST) modifiedLineup[s + 11] = NO_LINEUP;
        }
        return modifiedLineup;
    }

    function daysToSecs(dayz) {
        return (dayz * 24 * 3600); 
    }

    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }
    
    function getDefaultTPs() {
        TP = 200;
        TPperSkill = Array.from(new Array(25), (x,i) => TP/5 - 3*i % 6);
        specialPlayer = 21;
        // make sure they sum to TP:
        for (bucket = 0; bucket < 5; bucket++){
            sum4 = 0;
            for (sk = 5 * bucket; sk < (5 * bucket + 4); sk++) {
                sum4 += TPperSkill[sk];
            }
            TPperSkill[5 * bucket + 4] = TP - sum4;
        }       
        return [TP, TPperSkill];
    }
    
    const createTeamState442 = async (engine, forceSkills, alignedEndOfLastHalfTwoVec = [false, false]) => {
        teamState = [];
        playerId = 123121;
        pot = 3;
        aggr = 0;
        alignedEndOfLastHalf = true;
        redCardLastGame = false;
        gamesNonStopping = 0;
        injuryWeeksLeft = 0;
        sumSkills = forceSkills.reduce((a, b) => a + b, 0);
        for (p = 0; p < 11; p++) {
            pSkills = await assets.encodePlayerSkills(forceSkills, dayOfBirth21, gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedEndOfLastHalfTwoVec[0], redCardLastGame, gamesNonStopping, 
                injuryWeeksLeft, subLastHalf, sumSkills).should.be.fulfilled 
            teamState.push(pSkills)
        }
        p = 10;
        pSkills = await assets.encodePlayerSkills(forceSkills, dayOfBirth21, gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedEndOfLastHalfTwoVec[1], redCardLastGame, gamesNonStopping, 
                injuryWeeksLeft, subLastHalf, sumSkills).should.be.fulfilled 
        for (p = 11; p < PLAYERS_PER_TEAM_MAX; p++) {
            teamState.push(pSkills)
        }        
        return teamState;
    };

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

    const createHardcodedTeam = function () {
        // returns 18 players generated with the following code. We hardcode it to avoid the "deployDate" time-dependency
        // teamState = [];
        // playerId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdx = 0, playerIdx = 0).should.be.fulfilled;
        // for (p = 0; p < 18; p++) {
        //     skills = await assets.getPlayerSkillsAtBirth(playerId0.toNumber() + p);
        //     teamState.push(skills);
        //     console.log(skills.toString(10))
        // }
        return [
            '14606248079918261338806855150670198598294524424421999',
            '14603325075249802958062362651785117246719383552393656',
            '14615017086954653606499907426763036762091679724733245',
            '14609171184243174825485386589332947715467405749846827',
            '14615017461189033969342085869889674545308663693968083',
            '14603325891317697566792669908219362044711638355411673',
            '14606249873734453245614329076439313941148075272765994',
            '14603324461979309998470701478621001103697221903123183',
            '14606248281321866413037179508268863783570851530343215',
            '14606249082057998697777445123967984023640370982880706',
            '14603327085801362263089568768708477093108613577769640',
            '14612095382001501327618929648053879079031002742916002',
            '14603326117112742701915784319947485139466656825672861',
            '14612093787498219632679532865607761507997232766977103',
            '14609173081200313275497388848716119026424650418029241',
            '14603326360330245023390630956127251848106222989410926',
            '14606249807529115937477333996086265720951632055960118',
            '14603326808435843856365497638008216685947959514366883'
        ];
    };
    
    function parseLog(tr) {
        arr = [
            tr.goalkeepersShoot,
            tr.goalkeepersSpeed,
            tr.goalkeepersPass,
            tr.goalkeepersDefence,
            tr.goalkeepersEndurance,
            // 
            tr.defendersShoot,
            tr.defendersSpeed,
            tr.defendersPass,
            tr.defendersDefence,
            tr.defendersEndurance,
            // 
            tr.midfieldersShoot,
            tr.midfieldersSpeed,
            tr.midfieldersPass,
            tr.midfieldersDefence,
            tr.midfieldersEndurance,
            // 
            tr.attackersShoot,
            tr.attackersSpeed,
            tr.attackersPass,
            tr.attackersDefence,
            tr.attackersEndurance,
            // 
            tr.specialPlayerShoot,
            tr.specialPlayerSpeed,
            tr.specialPlayerPass,
            tr.specialPlayerDefence,
            tr.specialPlayerEndurance,
        ];    
        for (i = 0; i < arr.length; i++) arr[i] = parseInt(arr[i]);        
        return arr;
    }
    
    function checkTPAssigment(TP, TPassigned25, verbose) {
        OK = true;
        if (verbose) console.log("Total Available: ", TP);
        for (bucket = 0; bucket < 5; bucket++) {
            sum = 0;
            for (i = bucket * 5; i < (bucket+1) * 5; i++) {
                sum += TPassigned25[i];
                thisOK = (10 * TPassigned25[i] <= 6 * TP);
                if (verbose && !thisOK) console.log("skill ", i, " exceeds 60 percent of TPs. TP_thisSkill / Available = ", TPassigned25[i]/TP);
                OK = OK && thisOK;
            }
            thisOK = (sum <= TP);
            if (verbose && !thisOK) console.log("bucket ", bucket, " exceeds available TPs. Sum / Available = ", sum/TP);
            OK = OK && thisOK;
        }        
        if (verbose) console.log("OK = ", OK);
        return OK;
    }
    
    function assertStr(cond, x, y, msg) {
        if (cond == "eq") assert(x.toString() == y.toString(), msg);
        else assert(!(x.toString() == y.toString()), msg);
    }

    function assertBN(cond, x, y, msg) {
        if (cond == "eq") assert(web3.utils.toBN(x).eq(web3.utils.toBN(y)), msg);
        else assert(!web3.utils.toBN(x).eq(web3.utils.toBN(y)), msg);
    }

    function zeroPadToLength(x, desiredLength) {
        return x.concat(Array.from(new Array(desiredLength - x.length), (x,i) => 0))
    }
    
    // var leagueData = {
    //     seeds: [], // [2 * nMatchDays]
    //     teamIds: [], // [nTeamsInLeague]
    //     startTimes: [], // [2 * nMatchDays]
    //     teamStates: [], // [2 * nMatchdays + 1][nTeamsInLeague][PLAYERS_PER_TEAM_MAX]
    //     matchLogs: [], // [2 * nMatchdays+ 1][nTeamsInLeague]
    //     results: [], // [nMatchesPerLeague][2]  -> goals per team per match
    //     points: [], // [2 * nMatchdays][nTeamsInLeague]
    //     tactics: [], // [2 * nMatchdays + 1][nTeamsInLeague]
    //     trainings: [] // [2 * nMatchdays + 1][nTeamsInLeague]
    // }
    // - Data[1024] = [League[512], Team$_{i, aft}$[32], Team$_{i, bef}$[32]]
    // League[128] = leafsLeague[128] = [Points[team=0,..,7], Goals[56][2]]
    // 
    // returns leafs AFTER having played the matches at matchday = day, half = half.
    //  - sorting results:
    //      - idx = day * nMatchesPerDay * 2 + matchInDay * 2 + teamHomeOrAway
    function buildLeafs(leagueData, day, half) {
        var isNoPointsYet = (half == 0) && (day == 0);
        if (isNoPointsYet) { 
            leafs =  Array.from(new Array(nTeamsInLeague), (x,i) => 0);
        } else {
            lastDayToCount = (half == 0) ? day - 1 : day;
            leafs = leagueData.points[lastDayToCount]; 
            for (d = 0; d < lastDayToCount; d++) {
                leafs.push(leagueData.results[d][0]);
                leafs.push(leagueData.results[d][1]);
            }
        }
        leafs = zeroPadToLength(leafs, 128);
        for (team = 0; team < nTeamsInLeague; team++) {
            for (extraHalf = 0; extraHalf < 2; extraHalf++) {
                teamData = [];
                for (p = 0; p < nPlayersInTeam; p++) {
                    teamData.push(leagueData.teamStates[2*day + half + extraHalf][team][p])
                }
                teamData.push(leagueData.tactics[2*day + half + extraHalf][team]);
                teamData.push(leagueData.trainings[2*day + half + extraHalf][team]);
                teamData.push(leagueData.matchLogs[2*day + half + extraHalf][team]);
                leafs = leafs.concat(zeroPadToLength(teamData, 32));
            }
        }
        return zeroPadToLength(leafs, nLeafs);
    }

    function vec2str(y) {
        yStr = [...y];
        for (i = 0; i < y.length; i++) yStr[i] = y[i].toString();
        return yStr;
    }
    
    beforeEach(async () => {
        training = await TrainingPoints.new().should.be.fulfilled;
        evo = await Evolution.new().should.be.fulfilled;
        play = await PlayAndEvolve.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        market = await Market.new().should.be.fulfilled;
        shop = await Shop.new().should.be.fulfilled;
        encodeLog = await EncodingMatchLog.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        applyBoosters = await EngineApplyBoosters.new().should.be.fulfilled;
        await engine.setPreCompAddr(precomp.address).should.be.fulfilled;
        await engine.setApplyBoostersAddr(applyBoosters.address).should.be.fulfilled;
        await training.setAssetsAddress(assets.address).should.be.fulfilled;
        await training.setMarketAddress(market.address).should.be.fulfilled;
        await play.setEngineAddress(engine.address).should.be.fulfilled;
        await play.setTrainingAddress(training.address).should.be.fulfilled;
        await play.setEvolutionAddress(evo.address).should.be.fulfilled;
        await play.setShopAddress(shop.address).should.be.fulfilled;
        
        tactics0 = await engine.encodeTactics(substitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, substitutions), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        tactics1 = await engine.encodeTactics(substitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, substitutions), 
            extraAttackNull, tacticId433).should.be.fulfilled;
        tactics1NoChanges = await engine.encodeTactics(noSubstitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, noSubstitutions), 
            extraAttackNull, tacticId433).should.be.fulfilled;
        tactics442 = await engine.encodeTactics(substitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, substitutions),
            extraAttackNull, tacticId442).should.be.fulfilled;
        tactics442NoChanges = await engine.encodeTactics(noSubstitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, noSubstitutions), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        teamStateAll50Half1 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2, aligned = [false, false]).should.be.fulfilled;
        teamStateAll1Half1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2, aligned = [false, false]).should.be.fulfilled;
        teamStateAll50Half2 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2, aligned = [true, false]).should.be.fulfilled;
        teamStateAll1Half2 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2, aligned = [true, false]).should.be.fulfilled;
        MAX_RND = await engine.MAX_RND().should.be.fulfilled;
        MAX_RND = MAX_RND.toNumber();
        kMaxRndNumHalf = Math.floor(MAX_RND/2)-200; 
        events1Half = Array.from(new Array(7), (x,i) => 0);
        events1Half = [events1Half,events1Half];

        TPperSkill =  Array.from(new Array(25), (x,i) => 0);
        almostNullTraning = await training.encodeTP(TP = 0, TPperSkill, specialPlayer = 21).should.be.fulfilled;

    });
  
    // leafsLeague[128] = [Points[team=0,..,7], ML[team = 0,1; matchInDay = 0,1,2,3; matchDay = 0,..13], 0,...]
    it2('create real data for an entire league', async () => {
        mode = WRITE_NEW_EXPECTED_RESULTS; // JUST_CHECK_AGAINST_EXPECTED_RESULTS for testing, 1 WRITE_NEW_EXPECTED_RESULTS
        // prepare a training that is not identical to the bignumber(0), but which works irrespective of the previously earned TP
        // => all assingments to 0, but with a special player chosen

        champs = await Championships.new().should.be.fulfilled;
        let secsBetweenMatches = 12*3600;
        var leagueData = {
            seeds: [], // [2 * nMatchDays]
            teamIds: [], // [nTeamsInLeague]
            startTimes: [], // [2 * nMatchDays]
            teamStates: [], // [1 + 2 * nMatchdays][nTeamsInLeague][PLAYERS_PER_TEAM_MAX]
            matchLogs: [], // [1 + 2 * nMatchdays][nTeamsInLeague]
            results: [], // [nMatchesPerLeague][2]  ->  per team per match
            points: [], // [2 * nMatchdays][nTeamsInLeague]
            tactics: [], // [2 * nMatchdays + 1][nTeamsInLeague]
            trainings: [] // [2 * nMatchdays + 1][nTeamsInLeague]
        }
        // on starting points: if we query computeLeagueLeaderBoard, I would get 
        // a non-null value, sorting because of all tied, which would depend on a seed.
        // we don't have that seed before a match starts, so we set all points to 0.

        leagueData.seeds = Array.from(new Array(2 * nMatchdays), (x,i) => web3.utils.keccak256(i.toString()).toString());
        leagueData.startTimes = Array.from(new Array(2 * nMatchdays), (x,i) => now + i * secsBetweenMatches);
        allMatchLogs = Array.from(new Array(nTeamsInLeague), (x,i) => 0);
        leagueData.matchLogs.push([...allMatchLogs]);
        teamState442 = await createTeamState442(engine, forceSkills= [1000,1000,1000,1000,1000]).should.be.fulfilled;
        teamState442 = vec2str(teamState442);
        allTeamsSkills = Array.from(new Array(nTeamsInLeague), (x,i) => teamState442);
        leagueData.teamStates.push([...allTeamsSkills]);
        // nosub = [NO_SUBST, NO_SUBST, NO_SUBST];
        // tact = await engine.encodeTactics(nosub , ro = [0, 0, 0], setNoSubstInLineUp(lineupConsecutive, nosub), extraAttackNull, tacticsId = 0).should.be.fulfilled;
        teamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry = 0);
        leagueData.teamIds = Array.from(new Array(nTeamsInLeague), (x,i) => teamId.toNumber() + i);
        leagueData.results = Array.from(new Array(nMatchesPerLeague), (x,i) => [0,0]);

        // tactics and trainings start at all 0 (undefined until we play the first match)
        leagueData.tactics.push(Array.from(new Array(nTeamsInLeague), (x,i) => 0));
        leagueData.trainings.push(Array.from(new Array(nTeamsInLeague), (x,i) => 0));
        // same tactics and trainings for all matchdays:
        tact = tactics442NoChanges.toString();
        for (day = 0; day < 2 * nMatchdays; day++) {
            leagueData.tactics.push(Array.from(new Array(nTeamsInLeague), (x,i) => tact));
            leagueData.trainings.push(Array.from(new Array(nTeamsInLeague), (x,i) => almostNullTraning.toString()));
        }

        // we just need to build, across the league: teamStates, points, teamIds
        // for (day = 0; day < 1; day++) {
        for (day = 0; day < nMatchdays; day++) {
            // 1st half
            for (matchIdxInDay = 0; matchIdxInDay < nMatchesPerDay; matchIdxInDay++) {
                var {0: t0, 1: t1} = await champs.getTeamsInLeagueMatch(day, matchIdxInDay).should.be.fulfilled;
                t0 = t0.toNumber();
                t1 = t1.toNumber();
                console.log("day:", day, ", matchIdxInDay:", matchIdxInDay, ", half 0,  teams:", t0, t1);
                var {0: newSkills, 1: newLogs} =  await play.play1stHalfAndEvolve(
                    leagueData.seeds[2 * day], leagueData.startTimes[2 * day], 
                    [allTeamsSkills[t0], allTeamsSkills[t1]], 
                    [leagueData.teamIds[t0], leagueData.teamIds[t1]], 
                    [leagueData.tactics[2 * day + 1][t0], leagueData.tactics[2 * day + 1][t1]], 
                    [allMatchLogs[t0], allMatchLogs[t1]],
                    [is2nd = false, isHom = true, isPlay = false],
                    [tp = 0, tp = 0]
                ).should.be.fulfilled;
                allTeamsSkills[t0] = vec2str(newSkills[0]);
                allTeamsSkills[t1] = vec2str(newSkills[1]);
                allMatchLogs[t0] = newLogs[0].toString();
                allMatchLogs[t1] = newLogs[1].toString();
            }
            leagueData.teamStates.push([...allTeamsSkills]);        
            leagueData.matchLogs.push([...allMatchLogs]);        
            // 2nd half
            for (matchIdxInDay = 0; matchIdxInDay < nMatchesPerDay; matchIdxInDay++) {
                var {0: t0, 1: t1} = await champs.getTeamsInLeagueMatch(day, matchIdxInDay).should.be.fulfilled;
                t0 = t0.toNumber();
                t1 = t1.toNumber();
                console.log("day:", day, ", matchIdxInDay:", matchIdxInDay, ", half 1,  teams:", t0, t1);
                var {0: newSkills, 1: newLogs} =  await play.play2ndHalfAndEvolve(
                    leagueData.seeds[2*day + 1], leagueData.startTimes[2*day + 1], 
                    [allTeamsSkills[t0], allTeamsSkills[t1]], 
                    [leagueData.teamIds[t0], leagueData.teamIds[t1]], 
                    [leagueData.tactics[2 * day + 1][t0], leagueData.tactics[2 * day + 1][t1]], 
                    [allMatchLogs[t0], allMatchLogs[t1]],
                    [is2nd = true, isHom = true, isPlay = false]
                ).should.be.fulfilled;
                allTeamsSkills[t0] = vec2str(newSkills[0]);
                allTeamsSkills[t1] = vec2str(newSkills[1]);
                allMatchLogs[t0] = newLogs[0].toString();
                allMatchLogs[t1] = newLogs[1].toString(); 
                goals0 = await encodeLog.getNGoals(newLogs[0]).should.be.fulfilled;
                goals1 = await encodeLog.getNGoals(newLogs[1]).should.be.fulfilled;
                leagueData.results[nMatchesPerDay * day + matchIdxInDay] = [goals0.toNumber(), goals1.toNumber()];
            }
            leagueData.teamStates.push([...allTeamsSkills]);        
            leagueData.matchLogs.push([...allMatchLogs]);   
            var {0: rnking, 1: lPoints} = await champs.computeLeagueLeaderBoard([...leagueData.results], day, leagueData.seeds[2*day + 1]).should.be.fulfilled;
            leagueData.points.push(vec2str(lPoints));   
        }
        if (mode == WRITE_NEW_EXPECTED_RESULTS) {
            fs.writeFileSync('test/testdata/fullleague.json', JSON.stringify(leagueData), function(err) {
                if (err) {
                    console.log(err);
                }
            });
        }
        expectedData = fs.readFileSync('test/testdata/fullleague.json', 'utf8');
        assert.equal(
            web3.utils.keccak256(expectedData),
            web3.utils.keccak256(JSON.stringify(leagueData)),
            "leafs do not coincide with expected"
        );
    });

    it2('read an entire league and organize data in the leaf format required', async () => {
        mode = WRITE_NEW_EXPECTED_RESULTS; // JUST_CHECK_AGAINST_EXPECTED_RESULTS for testing, 1 WRITE_NEW_EXPECTED_RESULTS
        leagueData = JSON.parse(fs.readFileSync('test/testdata/fullleague.json', 'utf8'));
        var leafs = [];
        for (day = 0; day < nMatchdays; day++) {
            dayLeafs = buildLeafs(leagueData, day, half = 0);
            leafs.push([...dayLeafs]);
            dayLeafs = buildLeafs(leagueData, day, half = 1);
            leafs.push([...dayLeafs]);
        }
        if (mode == WRITE_NEW_EXPECTED_RESULTS) {
            fs.writeFileSync('test/testdata/leafsPerHalf.json', JSON.stringify(leafs), function(err) {
                if (err) {
                    console.log(err);
                }
            });
        }
        expectedLeafs = fs.readFileSync('test/testdata/leafsPerHalf.json', 'utf8');
        assert.equal(
            web3.utils.keccak256(expectedLeafs),
            web3.utils.keccak256(JSON.stringify(leafs)),
            "leafs do not coincide with expected"
        );
    });
    
    it2('test day 0, half 0', async () => {
        leafs = JSON.parse(fs.readFileSync('test/testdata/leafsPerHalf.json', 'utf8'));
        assert.equal(leafs.length, nMatchdays * 2);
        assert.equal(leafs[0].length, nLeafs);
        // at end of 1st half we still do not have end-game results
        for (i = 0; i < 128; i++) {
            console.log(leafs[0][i]);
            assertBN('eq', leafs[0][i], 0, "unexpected non-null leaf at start of league");
        }
        for (team = 0; team < nTeamsInLeague; team++) {
            // BEFORE first half ---------
            off = 128 + 64 * team;
            // ...player 0...10 are non-null, and different among them because of the different playerId
            for (i = off; i < off + 11; i++) assertBN('neq', leafs[0][i], 0, "unexpected teamstate leaf at start of league");
            // ...player 11...25 are identical because we used the same playerId for all of them
            for (i = off + 12; i < off + 25; i++) assertBN('eq',leafs[0][i], leafs[0][off+12], "unexpected teamstate leaf at start of league");
            assertBN('eq', leafs[0][off + 25], 0, "unexpected nonnull tactics leaf at start of league");
            assertBN('eq', leafs[0][off + 26], 0, "unexpected nonnull training leaf at start of league");
            assertBN('eq', leafs[0][off + 27], 0, "unexpected nonnull matchLog leaf at start of league");
            // AFTER first half ---------
            off += 32;
            // ...player 0...10 are non-null, and different among them because of the different playerId
            for (i = off; i < off + 11; i++) assertBN('neq', leafs[0][i], 0, "unexpected teamstate leaf at start of league");
            // ...player 11...25 are identical because we used the same playerId for all of them
            for (i = off + 12; i < off + 25; i++) assertBN('eq', leafs[0][i], leafs[0][off+12], "unexpected teamstate leaf at start of league");
            console.log(leafs[0][off + 26], almostNullTraning)
            console.log(web3.utils.toBN(leafs[0][off + 26]), web3.utils.toBN(almostNullTraning))
            console.log(leafs[0][off + 26].toString(), almostNullTraning.toString())
            assertBN('eq', leafs[0][off + 25], tactics442NoChanges, "unexpected tactics leaf after 1st half of league");
            assertBN('eq', leafs[0][off + 26], almostNullTraning, "unexpected training leaf after 1st half of league");
            assertBN('neq', leafs[0][off + 27], 0, "unexpected null matchLog leaf after 1st half of league");
        }
    });
    
});