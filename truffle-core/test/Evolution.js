const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const logUtils = require('./matchLogUtils.js');

const Evolution = artifacts.require('Evolution');
const Assets = artifacts.require('Assets');
const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');

contract('Evolution', (accounts) => {
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
    
    const it2 = async(text, f) => {};
       
    function daysToSecs(dayz) {
        return (dayz * 24 * 3600); 
    }

    function secsToDays(secs) {
        return secs/ (24 * 3600);
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
            pSkills = await engine.encodePlayerSkills(forceSkills, dayOfBirth21, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedEndOfLastHalfTwoVec[0], redCardLastGame, gamesNonStopping, 
                injuryWeeksLeft, subLastHalf, sumSkills).should.be.fulfilled 
            teamState.push(pSkills)
        }
        p = 10;
        pSkills = await engine.encodePlayerSkills(forceSkills, dayOfBirth21, playerId + p, [pot, fwd442[p], left442[p], aggr],
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

    beforeEach(async () => {
        evolution = await Evolution.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        encodeLog = await EncodingMatchLog.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        await engine.setCardsAndInjuries(precomp.address).should.be.fulfilled;
        tactics0 = await engine.encodeTactics(substitutions, subsRounds, lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        tactics1 = await engine.encodeTactics(substitutions, subsRounds, lineupConsecutive, extraAttackNull, tacticId433).should.be.fulfilled;
        tactics1NoChanges = await engine.encodeTactics(noSubstitutions, subsRounds, lineupConsecutive, extraAttackNull, tacticId433).should.be.fulfilled;
        tactics442 = await engine.encodeTactics(substitutions, subsRounds, lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        tactics442NoChanges = await engine.encodeTactics(noSubstitutions, subsRounds, lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        teamStateAll50Half1 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2, aligned = [false, false]).should.be.fulfilled;
        teamStateAll1Half1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2, aligned = [false, false]).should.be.fulfilled;
        teamStateAll50Half2 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2, aligned = [true, false]).should.be.fulfilled;
        teamStateAll1Half2 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2, aligned = [true, false]).should.be.fulfilled;
        MAX_RND = await engine.MAX_RND().should.be.fulfilled;
        MAX_RND = MAX_RND.toNumber();
        kMaxRndNumHalf = Math.floor(MAX_RND/2)-200; 
        events1Half = Array.from(new Array(7), (x,i) => 0);
        events1Half = [events1Half,events1Half];
        POINTS_FOR_HAVING_PLAYED = await evolution.POINTS_FOR_HAVING_PLAYED().should.be.fulfilled;
        POINTS_FOR_HAVING_PLAYED = POINTS_FOR_HAVING_PLAYED.toNumber();
        
        
    });
    
    it2('training points with random inputs', async () => {
        typeOut = [3, 0];
        outRounds = [7, 0];
        outGames = [9, 14]
        yellows1 = [14, 0]
        yellows2 = [0, 0]
        defs1 = 4; 
        defs2 = 0; 
        numTot = 10; 
        win = 0; 
        isHome = true;
        
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 3, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outGames, outRounds, typeOut, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellows1, yellows2, 
            halfTimeSubstitutions, defs1, defs2, numTot, win, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        expected = [30, 19];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });

    it2('training points with no goals nor anything else', async () => {
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 0, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + cleanSheet(25+8) = 42
        expected = [42, 42];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
            // console.log(points.toNumber())//.should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals but bothing else', async () => {
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(5)  
        expected = [25, 25];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
            // console.log(points.toNumber())//.should.be.equal(expected[team]);
        }
    });    
});