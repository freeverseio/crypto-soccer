const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const logUtils = require('../utils/matchLogUtils.js');
const debug = require('../utils/debugUtils.js');

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
    
    const it2 = async(text, f) => {};

    function setNoSubstInLineUp(lineup, substitutions) {
        modifiedLineup = [...lineup];
        NO_SUBST = 11;
        for (s = 0; s < 3; s++) {
            if (substitutions[s] == NO_SUBST) modifiedLineup[s + 11] = 25 + s;
        }
        return modifiedLineup;
    }

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
            pSkills = await engine.encodePlayerSkills(forceSkills, dayOfBirth21, gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedEndOfLastHalfTwoVec[0], redCardLastGame, gamesNonStopping, 
                injuryWeeksLeft, subLastHalf, sumSkills).should.be.fulfilled 
            teamState.push(pSkills)
        }
        p = 10;
        pSkills = await engine.encodePlayerSkills(forceSkills, dayOfBirth21, gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
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
            skills, dayOfBirth21, gen = 0, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[0], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 0; player < 11; player++) {
            teamState.push(playerStateTemp)
        }

        playerStateTemp = await engine.encodePlayerSkills(
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

    
    beforeEach(async () => {
        evolution = await Evolution.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        encodeLog = await EncodingMatchLog.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        await engine.setPreCompAddr(precomp.address).should.be.fulfilled;
        await evolution.setAssetsAddress(assets.address).should.be.fulfilled;
        await evolution.setEngine(engine.address).should.be.fulfilled;
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
        POINTS_FOR_HAVING_PLAYED = await evolution.POINTS_FOR_HAVING_PLAYED().should.be.fulfilled;
        POINTS_FOR_HAVING_PLAYED = POINTS_FOR_HAVING_PLAYED.toNumber();
        MAX_WEIGHT = await evolution.MAX_WEIGHT().should.be.fulfilled;
        MAX_WEIGHT = MAX_WEIGHT.toNumber();
        MIN_WEIGHT = await evolution.MIN_WEIGHT().should.be.fulfilled;
        MIN_WEIGHT = MIN_WEIGHT.toNumber();
    });

    it('evolution leading to a roster player', async () => {
        playerSkills = await engine.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 3,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 40;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        weights = [10, 20, 30, 40, 50];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime);

        // checks that the generation increases by 1. 
        // It sets a "32" at the beginning if it is a Roster player, otherwise it is a child
        // In this case, the randomness leads to a Roster player
        result = await assets.getGeneration(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(32 + gen + 1)

        expected = [ 809, 1199, 947, 799, 1244 ];
        results = []
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getPass(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        results.push(result)
        debug.compareArrays(results, expected, toNum = true, verbose = false);
        
        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
        
    });

    it('evolution leading to an actual son', async () => {
        // all inputs are identical to the previous test, except for a +2 in matchStatTime,
        // which changes the entire randomness
        playerSkills = await engine.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 3,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 40;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        weights = [10, 20, 30, 40, 50];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime + 2);

        // checks that the generation increases by 1. It sets a "32" at the beginning if it is a Roster player, otherwise it is a child.
        // In this case, randomness leads to a son.
        result = await assets.getGeneration(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(0 * 32 + gen + 1)

        expected = [ 743, 1459, 1070, 757, 970 ];
        results = []
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getPass(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        results.push(result)
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        results.push(result)
        debug.compareArrays(results, expected, toNum = true, verbose = false);
        
        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
        
    });
    
    
    it('getTeamEvolvedSkills', async () => {
        weights = Array.from(new Array(25), (x,i) => 3*i % 14);
        specialPlayer = 21;
        // make sure they sum to MAX_WEIGHT:
        for (bucket = 0; bucket < 5; bucket++){
            sum4 = 0;
            for (sk = 5 * bucket; sk < (5 * bucket + 4); sk++) {
                sum4 += weights[sk];
            }
            weights[5 * bucket + 4] = MAX_WEIGHT - sum4;
        }        
        assignment = await evolution.encodeTP(weights, specialPlayer).should.be.fulfilled;
        matchStartTime = now;
        newSkills = await evolution.getTeamEvolvedSkills(teamStateAll50Half2, 200, assignment, matchStartTime);
        for (p = 0; p < 25; p++) {
            result = await evolution.getShoot(newSkills[p]);
            if (p == specialPlayer) result.toNumber().should.be.equal(77);
            else result.toNumber().should.be.equal(74);
        }
    });
    
    it('getTeamEvolvedSkills with realistic team and zero TPs', async () => {
        teamState = createHardcodedTeam();
        for (p = 18; p < 25; p++) teamState.push(0);
        weights = [ 0, 3, 6, 9, 57, 1, 4, 7, 10, 53, 2, 5, 8, 11, 49, 3, 6, 9, 12, 45, 4, 7, 10, 13, 41 ];
        assignment = await evolution.encodeTP(weights, specialPlayer = 12).should.be.fulfilled;
        matchStartTime = now;
        newSkills = await evolution.getTeamEvolvedSkills(teamState, TPs = 0, assignment, matchStartTime);
        initShoot = [];
        newShoot = [];
        expectedNewShoot  = [ 623, 440, 829, 811, 723, 702, 554, 735, 815, 1466, 680, 930, 1181, 1095, 697, 622, 566, 931 ];
        expectedInitShoot = [ 623, 440, 829, 811, 723, 729, 554, 751, 815, 1474, 680, 930, 1181, 1103, 697, 622, 566, 931 ];
        for (p = 0; p < 18; p++) {
            result0 = await evolution.getShoot(teamState[p]);
            result1 = await evolution.getShoot(newSkills[p]);
            initShoot.push(result0)
            newShoot.push(result1)
        }
        debug.compareArrays(newShoot, expectedNewShoot, toNum = true, verbose = false);
        debug.compareArrays(initShoot, expectedInitShoot, toNum = true, verbose = false);
    });
    
    it('getTeamEvolvedSkills with realistic team and non-zero TPs', async () => {
        teamState = createHardcodedTeam();
        for (p = 18; p < 25; p++) teamState.push(0);
        weights = [ 0, 3, 6, 9, 57, 1, 4, 7, 10, 53, 2, 5, 8, 11, 49, 3, 6, 9, 12, 45, 4, 7, 10, 13, 41 ];
        assignment = await evolution.encodeTP(weights, specialPlayer = 12).should.be.fulfilled;
        matchStartTime = now;
        newSkills = await evolution.getTeamEvolvedSkills(teamState, TPs = 30, assignment, matchStartTime);
        initShoot = [];
        newShoot = [];
        expectedNewShoot  = [ 624, 441, 830, 819, 737, 703, 555, 736, 825, 1468, 691, 943, 1199, 1097, 710, 634, 568, 943 ];
        expectedInitShoot = [ 623, 440, 829, 811, 723, 729, 554, 751, 815, 1474, 680, 930, 1181, 1103, 697, 622, 566, 931 ];
        for (p = 0; p < 18; p++) {
            result0 = await evolution.getShoot(teamState[p]);
            result1 = await evolution.getShoot(newSkills[p]);
            initShoot.push(result0)
            newShoot.push(result1)
        }
        debug.compareArrays(newShoot, expectedNewShoot, toNum = true, verbose = false);
        debug.compareArrays(initShoot, expectedInitShoot, toNum = true, verbose = false);
    });

    it('test evolvePlayer at zero potential', async () => {
        playerSkills = await engine.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 0, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 16;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        weights = [10, 20, 30, 40, 50];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime);
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        expected = [102,104,106,108,110];
        result.toNumber().should.be.equal(expected[0]);
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[1]);
        result = await engine.getPass(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[2]);
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[3]);
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[4]);
    });
    
    it('test evolvePlayer with TPs= 0', async () => {
        playerSkills = await engine.encodePlayerSkills(
            skills = [12, 13, 155, 242, 32], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 6, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 16;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 0;
        weights = [10, 20, 30, 40, 50];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime);
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        expected = skills;
        result.toNumber().should.be.equal(expected[0]);
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[1]);
        result = await engine.getPass(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[2]);
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[3]);
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[4]);
    });
    
    
    it('test evolvePlayer at non-zero potential', async () => {
        playerSkills = await engine.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 1, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 16;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        weights = [10, 20, 30, 40, 50];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime);
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        expected = [102,105,107,110,113];
        result.toNumber().should.be.equal(expected[0]);
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[1]);
        result = await engine.getPass(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[2]);
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[3]);
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[4]);

        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
    });

    it('test evolvePlayer at non-zero potential and age', async () => {
        playerSkills = await engine.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 17;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        weights = [10, 20, 30, 40, 50];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime);
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        expected = [104, 108, 112, 117, 121];
        result.toNumber().should.be.equal(expected[0]);
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[1]);
        result = await engine.getPass(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[2]);
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[3]);
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[4]);
        
        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
    });

    it('test evolvePlayer with old age', async () => {
        playerSkills = await engine.encodePlayerSkills(
            skills = [1000, 2000, 3000, 4000, 5000], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 35;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        weights = [0, 0, 0, 0, 0];
        newSkills = await evolution.evolvePlayer(playerSkills, TPs, weights, matchStartTime);
        result = await engine.getShoot(newSkills).should.be.fulfilled;
        expected = [968, 1968, 2968, 3968, 4968]; // -32 per game
        result.toNumber().should.be.equal(expected[0]);
        result = await engine.getSpeed(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[1]);
        result = await engine.getPass(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[2]);
        result = await engine.getDefence(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[3]);
        result = await engine.getEndurance(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expected[4]);
    });

    it('test that we can a 2nd half and include the evolution points too', async () => {
        matchLog = await evolution.play2ndHalfAndEvolve(
            123456, now, [teamStateAll50Half2, teamStateAll50Half2], [tactics0, tactics1], [0, 0], 
            [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;

            expectedResult = [2, 2];
            expectedPoints = [15, 18];
            for (team = 0; team < 2; team++) {
            nGoals = await encodeLog.getNGoals(matchLog[team]);
            nGoals.toNumber().should.be.equal(expectedResult[team]);
            points = await encodeLog.getTrainingPoints(matchLog[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expectedPoints[team]);
        }
    });

    it('training points: estimate cost', async () => {
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 0, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        logFinal = await evolution.computeTrainingPointsWithCost([log0, log0]);
    });    
        
    it('training points with random inputs', async () => {
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
        expected = [36, 25];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });

    it('training points with no goals nor anything else', async () => {
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 0, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + cleanSheet(24+8) = 42
        expected = [42, 42];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals by attackers', async () => {
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
        }
    });    

    it('training points with many goals by mids', async () => {
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 6);
        shoot   = Array.from(new Array(goals), (x,i) => 6);
        fwd     = Array.from(new Array(goals), (x,i) => 2);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_MIDS(5 * 5) - GOALS_OPPONENT(5)  
        expected = [30, 30];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals by defs with assists', async () => {
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 6);
        shoot   = Array.from(new Array(goals), (x,i) => 2);
        fwd     = Array.from(new Array(goals), (x,i) => 1);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_DEFS(4 * 5) + ASSISTS(3*5) - GOALS_OPPONENT(5)  
        expected = [50, 50];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
            // console.log(points.toNumber())//.should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals with a winner at home', async () => {
        win = 0;
        isHome = true;

        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);

        goals = 4;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);
            
        logFinal = await evolution.computeTrainingPoints([log0, log1])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + WIN_AT_HOME(11) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(4)  
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 4) - GOALS_OPPONENT(5)  
        expected = [37, 21];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals with a winner away', async () => {
        win = 1;
        isHome = true;

        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);

        goals = 6;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);
            
        logFinal = await evolution.computeTrainingPoints([log0, log1])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(6)  
        // expect: POINTS_FOR_HAVING_PLAYED(10) + WIN_AWAY(22) + GOALS_BY_ATTACKERS(4 * 6) - GOALS_OPPONENT(5)  
        expected = [24, 51];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
            // console.log(points.toNumber())//.should.be.equal(expected[team]);
        }
    });    
    
    it('training points with no goals but cards', async () => {
        outGames    = [4, 6];
        types       = [RED_CARD, RED_CARD];
        yellows1    = [3, 7];
        yellows2    = [1, 2];
        
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 0, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outGames, outOfGameRounds, types, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellows1, yellows2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await evolution.computeTrainingPoints([log0, log0])
        // expect: POINTS_FOR_HAVING_PLAYED(10) + cleanSheet(23+8) - REDS(3*2) - YELLOWS(4) 
        expected = [31, 31];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    
    
    it('training points with many goals by attackers... and different teamSumSkills', async () => {
        // first get the resulting Traning points with teamSkills difference: [25, 25]
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
        }

        // second: get the resulting Traning points with teamSkills difference
        teamSumSkills = 1000;
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
        teamSumSkills = 2000;
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
            
        logFinal = await evolution.computeTrainingPoints([log0, log1])
        expected = [50, 12];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
        // third: same as above but inverse
        teamSumSkills = 2000;
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
        teamSumSkills = 1000;
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
            
        logFinal = await evolution.computeTrainingPoints([log0, log1])
        expected = [12, 50];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }

    });    
});