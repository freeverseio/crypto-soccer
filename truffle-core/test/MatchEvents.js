const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const logUtils = require('../utils/matchLogUtils.js');
const debug = require('../utils/debugUtils.js');

const Engine = artifacts.require('Engine');
const Assets = artifacts.require('Assets');
const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const EnginePreComp = artifacts.require('EnginePreComp');
const EncodingSkillsSetters = artifacts.require('EncodingSkillsSetters');
const Evolution = artifacts.require('Evolution');

contract('Engine', (accounts) => {
    const UNDEF = undefined;
    const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
    const substitutions = [6, 10, 0];
    const subsRounds = [3, 7, 1];
    const noSubstitutions = [11, 11, 11];
    const lineup0 = [0, 3, 4, 5, 6, 9, 10, 11, 12, 15, 16, 7, 13, 17];
    const lineup1 = [0, 3, 4, 5, 6, 9, 10, 11, 16, 17, 18, 7, 13, 17];
    const lineupConsecutive = Array.from(new Array(14), (x,i) => i); 
    const extraAttackNull =  Array.from(new Array(10), (x,i) => 0);
    const tacticId442 = 0; // 442
    const tacticId433 = 2; // 433
    const playersPerZone442 = [1,2,1,1,2,1,0,2,0];
    const playersPerZone433 = [1,2,1,1,1,1,1,1,1];
    const PLAYERS_PER_TEAM_MAX = 25;
    const firstHalfLog = [0, 0];
    const subLastHalf = false;
    const is2ndHalf = false;
    const isHomeStadium = false;
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
    const dayOfBirthOld = secsToDays(now) - Math.floor(37*365/7);
    const MAX_PENALTY = 10000;
    const DRAW = 2;
    const WINNER_HOME = 0;
    const WINNER_AWAY = 1;
    const teamSumSkillsDefault = 3256244;
    const MAX_GOALS = 12;
    const it2 = async(text, f) => {};
    const trainingPointsDefault = 12;
    
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
        pSkills = await engine.encodePlayerSkills(forceSkills, dayOfBirth21,  gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
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

    beforeEach(async () => {
        encodingSet = await EncodingSkillsSetters.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        encodingLog = await EncodingMatchLog.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        await engine.setPreCompAddr(precomp.address).should.be.fulfilled;
        evolution = await Evolution.new().should.be.fulfilled;
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
    });

    it('hardcoded test - half 1', async () => {
        // We test that this hardcoded-inputs game ends 1-1
        // At start of first half, all players have indentical state.
        team0 = [];
        for (p = 0; p < 25; p++) team0.push('16573429227295117480385309339445376240739796176995438');
        team1 = [];
        for (p = 0; p < 25; p++) team1.push('16573429227295117480385309340654302060354425351701614');
        tactics = '340596594427581673436941882753025';
        seed0 = '0xd971ed54244f38710fc14e8b5fd0ad9491eba0a0b62fb04f575c9875a3d75739';
        homeId = 1;
        awayId = 2;
        seed0 = await engine.generateMatchSeed(seed0, homeId, awayId);
        log0 = '0';
        log1 = '0';
        startTime = 1570147200;
        evs = await engine.playHalfMatch(seed0, startTime, [team0, team1], [tactics, tactics], [log0, log1], [is2nd = false, isHome = true,  playoff = false]).should.be.fulfilled;
        nG = await engine.getNGoals(evs[0]).should.be.fulfilled;
        nG.toNumber().should.be.equal(1);
        nG = await engine.getNGoals(evs[1]).should.be.fulfilled;
        nG.toNumber().should.be.equal(1);
        expectedLog = '205261884989140840566329426675285474869510664919538653543599802679409'
        evs[0].should.be.bignumber.equal(expectedLog)
    });
    
    it('hardcoded test - half 2', async () => {
        // the first half ended 1-1, and the second, 4-1. This test checks that team0 scores 3 goals in the events list.
        // the states we write here are the output states after 1st half. Basically, players who were linedup have a different state
        linedUp = '16573434936285888304224833572589254038720341707981934'
        notLinedUp = '16573429227295117480385309339445376240739796176995438'
        team0 = [];
        for (p = 0; p < 25; p++) team0.push(linedUp);
        for (p = 13; p < 25; p++) team0[p] = notLinedUp;
        for (p = 1; p <= 2; p++) team0[p] = notLinedUp;
        
        linedUp = '16573434936285888304224833573798179858334970882688110'
        notLinedUp = '16573429227295117480385309340654302060354425351701614'
        team1 = [];
        for (p = 0; p < 25; p++) team1.push(linedUp);
        for (p = 13; p < 25; p++) team1[p] = notLinedUp;
        for (p = 1; p <= 2; p++) team1[p] = notLinedUp;

        //       
        tactics = '340596594427581673436941882753025';
        seed0 = '0xd971ed54244f38710fc14e8b5fd0ad9491eba0a0b62fb04f575c9875a3d75739';
        homeId = 1;
        awayId = 2;
        seed0 = await engine.generateMatchSeed(seed0, homeId, awayId);
        log0 = '205261884989140840566329426675285474869510664919538653543599802679409';
        log1 = '205261884989140840566329426675285474869510664919538662550799057420433';
        startTime = 1570147200;

        // play using evolution, which automatically adds training points, and check againsted expected log coming from Go.
        evs2 = await evolution.play2ndHalfAndEvolve(seed0, startTime, [team0, team1], [tactics, tactics], [log0, log1], [is2nd = true, isHome = true,  playoff = false]).should.be.fulfilled;
        nG = await engine.getNGoals(evs2[0]).should.be.fulfilled;
        nG.toNumber().should.be.equal(4)
        expectedLog = '1270126589710522258182323394732234344281996953062426276584982961367525236'
        evs2[0].should.be.bignumber.equal(expectedLog)

        // playing using matchEvents to get the events (besides the logs). Check that the logs, after adding the training points,
        // lead to the same logs as above
        evs = await engine.playHalfMatch(seed0, startTime, [team0, team1], [tactics, tactics], [log0, log1], [is2nd = true, isHome = true,  playoff = false]).should.be.fulfilled;
        nG = await engine.getNGoals(evs[0]).should.be.fulfilled;
        nG.toNumber().should.be.equal(4)
        nG = await engine.getNGoals(evs[1]).should.be.fulfilled;
        nG.toNumber().should.be.equal(1)
        newLogs = await evolution.computeTrainingPoints([evs[0],evs[1]]).should.be.fulfilled;
        newLogs[0].should.be.bignumber.equal(expectedLog)

        // comment/uncomment to get new test results
        // st2 = ''
        // for (e = 0; e < 12; e++) {
        //     st = 'event: ' + e
        //     for (p = 0; p < 5; p++) {
        //         st += " " + evs[2+5*e+p].toNumber();
        //         st2 += ', ' + evs[2+5*e+p].toNumber();
        //     }
        //     if (evs[2+5*e+3].toNumber() == 1) {
        //         st += "  -> GOAL"
        //     }
        //     console.log(st);
        // }
        // console.log(st2);
        
        expected = [1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 7, 1, 7, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 10, 1, 10, 0, 1, 7, 1, 7, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0];
        debug.compareArrays(evs.slice(2), expected, toNum = true, verbose = false);
    });
    
    it('wasPlayerAlignedEndOfLastHalf', async () => {
        seedForRedCard = seed + 83;
        substis = [2, 9, 1];
        rounds = [4, 2, 6];
        // as seen in a test below, there is a redCard for player 9 at round 1
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0).should.be.fulfilled;
        matchEvents = await engine.playHalfMatch(seedForRedCard, now, [teamStateAll50Half1, teamStateAll50Half1], [tactics, tactics], [0, 0], [is2nd = false, isHomeStadium,  playoff = false]).should.be.fulfilled;
        newLog = [];
        newLog[0] = matchEvents[0];
        newLog[1] = matchEvents[1];
        expected = Array.from(new Array(14), (x,i) => true);
        expected[2] = false; 
        expected[1] = false; 
        expected[12] = false; 
        for (p = 0; p < 14; p++) {
            result = await engine.wasPlayerAlignedEndOfLastHalf(p, tactics, newLog[0]).should.be.fulfilled;
            result.should.be.equal(expected[p]);
        }
        // reassuring that the red card was as described above:
        expectedOut = [9, 0];
        expectedOutRounds = [1, 0]; // note that this 1 would be 9 otherwise
        expectedYellows1 = [1, 12];
        expectedYellows2 = [0, 0];
        expectedType = [3, 0]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 2, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [true, false];
        await logUtils.checkExpectedLog(encodingLog, newLog[0], nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt = UNDEF, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
            
        // for each event: 0: teamThatAttacks, 1: managesToShoot, 2: shooter, 3: isGoal, 4: assister
        expected = [
            1, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 
            1, 0, 0, 0, 0, 
            0, 1, 1, 1, 6, 
            0, 0, 0, 0, 0, 
            1, 1, 9, 1, 8, 
            1, 0, 0, 0, 0, 
            1, 1, 9, 1, 9, 
            1, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0
        ];
        debug.compareArrays(matchEvents.slice(2), expected, toNum = true, verbose = false);
    });

    it('computeExceptionalEvents no clashes with redcards', async () => {
        // there is a red card with this seed, to player 9, but he's not involved in any change
        seedForRedCard = seed + 83;
        substis = [2, 6, 1];
        rounds = [4, 2, 6];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half2, tactics, is2nd = true, seedForRedCard).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [0, 9];
        expectedOutRounds = [0, 1];
        expectedYellows1 = [0, 0];
        expectedYellows2 = [1, 12];
        expectedType = [0, 3]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
    });

    
    it('computeExceptionalEvents clashing with redcards before changing player', async () => {
        // there is a red card with this seed, to player 9. Since he's involved in a change, 
        // the round for which he saw the card should be before the proposed change round (2) 
        seedForRedCard = seed + 83;
        substis = [2, 9, 1];
        rounds = [4, 2, 6];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half2, tactics, is2nd = true, seedForRedCard).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [0, 9];
        expectedOutRounds = [0, 1]; // note that this 1 would be 9 otherwise
        expectedYellows1 = [0, 0,];
        expectedYellows2 = [1, 12];
        expectedType = [0, 3]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [1, 2, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
    });

    it('computeExceptionalEvents clashing with redcards after changing player', async () => {
        // there is a red card with this seed, to player 13, which is by definition one of the players to join during the game. 
        // the round for which he saw the card (6) should be after the proposed change round (6 too) 
        seedForRedCardInSubstitutes = seed + 357;
        substis = [2, 9, 1];
        rounds = [4, 2, 6];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half2, tactics, is2nd = true, seedForRedCardInSubstitutes).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [0, 13];
        expectedOutRounds = [0, 6]; // note that it'd be 0, 9 otherwise
        expectedYellows1 = [0, 0];
        expectedYellows2 = [14, 13];
        expectedType = [0, 3]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
    });

    it('computeExceptionalEvents clashing with redcards after changing player forcing last minute', async () => {
        // note that in the first half, player 13 joined, and saw both a yellow and a red card (!!)
        // same as previous but pushing it to the limit, so that the round is 10
        seedForRedCardInSubstitutes = seed + 357;
        substis = [2, 9, 1];
        rounds = [4, 2, 10];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half2, tactics, is2nd = true, seedForRedCardInSubstitutes).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [0, 13];
        expectedOutRounds = [0, 10]; 
        expectedYellows1 = [0, 0];
        expectedYellows2 = [14, 13];
        expectedType = [0, 3]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
    });

    it('computeExceptionalEvents clashing with redcards after changing player forcing last minute (first half)', async () => {
        // note that in the first half, player 13 joined, and saw both a yellow and a red card (!!)
        // same as previous but pushing it to the limit, so that the round is 10
        seedForRedCardInSubstitutes = seed + 357;
        substis = [2, 9, 1];
        rounds = [4, 2, 10];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half1, tactics, is2nd = false, seedForRedCardInSubstitutes).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [13, 0];
        expectedOutRounds = [10, 0];
        expectedYellows1 = [14, 13,];
        expectedYellows2 = [0, 0];
        expectedType = [3, 0]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, true];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
    });
    
    it('check that nDefs is reduced by one when a defender misses in the 2nd half', async () => {
        // note that in the first half, player 13 joined, and saw both a yellow and a red card (!!)
        // same as previous but pushing it to the limit, so that the round is 10
        seedForRedCardInSubstitutes = seed + 357;
        substis = [2, 9, 1];
        rounds = [4, 2, 10];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half1, tactics, is2nd = false, seedForRedCardInSubstitutes).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [13, 0];
        expectedOutRounds = [10, 0];
        expectedYellows1 = [14, 13];
        expectedYellows2 = [0, 0];
        expectedType = [3, 0]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, true];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
           
        // the player with shirt = 1 was substituted by player 13, who was red-carded
        // in the 2nd half there is a defender less than usual
        teamStateAll50Half2[1] = 0;
        seedDraw = 12;
        log2 = await engine.playHalfMatch(seedDraw, now, [teamStateAll50Half2, teamStateAll50Half2], [tactics442NoChanges, tactics442NoChanges], [newLog, newLog], [is2nd = true, isHomeStadium,  playoff = false]).should.be.fulfilled;
        for (team = 0; team < 2; team++){
            nDefs = await encodingLog.getNDefs(log2[team], is2nd = false);
            nDefs.toNumber().should.be.equal(0); // 0 because we did not playHalfMatch in 1st half
            nDefs = await encodingLog.getNDefs(log2[team], is2nd = true);
            nDefs.toNumber().should.be.equal(3); // 3 because it's 1 less than in a 442 tactics
        }   
    });
    
    it('computeExceptionalEvents clashing 2nd against 1st', async () => {
        // first half:
        //      - there is a red card with this seed, to player 9 at round 2. 
        //      - there are two yellow cards, for player 1, and for subtituted 12.
        seedForRedCard = seed + 83;
        substis = [2, 9, 1];
        rounds = [4, 2, 6];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half1, tactics, is2nd = false, seedForRedCard).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [9, 0];
        expectedOutRounds = [1, 0]; // note that this 1 would be 9 otherwise
        expectedYellows1 = [1, 12];
        expectedYellows2 = [0, 0];
        expectedType = [3, 0]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 2, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [true, false];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);

        // second half
        tactics = await engine.encodeTactics(substis = [0,0,0], rounds = [0,0,0], lineupConsecutive, extraAttackNull, tacticsId = 0);
        finalLog = await precomp.computeExceptionalEvents(newLog, teamStateAll50Half2, tactics, is2nd = true, seedForRedCard).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [9, 12]; 
        expectedOutRounds = [1, 1]; // note that this 1 would be 9 otherwise
        expectedYellows1 = [1, 12]; // note that this 1 is OK, he's a different guy, as he was substituted in 1st half
        expectedYellows2 = [1, 14]; // note that this 1 is OK, he's a different guy, as he was substituted in 1st half
        expectedType = [3, 3]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 2, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [true, false];
        await logUtils.checkExpectedLog(encodingLog, finalLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);

    });
    
    it('computeExceptionalEvents clashing 2nd against 1st, with no substitution in the middle', async () => {
        // first half:
        //      - there is a red card with this seed, to player 9 at round 2. 
        //      - there are two yellow cards, for player 1, and for subtituted 12.
        seedForRedCard = seed + 83;
        substis = [2, 3, 4];
        rounds = [4, 2, 6];
        tactics = await engine.encodeTactics(substis, rounds, lineupConsecutive, extraAttackNull, tacticsId = 0);
        newLog = await precomp.computeExceptionalEvents(log = 0, teamStateAll50Half1, tactics, is2nd = false, seedForRedCard).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [9, 0];
        expectedOutRounds = [1, 0]; 
        expectedYellows1 = [1, 12];
        expectedYellows2 = [0, 0];
        expectedType = [3, 0]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, newLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);

        // second half
        tactics = await engine.encodeTactics(substis = [0,0,0], rounds = [0,0,0], lineupConsecutive, extraAttackNull, tacticsId = 0);
        finalLog = await precomp.computeExceptionalEvents(newLog, teamStateAll50Half2, tactics, is2nd = true, seedForRedCard).should.be.fulfilled;
        isHomeSt = false;
        expectedOut = [9, 1]; // note that the red card comes from two yellows.
        expectedOutRounds = [1, 1]; 
        expectedYellows1 = [1, 12]; // note that he'd like to yellow card [1,12] again, but the 1 goes immediately to redCard above.
        expectedYellows2 = [14, 14]; // note that he'd like to yellow card [1,12] again, but the 1 goes immediately to redCard above.
        expectedType = [3, 3]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [1, 1, 1]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, finalLog, nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
    });    

    it('play a match with a special playerId that made it fail before fixing a bug', async () => {
        playerId = 274877907169;
        skills = await assets.getPlayerSkillsAtBirth(playerId).should.be.fulfilled;
        for (i = 0; i< PLAYERS_PER_TEAM_MAX; i++) teamStateAll50Half1[i] = skills;
        result = await engine.playHalfMatch(seed, now, [teamStateAll50Half1, teamStateAll50Half1], [tactics0, tactics0], firstHalfLog, matchBools).should.be.fulfilled;
    });

    it('penaltyPerAge', async () => {
        ageInDays       = [31*365, 31*365+1, 31*365+2, 41*365-4, 41*365-3, 41*365-2, 41*365-1];
        expectedPenalty = [1000000, 998904, 998904, 1544, 1544, 0, 0, 0]
        for (i = 0; i < ageInDays.length; i++) {
            dayOfBirth = Math.round(secsToDays(now) - ageInDays[i]/7);
            playerSkills = await engine.encodePlayerSkills(
                skills = [1,1,1,1,1], 
                dayOfBirth, 
                gen = 0,
                playerId = 2132321,
                [potential = 3,
                forwardness,
                leftishness,
                aggr = 0],
                alignedEndOfLastHalf = true,
                redCardLastGame = false,
                gamesNonStopping = 0,
                injuryWeeksLeft = 0,
                subLastHalf,
                sumSkills = 5
            ).should.be.fulfilled;
            result = await engine.penaltyPerAge(playerSkills, now).should.be.fulfilled;
            result.toNumber().should.be.equal(expectedPenalty[i]);
        }
    });

    it('check that penalties are played in playoff games and excluding redcarded players', async () => {
        // cook data so that the first half ended up in a way that:
        //  - there are red cards
        //  - there are the right goals to, then, in 2nd half, end up in draw.
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
        
        log0 = await logUtils.encodeLog(encodingLog, nGoals = 3, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeStadium, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsDefault);
        
        seedDraw = 12;
        teamStateAll50Half2[9] = 0;
        log2 = await engine.playHalfMatch(seedDraw, now, [teamStateAll50Half2, teamStateAll50Half2], [tactics442, tactics1], [log0, log0], [is2nd = true, isHomeStadium,  playoff = true]).should.be.fulfilled;
        nGoals0 = await encodingLog.getNGoals(log2[0]).should.be.fulfilled;
        nGoals1 = await encodingLog.getNGoals(log2[1]).should.be.fulfilled;
        nGoals0.toNumber().should.be.equal(nGoals1.toNumber());
        expected1 = [true, false, true, true, true, true, true]
        expected2 = [true, false, true, true, true, true, false];
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log2[0], i).should.be.fulfilled;
            pen.should.be.equal(expected1[i]);
            pen = await encodingLog.getPenalty(log2[1], i).should.be.fulfilled;
            pen.should.be.equal(expected2[i]);
        }
        for (team = 0; team < 2; team++){
            win = await encodingLog.getWinner(log2[team]).should.be.fulfilled;
            win.toNumber().should.be.equal(0);
            nDefs = await encodingLog.getNDefs(log2[team], is2nd = false);
            nDefs.toNumber().should.be.equal(4);
            nDefs = await encodingLog.getNDefs(log2[team], is2nd = true);
            nDefs.toNumber().should.be.equal(4);
        }   
    });
    
    it('computePenalties', async () => {
        // one team much better than the other:
        log = await precomp.computePenalties(log = [0,0], [teamStateAll50Half2, teamStateAll1Half2], 50, 1, seed);
        expected = [true, true, true, true, true, false, false]
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log[team = 0], i).should.be.fulfilled;
            pen.should.be.equal(expected[i]);
        }
        expected = [false, false, false, false, false, false, false]
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log[team = 1], i).should.be.fulfilled;
            pen.should.be.equal(expected[i]);
        }
        for (team = 0; team < 2; team++){
            win = await encodingLog.getWinner(log[team]).should.be.fulfilled;
            win.toNumber().should.be.equal(0);
        }   

        // both teams similar:
        log = await precomp.computePenalties(log = [0,0], [teamStateAll50Half2, teamStateAll50Half2], 50, 50, seed);
        expected = [false, true, true, true, true, false, false]
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log[team = 0], i).should.be.fulfilled;
            pen.should.be.equal(expected[i]);
        }
        expected = [true, true, true, true, true, false, false]
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log[team = 1], i).should.be.fulfilled;
            pen.should.be.equal(expected[i]);
        }
        for (team = 0; team < 2; team++){
            win = await encodingLog.getWinner(log[team]).should.be.fulfilled;
            win.toNumber().should.be.equal(1);
        }   

        // both teams really incredible goalkeepers:
        log = await precomp.computePenalties(log = [0,0], [teamStateAll50Half2, teamStateAll50Half2], 5000000, 5000000, seed);
        expected = [false, false, false, false, false, false, false]
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log[team = 0], i).should.be.fulfilled;
            pen.should.be.equal(expected[i]);
        }
        expected = [false, false, false, false, false, false, true]
        for (i = 0; i < 7; i++) {
            pen = await encodingLog.getPenalty(log[team = 1], i).should.be.fulfilled;
            pen.should.be.equal(expected[i]);
        }
        for (team = 0; team < 2; team++){
            win = await encodingLog.getWinner(log[team]).should.be.fulfilled;
            win.toNumber().should.be.equal(1);
        }   
    });

    it('teamSkills are added from 1st to 2nd half', async () => {
        seedDraw = 12;
        subs = [3,1,11];
        tactics442TwoChanges = await engine.encodeTactics(subs, subsRounds, setNoSubstInLineUp(lineupConsecutive, subs), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        subs = [11,11,11];
        tactics442WithNoChanges = await engine.encodeTactics(subs, subsRounds, setNoSubstInLineUp(lineupConsecutive, subs), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        log0 =  await engine.playHalfMatch(seedDraw,  now, [teamStateAll50Half1, teamStateAll50Half1], [tactics442TwoChanges, tactics442WithNoChanges], log = [0, 0], [is2nd = false, isHomeStadium, isPlayoff]).should.be.fulfilled;
        log0 = [log0[0], log0[1]];
        expected = [3250, 2750];
        for (team = 0; team < 2; team++) {
            teamSkills = await encodingLog.getTeamSumSkills(log0[team]).should.be.fulfilled;
            teamSkills.toNumber().should.be.equal(expected[team]);
        }
        subs = [3,11,11];
        tactics442OneChange = await engine.encodeTactics(subs, subsRounds, setNoSubstInLineUp(lineupConsecutive, subs), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        log12 = await engine.playHalfMatch(seedDraw,  now, [teamStateAll50Half2, teamStateAll50Half2], [tactics442OneChange, tactics442WithNoChanges], log0, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expected = [3322, 2750];
        for (team = 0; team < 2; team++) {
            teamSkills = await encodingLog.getTeamSumSkills(log12[team]).should.be.fulfilled;
            teamSkills.toNumber().should.be.equal(expected[team]);
        }
    });
    
    

    it('goals from 1st half are added in the 2nd half', async () => {
        seedDraw = 12;
        log0 =  await engine.playHalfMatch(seedDraw,  now, [teamStateAll50Half1, teamStateAll50Half1], [tactics442NoChanges, tactics1NoChanges], log = [0, 0], [is2nd = false, isHomeStadium, isPlayoff]).should.be.fulfilled;
        log0 = [log0[0], log0[1]];
        log12 = await engine.playHalfMatch(seedDraw,  now, [teamStateAll50Half2, teamStateAll50Half2], [tactics442, tactics1], log0, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // for this seedDraw, they all score one goal in each half
        expected = [1, 1]
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log0[team]);
            nGoals.toNumber().should.be.equal(expected[team]);
        }
        // so the final result should be 2-2
        expected = [2, 2]
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log12[team]);
            nGoals.toNumber().should.be.equal(expected[team]);
            winner = await encodingLog.getWinner(log12[team]);
            winner.toNumber().should.be.equal(DRAW);
            nDefs = await encodingLog.getNDefs(log12[team], is2nd = false);
            nDefs.toNumber().should.be.equal(4);
            nDefs = await encodingLog.getNDefs(log12[team], is2nd = true);
            nDefs.toNumber().should.be.equal(4);
        }
    });

    it('red cards in first half force lineups of 10 players in 2nd half', async () => {
        // choose a seed that gives a red card for player 9.
        seedForRedCard = seed + 83;
        log0 =  await engine.playHalfMatch(seedForRedCard,  now, [teamStateAll50Half1, teamStateAll50Half1], [tactics442NoChanges, tactics1NoChanges], log = [0, 0], [is2nd = false, isHomeStadium, isPlayoff]).should.be.fulfilled;
        log0 = [log0[0], log0[1]]
        isHomeSt = false;
        expectedOut = [9, 0];
        expectedOutRounds = [1, 0]; 
        expectedYellows1 = [1, 10];
        expectedYellows2 = [0, 0];
        expectedType = [3, 0]; // 0 = no event, 3 = redCard
        expectedInGameSubs1 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        expectedInGameSubs2 = [0, 0, 0]; // 0: no subs requested, 1: change takes place, 2: change cancelled
        yellowedCouldNotFinish = [false, false];
        await logUtils.checkExpectedLog(encodingLog, log0[0], nGoals = UNDEF, ass = UNDEF, sho = UNDEF, fwdPos = UNDEF, penalties = UNDEF,
            expectedOut, expectedOutRounds, expectedType, yellowedCouldNotFinish,
            isHomeSt, expectedInGameSubs1, expectedInGameSubs2, expectedYellows1, expectedYellows2, 
            halfTimeSubstitutions = UNDEF, nDefs1 = UNDEF, nDefs2 = UNDEF, nTot = UNDEF, winner = UNDEF, teamSumSkills = UNDEF, trainPo = UNDEF);
        
        teamStateAll50Half2[9] = await encodingSet.setRedCardLastGame(teamStateAll50Half2[9], true);    
        result = await precomp.verifyCanPlay(linedUp = 9, teamStateAll50Half2[9], is2nd = true, isSubst = false).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
        log2 = await engine.playHalfMatch(seedForRedCard, now, [teamStateAll50Half2, teamStateAll50Half2], [tactics442, tactics1], log0, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        for (team = 0; team < 2; team++) {
            nDefs = await encodingLog.getNDefs(log2[team], is2nd = false);
            nDefs.toNumber().should.be.equal(4);
            nDefs = await encodingLog.getNDefs(log2[team], is2nd = true);
            nDefs.toNumber().should.be.equal(4);
            teamSkills = await encodingLog.getTeamSumSkills(log2[team]).should.be.fulfilled;
            teamSkills.toNumber().should.be.equal(2814);
        }
    });
    
    it('play 2nd half with 3 changes is OK, but more than 3 is rejected, by lying in the team-states', async () => {
        // create a 2nd half using 3 players that already played in the 1st half... should work
        messi = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 250).should.be.fulfilled;            
        for (p = 0; p < 3; p++) teamStateAll50Half2[p] = messi; 
        result = await engine.playHalfMatch(seed, now, [teamStateAll50Half2, teamStateAll1Half2], [tactics442NoChanges, tactics442NoChanges], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // create a 2nd half using 4 players that already played in the 1st half... should fail
        teamStateAll50Half2[5] = messi; 
        result = await engine.playHalfMatch(seed, now, [teamStateAll50Half2, teamStateAll1Half2], [tactics442NoChanges, tactics442NoChanges], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.rejected;
    });

    it('play 2nd half with 3 changes is OK, but more than 3 is rejected, by lying in the substitutions', async () => {
        // create a 2nd half using 1 players that already played in the 1st half, and 2 changes only... should work
        messi = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 250).should.be.fulfilled;            
        teamStateAll50Half2[lineupConsecutive[1]] = messi; 
        subst = [3,1,11];
        tactics442TwoChanges = await engine.encodeTactics(subst, subsRounds, setNoSubstInLineUp(lineupConsecutive, subst),
            extraAttackNull, tacticId442).should.be.fulfilled;
        result = await engine.playHalfMatch(seed, now, [teamStateAll50Half2, teamStateAll1Half2], [tactics442TwoChanges, tactics442NoChanges], firstHalfLog, 
            [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // create a 2nd half using 1 players that already played in the 1st half, and 3 changes... should fail
        subst = [3,1,5];
        tactics442ThreeChanges = await engine.encodeTactics(subst, subsRounds, setNoSubstInLineUp(lineupConsecutive, subst),
            extraAttackNull, tacticId442).should.be.fulfilled;
        result = await engine.playHalfMatch(seed, now, [teamStateAll50Half2, teamStateAll1Half2], [tactics442ThreeChanges, tactics442NoChanges], firstHalfLog, 
            [is2nd = true, isHomeStadium, isPlayoff]).should.be.rejected;
    });

    it('play with an injured / red carded / free-slot player', async () => {
        // legit works:
        result = await engine.playHalfMatch(seed, now, [teamStateAll50Half2, teamStateAll1Half2], [tactics442, tactics1], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // red card fails:
        teamStateAll50Half2[5] = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0],
            alignedEndOfLastHalf = false, redCardLastGame = true, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 250).should.be.fulfilled;    

        result = await precomp.verifyCanPlay(linedUp = 9, teamStateAll50Half2[9], is2nd = true, isSubst = false).should.be.fulfilled;
        result.should.not.be.bignumber.equal('0');
        result = await precomp.verifyCanPlay(linedUp = 5, teamStateAll50Half2[5], is2nd = true, isSubst = false).should.be.fulfilled;
        result.should.be.bignumber.equal('0');

        // injured fails
        teamStateAll50Half2[5] = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0],
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 2, subLastHalf, sumSkills = 250).should.be.fulfilled;            
        result = await precomp.verifyCanPlay(linedUp = 5, teamStateAll50Half2[5], is2nd = true, isSubst = false).should.be.fulfilled;
        result.should.be.bignumber.equal('0');
        });

    it('computePenaltyBadPositionAndCondition for GK ', async () => {
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0, gen = 0,  playerId = 232131, [potential = 1,
            forwardness = 0, leftishness = 0, aggr = 0], 
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
        ).should.be.fulfilled;            
        expected = Array.from(new Array(11), (x,i) => MAX_PENALTY);
        expected[0] = 0;
        for (p=0; p < 11; p++) {
            penalty = await precomp.computePenaltyBadPositionAndCondition(p, playersPerZone442, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected[p]);
        }
    });

    it('computePenaltyBadPositionAndCondition for DL ', async () => {
            // for a DL:
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0, gen = 0,  playerId = 312321, [potential = 1,
            forwardness = 1, leftishness = 4, aggr = 0], alignedEndOfLastHalf = false, 
            redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
        ).should.be.fulfilled;            
        expected442 = [MAX_PENALTY, 
            0, 1000, 1000, 2000, 
            1000, 2000, 2000, 3000, 
            3000, 3000 
        ];
        expected433 = [MAX_PENALTY, 
            0, 1000, 1000, 2000, 
            1000, 2000, 3000,  
            2000, 3000, 4000
        ];
        for (p=0; p < 11; p++) {
            penalty = await precomp.computePenaltyBadPositionAndCondition(p, playersPerZone442, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected442[p]);
            penalty = await precomp.computePenaltyBadPositionAndCondition(p, playersPerZone433, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected433[p]);
        }
    });

    it('computePenaltyBadPositionAndCondition for DL with gamesNonStopping', async () => {
        // for a DL:
        expected442 = [MAX_PENALTY, 
            0, 1000, 1000, 2000, 
            1000, 2000, 2000, 3000, 
            3000, 3000 
        ];
        expected433 = [MAX_PENALTY, 
            0, 1000, 1000, 2000, 
            1000, 2000, 3000,  
            2000, 3000, 4000
        ];
        for (games = 1; games < 9; games+=2) {
            playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0, gen = 0,  playerId = 1323121, [potential = 1,
                forwardness = 1, leftishness = 4, aggr = 0], alignedEndOfLastHalf = false, 
                redCardLastGame = false, games, injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
            ).should.be.fulfilled;            
            for (p=0; p < 11; p+=3) {
                penalty = await precomp.computePenaltyBadPositionAndCondition(p, playersPerZone442, playerSkills).should.be.fulfilled;
                if (expected442[p] == MAX_PENALTY) {
                    penalty.toNumber().should.be.equal(0);
                } else {
                    penalty.toNumber().should.be.equal(10000 - Math.min(5000, games*1000) - expected442[p]);
                }
            }
        }
    });


    it('computePenaltyBadPositionAndCondition for MFLCR ', async () => {
        // for a DL:
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0, gen = 0,  playerId = 312321, [potential = 1,
            forwardness = 5, leftishness = 7, aggr = 0], alignedEndOfLastHalf = false, 
            redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
        ).should.be.fulfilled;            
        expected442 = [MAX_PENALTY, 
            1000, 1000, 1000, 1000, 
            0, 0, 0, 0, 
            0, 0 
        ];
        expected433 = expected442;
        for (p=0; p < 11; p++) {
            penalty = await precomp.computePenaltyBadPositionAndCondition(p, playersPerZone442, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected442[p]);
            penalty = await precomp.computePenaltyBadPositionAndCondition(p, playersPerZone433, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected433[p]);
        }
    });
    
    it('teams get tired', async () => {
        const result = await engine.teamsGetTired([10,20,30,40,100], [20,40,60,80,50]).should.be.fulfilled;
        result[0][0].toNumber().should.be.equal(10);
        result[0][1].toNumber().should.be.equal(20);
        result[0][2].toNumber().should.be.equal(30);
        result[0][3].toNumber().should.be.equal(40);
        result[0][4].toNumber().should.be.equal(100);
        result[1][0].toNumber().should.be.equal(10);
        result[1][1].toNumber().should.be.equal(20);
        result[1][2].toNumber().should.be.equal(30);
        result[1][3].toNumber().should.be.equal(40);
        result[1][4].toNumber().should.be.equal(50);
    });
    
    it('play a match in home stadium', async () => {
        log = await engine.playHalfMatch(seed, now, [teamStateAll50Half1, teamStateAll1Half1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHome = true, isPlayoff]).should.be.fulfilled;
        expected = [10, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expected[team]);
        }
    });
    
    it('play a match', async () => {
        log = await engine.playHalfMatch(seed, now, [teamStateAll50Half1, teamStateAll1Half1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expected = [10, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expected[team]);
        }
    });

    it('manages to score with a really old player vs a young one', async () => {
        // a Young Messi manages to score:
        teamState = await createTeamState442(engine, forceSkills= [20,20,20,20,20]).should.be.fulfilled;
        messi = await engine.encodePlayerSkills([100,100,100,100,100], dayOfBirth21, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
        ).should.be.fulfilled;            
        teamState[10] = messi;
        teamThatAttacks = 0;
        log = [0, 0]
        scoreData = await engine.managesToScore(now, 0, teamState, playersPerZone442, extraAttackNull, blockShoot = 20, [kMaxRndNumHalf, kMaxRndNumHalf, kMaxRndNumHalf]).should.be.fulfilled;
        log[teamThatAttacks] = scoreData[0]
        expectedGoals       = [1, 0];
        expectedShooters    = [10, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expectedGoals[team]);
            sho = await encodingLog.getShooter(log[team], 0).should.be.fulfilled;
            sho.toNumber().should.be.equal(expectedShooters[team]);
        }
        // an old Messi does not manage to score:
        oldMessi = await engine.encodePlayerSkills([100,100,100,100,100], dayOfBirthOld, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
        ).should.be.fulfilled;            
        teamState[10] = oldMessi;
        teamThatAttacks = 0;
        log = [0, 0]
        scoreData = await engine.managesToScore(now, 0, teamState, playersPerZone442, extraAttackNull, blockShoot = 20, [kMaxRndNumHalf, kMaxRndNumHalf, kMaxRndNumHalf]).should.be.fulfilled;
        log[teamThatAttacks] = scoreData[0]
        // for this case, there should be a goal, so: 1-0    
        expectedGoals       = [0, 0];
        expectedShooters    = [0, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expectedGoals[team]);
            sho = await encodingLog.getShooter(log[team], 0).should.be.fulfilled;
            sho.toNumber().should.be.equal(expectedShooters[team]);
        }
    });
    
    it('manages to score with select shooter without modifiers', async () => {
        // lets put a Messi and check that it surely scores:
        teamState = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        messi = await engine.encodePlayerSkills([100,100,100,100,100], dayOfBirth21, gen = 0, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedEndOfLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills = 5
        ).should.be.fulfilled;            
        teamState[10] = messi;
        result = await engine.selectShooter(now, teamState, playersPerZone442, extraAttackNull, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(10);
        teamThatAttacks = 0;
        log = [0, 0]
        scoreData = await engine.managesToScore(now, 0, teamState, playersPerZone442, extraAttackNull, blockShoot = 1, [kMaxRndNumHalf, kMaxRndNumHalf, kMaxRndNumHalf]).should.be.fulfilled;
        log[teamThatAttacks] = scoreData[0]         // for this case, there should be a goal, so: 1-0    
        expectedGoals       = [1, 0];
        expectedShooters    = [10, 0];
        expectedAssisters   = [10, 0];
        expectedFwd         = [3, 0];
        // scoreData: 0: matchLog, 1: shooter, 2: isGoal, 3: assister
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expectedGoals[team]);
            ass = await encodingLog.getAssister(log[team], 0).should.be.fulfilled;
            sho = await encodingLog.getShooter(log[team], 0).should.be.fulfilled;
            fwd = await encodingLog.getForwardPos(log[team], 0).should.be.fulfilled;
            ass.toNumber().should.be.equal(expectedAssisters[team]);
            sho.toNumber().should.be.equal(expectedShooters[team]);
            fwd.toNumber().should.be.equal(expectedFwd[team]);
            if (team == teamThatAttacks) {
                scoreData[1].toNumber().should.be.equal(expectedShooters[team])
                scoreData[2].toNumber().should.be.equal(1)
                scoreData[3].toNumber().should.be.equal(expectedAssisters[team])
            }
        }
        // let's put a radically good GK, and check that it doesn't score
        log = [0, 0]
        teamThatAttacks = 0;
        scoreData = await engine.managesToScore(now, 0, teamState, playersPerZone442, extraAttackNull, blockShoot = 1000, [kMaxRndNumHalf, kMaxRndNumHalf, kMaxRndNumHalf]).should.be.fulfilled;
        log[teamThatAttacks] = scoreData[0]
        expectedGoals       = [0, 0];
        expectedShooters    = [0, 0];
        expectedAssisters   = [0, 0];
        expectedFwd         = [0, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expectedGoals[team]);
            ass = await encodingLog.getAssister(log[team], 0).should.be.fulfilled;
            sho = await encodingLog.getShooter(log[team], 0).should.be.fulfilled;
            fwd = await encodingLog.getForwardPos(log[team], 0).should.be.fulfilled;
            ass.toNumber().should.be.equal(expectedAssisters[team]);
            sho.toNumber().should.be.equal(expectedShooters[team]);
            fwd.toNumber().should.be.equal(expectedFwd[team]);
        }
        // Finally, check that even with a super-goalkeeper, there are chances of scoring (e.g. if the rnd is super small, in this case)
        log = [0, 0]
        scoreData = await engine.managesToScore(now, 0, teamState, playersPerZone442, extraAttackNull, blockShoot = 1000, [kMaxRndNumHalf, 1, kMaxRndNumHalf]).should.be.fulfilled;
        log[teamThatAttacks] = scoreData[0]
        expectedGoals       = [1, 0];
        expectedShooters    = [10, 0];
        expectedAssisters   = [10, 0];
        expectedFwd         = [3, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(log[team]);
            nGoals.toNumber().should.be.equal(expectedGoals[team]);
            ass = await encodingLog.getAssister(log[team], 0).should.be.fulfilled;
            sho = await encodingLog.getShooter(log[team], 0).should.be.fulfilled;
            fwd = await encodingLog.getForwardPos(log[team], 0).should.be.fulfilled;
            ass.toNumber().should.be.equal(expectedAssisters[team]);
            sho.toNumber().should.be.equal(expectedShooters[team]);
            fwd.toNumber().should.be.equal(expectedFwd[team]);
        }
    });

    it('select shooter with modifiers', async () => {
        teamState = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        extraAttack = [
            true, false, false, true,
            false, true, true, false,
            true, false,
        ];
        expectedRatios = [1,
            15000, 5000, 5000, 15000,
            25000, 50000, 50000, 25000,
            75000, 75000
        ]
        sum = expectedRatios.reduce((a,b) => a + b, 0)
        k = 0;
        for (p = 0; p < 11; p++) {
            k += Math.floor(MAX_RND*expectedRatios[p]/sum);
            result = await engine.selectShooter(now, teamState, playersPerZone442, extraAttack, k).should.be.fulfilled;
            result.toNumber().should.be.equal(p);
            if (p < 10) {
                result = await engine.selectShooter(now, teamState, playersPerZone442, extraAttack, k + p + 1).should.be.fulfilled;
                result.toNumber().should.be.equal(p+1);
            }
        }
    });
    
    it('select assister with modifiers', async () => {
        console.log("warning: This test takes a few secs...")
        teamState = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        extraAttack = [
            true, false, false, true,
            false, true, true, false,
            true, false,
        ];
        nPartitions = 200;
        expectedTrans = [ 5, 65, 15, 20, 65, 80, 110, 115, 220, 155, 150 ];
        transtions = [];
        t=0;
        rndOld = 0; 
        result = await engine.selectAssister(now, teamState, playersPerZone442, extraAttack, shooter = 8, rnd = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        prev = result.toNumber();
        for (p = 0; p < nPartitions; p++) {
            rnd = Math.floor(p * MAX_RND/ nPartitions);
            result = await engine.selectAssister(now, teamState, playersPerZone442, extraAttack, shooter = 8, rnd).should.be.fulfilled;
            if (result.toNumber() != prev) {
                percentageForPrevPlayer = Math.round((rnd-rndOld)/MAX_RND*1000);
                // console.log(prev, percentageForPrevPlayer);
                transtions.push(percentageForPrevPlayer);
                prev = result.toNumber();
                t++;
                rndOld = rnd;
            }
        }
        percentageForPrevPlayer = Math.round((MAX_RND-rndOld)/MAX_RND*1000);
        // console.log(prev, percentageForPrevPlayer);
        transtions.push(percentageForPrevPlayer);
        for (t = 0; t < expectedTrans.length; t++) {
            (result.toNumber()*0 + transtions[t]).should.be.equal(expectedTrans[t]);
        }
    });

    it('select assister with modifiers and one Messi', async () => {
        console.log("warning: This test takes a few secs...")
        teamState = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        messi = await engine.encodePlayerSkills([2,2,2,2,2], dayOfBirth21, gen = 0, id = 1323121, [pot = 3, fwd = 3, left = 7, aggr = 0],
            alignedEndOfLastHalf = false, redCardLastGame = false, 
            gamesNonStopping = 0, injuryWeeksLeft = 0, subLastHalf, sumSkills = 10).should.be.fulfilled;            
        teamState[8] = messi;
        extraAttack = [
            true, false, false, true,
            false, true, true, false,
            true, false,
        ];
        nPartitions = 200;
        expectedTrans = [ 5, 40, 10, 10, 40, 45, 70, 70, 530, 90, 90 ];
        transtions = [];
        t=0;
        rndOld = 0;
        result = await engine.selectAssister(now, teamState, playersPerZone442, extraAttack, shooter = 8, rnd = 0).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        prev = result.toNumber();
        for (p = 0; p < nPartitions; p++) {
            rnd = Math.floor(p * MAX_RND/ nPartitions);
            result = await engine.selectAssister(now, teamState, playersPerZone442, extraAttack, shooter = 8, rnd).should.be.fulfilled;
            if (result.toNumber() != prev) {
                percentageForPrevPlayer = Math.round((rnd-rndOld)/MAX_RND*1000);
                // console.log(prev, percentageForPrevPlayer);
                transtions.push(percentageForPrevPlayer);
                prev = result.toNumber();
                t++;
                rndOld = rnd;
            }
        }
        percentageForPrevPlayer = Math.round((MAX_RND-rndOld)/MAX_RND*1000);
        // console.log(prev, percentageForPrevPlayer);
        transtions.push(percentageForPrevPlayer);
        // console.log(transtions)
        for (t = 0; t < expectedTrans.length; t++) {
            (result.toNumber()*0 + transtions[t]).should.be.equal(expectedTrans[t]);
        }
    });

    it('throws dice array11 fine grained testing', async () => {
        // interface: throwDiceArray(uint[11] memory weights, uint rndNum)
        weights = Array.from(new Array(11), (x,i) => 100);
        sum = 100 * 11;
        k = 0;
        for (p = 0; p < 11; p++) {
            k += Math.floor(MAX_RND*weights[p]/sum);
            result = await engine.throwDiceArray(weights, k).should.be.fulfilled;
            result.toNumber().should.be.equal(p);
            if (p < 10) {
                result = await engine.throwDiceArray(weights, k+p+1).should.be.fulfilled;
                result.toNumber().should.be.equal(p+1);
            }
        }
    });

    it('throws dice array11', async () => {
        // interface: throwDiceArray(uint[11] memory weights, uint rndNum)
        weights = Array.from(new Array(11), (x,i) => 1);
        weights[8] = 1000;
        let result = await engine.throwDiceArray(weights, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(8);
        weights[8] = 1;
        weights[9] = 1000;
        result = await engine.throwDiceArray(weights, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(9);
        weights[9] = 1;
        weights[10] = 1000;
        result = await engine.throwDiceArray(weights, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(10);
    });
    
    it('manages to shoot', async () => {
        // interface: managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum)
        let  = 8000; // the max allowed random number is 16383, so this is about half of it
        let globSkills = [[100,100,100,100,100], [1,1,1,1,1]];
        let result = await engine.managesToShoot(0,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
        result = await engine.managesToShoot(1,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        globSkills = [[1,1,1,1,1], [100,100,100,100,100]];
        result = await engine.managesToShoot(0,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        result = await engine.managesToShoot(1,globSkills,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('throws dice', async () => {
        // interface: throwDice(uint weight1, uint weight2, uint rndNum)
        let result = await engine.throwDice(1,10,kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await engine.throwDice(10,1,kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDice(10,10,kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDice(10,10,2*kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
    });


    it('gets n rands from a seed', async () => {
        ROUNDS_PER_MATCH = await engine.ROUNDS_PER_MATCH().should.be.fulfilled
        const result = await engine.getNRandsFromSeed(seed, 4*ROUNDS_PER_MATCH).should.be.fulfilled;
        expectLen = 4*ROUNDS_PER_MATCH.toNumber();
        result.length.should.be.equal(expectLen);
        prevRnds = [];
        // checks that all rnds are actually different:
        for (r = 0; r < result.length; r++) {
            for (prev = 0; prev < prevRnds.length; prev++){
                result[r].should.be.bignumber.not.equal(prevRnds[prev]);
            }
            prevRnds.push(result[r]);
        }
        result[0].should.be.bignumber.equal("32690506113");
        result[expectLen-1].should.be.bignumber.equal("62760289461");
    });

    it('computes team global skills by aggregating across all players in team', async () => {
        // If all skills where 1 for all players, and tactics = 442 =>
        // move2attack =    defence(defenders + 2*midfields + attackers) +
        //                  speed(defenders + 2*midfields) +
        //                  pass(defenders + 3*midfields) 
        //             =    14 + 12 + 16 = 42
        // createShoot =    speed(attackers) + pass(attackers) = 2 + 2 = 4
        // defendShoot =    speed(defenders) + defence(defenders) = 4 + 4 = 8 
        // blockShoot  =    shoot(keeper); 1
        // endurance   =    70;
        // attackersSpeed = [1,1]
        // attackersShoot = [1,1]
        
        teamState442 = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        globSkills = await precomp.getTeamGlobSkills(teamState442, playersPerZone442, extraAttackNull, now).should.be.fulfilled;
        expectedGlob = [42, 4, 8, 1, 70];
        for (g = 0; g < 5; g++) globSkills[g].toNumber().should.be.equal(expectedGlob[g]);
    });

    it('getLineUpAndPlayerPerZone', async () => {
        teamState442 = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        result = await engine.getLineUpAndPlayerPerZone(teamState442, tactics1, is2ndHalf, log = [0,0], seed).should.be.fulfilled;
        let {0: matchLog, 1: states} = result;
        for (p = 0; p < 11; p++) states[p].should.be.bignumber.equal(teamState442[lineupConsecutive[p]]);
    });

    it('play match with wrong tactic', async () => {
        tacticsWrong = await engine.encodeTactics(substitutions, subsRounds, lineup1, extraAttackNull, tacticIdTooLarge = 6);
        await engine.playHalfMatch(seed, now, teamStateAll50Half1, teamStateAll50Half1, [tacticsWrong, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.rejected;
    });


    it('different team state => different result', async () => {
        matchLog = await engine.playHalfMatch(123456, now, [teamStateAll50Half1, teamStateAll50Half1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [2, 2];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(matchLog[team]);
            nGoals.toNumber().should.be.equal(expectedResult[team]);
        }

        matchLog = await engine.playHalfMatch(123456, now, [teamStateAll50Half1, teamStateAll1Half1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [11, 0];
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(matchLog[team]);
            nGoals.toNumber().should.be.equal(expectedResult[team]);
        }
    });

    it('different seeds => different result', async () => {
        matchLog = await engine.playHalfMatch(123456, now, [teamStateAll50Half1, teamStateAll50Half1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [2, 2];
        result = []
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(matchLog[team]);
            result.push(nGoals);
        }
        debug.compareArrays(result, expectedResult, toNum = true, verbose = false);
        matchLog = await engine.playHalfMatch(654322, now, [teamStateAll50Half1, teamStateAll50Half1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [1, 1];
        result = []
        for (team = 0; team < 2; team++) {
            nGoals = await encodingLog.getNGoals(matchLog[team]);
            result.push(nGoals);
        }
        debug.compareArrays(result, expectedResult, toNum = true, verbose = false);
        // for each event: 0: teamThatAttacks, 1: managesToShoot, 2: shooter, 3: isGoal, 4: assister
        expected = [ 
            1, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 
            1, 1, 7, 1, 9, 
            0, 0, 0, 0, 0, 
            0, 0, 0, 0, 0, 
            1, 0, 0, 0, 0, 
            0, 1, 10, 1, 10, 
            1, 0, 0, 0, 0, 
            1, 0, 0, 0, 0, 
            1, 1, 8, 0, 0, 
            1, 1, 10, 0, 0, 
            1, 0, 0, 0, 0 
        ];
        goals = [0,0];
        for (i=0;i< expected.length/5;i++) goals[expected[5*i]] += expected[5*i+3] + 0*result[0] ;
        debug.compareArrays(goals, expectedResult, toNum = false, verbose = false);
        debug.compareArrays(matchLog.slice(2), expected, toNum = true, verbose = false);
    });
});