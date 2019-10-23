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
    });
    
    it('test1', async () => {
        assistersIdx = Array.from(new Array(MAX_GOALS), (x,i) => i);
        shootersIdx  = Array.from(new Array(MAX_GOALS), (x,i) => 1);
        shooterForwardPos  = Array.from(new Array(MAX_GOALS), (x,i) => 1);
        penalties  = Array.from(new Array(7), (x,i) => false);
        typesOutOfGames = [3, 0];
        outOfGameRounds = [7, 0];
        yellowCardedDidNotFinish1stHalf = [false, false];
        ingameSubs1 = [0, 0, 0]
        ingameSubs2 = [0, 0, 0]
        outOfGames = [9, 14]
        yellowCards1 = [14, 0]
        yellowCards2 = [0, 0]
        halfTimeSubstitutions = [14, 14, 14]
        nDefs1 = 4; 
        nDefs2 = 0; 
        nTot = 10; 
        winner = 0; 
        isHomeSt = true;
        teamSumSkillsDefault = 0;
        trainingPointsInit = 0;
        
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 3, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        expected = [30, 9];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });

});