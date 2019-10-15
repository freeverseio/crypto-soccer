const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Engine = artifacts.require('Engine');
const Assets = artifacts.require('Assets');
const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const EnginePreComp = artifacts.require('EnginePreComp');

contract('Engine', (accounts) => {
    // const seed = 610106;
    const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
    const substitutions = [6, 10, 0];
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
    const MAX_PENALTY = 10000;

    const createTeamState = async (seed, engine, assets, forceSkills, forceFwd, forceLeft) => {
        teamState = []
        for (p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            result = await assets.computeSkills(seed, shirtNum = p).should.be.fulfilled;
            let {0: skills, 1: potential, 2: forwardness, 3: leftishness} = result;
            if (forceSkills) skills = forceSkills;
            if (forceFwd) forwardness = forceFwd;
            if (forceLeft) leftishness = forceLeft;
            playerSkillsTemp = await engine.encodePlayerSkills(
                skills, 
                monthOfBirth = 0, 
                playerId = 1312312,
                [potential,
                forwardness,
                leftishness,
                aggressiveness = 0],
                alignedLastHalf = true,
                redCardLastGame = false,
                gamesNonStopping = 0,
                injuryWeeksLeft = 0
            ).should.be.fulfilled;            
            teamState.push(playerSkillsTemp)
        }        
        return teamState;
    };

    function daysToSecs(dayz) {
        return (dayz * 24 * 3600); 
    }

    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }

    
    const createTeamState442 = async (engine, forceSkills) => {
        teamState = [];
        playerId = 123121;
        pot = 3;
        aggr = 0;
        alignedLastHalf = true;
        redCardLastGame = false;
        gamesNonStopping = 0;
        injuryWeeksLeft = 0;
        for (p = 0; p < 11; p++) {
            pSkills = await engine.encodePlayerSkills(forceSkills, dayOfBirth21, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft).should.be.fulfilled 
            teamState.push(pSkills)
        }
        for (p = 11; p < PLAYERS_PER_TEAM_MAX; p++) {
            teamState.push(pSkills)
        }        
        return teamState;
    };

    const createTeamStateFromSinglePlayer = async (skills, engine, forwardness = 3, leftishness = 2) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            dayOfBirth21,
            playerId = 2132321,
            [potential = 3,
            forwardness,
            leftishness,
            aggr = 0],
            alignedLastHalf = true,
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

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        encodingLog = await EncodingMatchLog.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        await engine.setCardsAndInjuries(precomp.address).should.be.fulfilled;
        tactics0 = await engine.encodeTactics(substitutions, lineup0, extraAttackNull, tacticId442).should.be.fulfilled;
        tactics1 = await engine.encodeTactics(substitutions, lineup1, extraAttackNull, tacticId433).should.be.fulfilled;
        tactics1NoChanges = await engine.encodeTactics(noSubstitutions, lineup1, extraAttackNull, tacticId433).should.be.fulfilled;
        tactics442 = await engine.encodeTactics(substitutions, lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        tactics442NoChanges = await engine.encodeTactics(noSubstitutions, lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2).should.be.fulfilled;
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2).should.be.fulfilled;
        MAX_RND = await engine.MAX_RND().should.be.fulfilled;
        MAX_RND = MAX_RND.toNumber();
        kMaxRndNumHalf = Math.floor(MAX_RND/2)-200; 
        events1Half = Array.from(new Array(7), (x,i) => 0);
        events1Half = [events1Half,events1Half];
    });
    
    it('play a match to estimate cost', async () => {
        const result = await engine.playMatchWithCost(seed, now, [teamStateAll50, teamStateAll1], [tactics0, tactics1], firstHalfLog, matchBools).should.be.fulfilled;
    });
return
    it('penaltyPerAge', async () => {
        ageInDays       = [31*365, 31*365+1, 31*365+2, 41*365-4, 41*365-3, 41*365-2, 41*365-1];
        expectedPenalty = [1000000, 998904, 998904, 1544, 1544, 0, 0, 0]
        for (i = 0; i < ageInDays.length; i++) {
            dayOfBirth = Math.round(secsToDays(now) - ageInDays[i]/7);
            playerSkills = await engine.encodePlayerSkills(
                skills = [1,1,1,1,1], 
                dayOfBirth, 
                playerId = 2132321,
                [potential = 3,
                forwardness,
                leftishness,
                aggr = 0],
                alignedLastHalf = true,
                redCardLastGame = false,
                gamesNonStopping = 0,
                injuryWeeksLeft = 0
            ).should.be.fulfilled;
            result = await engine.penaltyPerAge(playerSkills, now).should.be.fulfilled;
            result.toNumber().should.be.equal(expectedPenalty[i]);
        }
    });
    

    it('play a match with penalties to estimate cost', async () => {
        const result = await engine.playMatchWithCost(seed, now, [teamStateAll50, teamStateAll1], [tactics0, tactics1], firstHalfLog, [is2nd = true, isHomeStadium,  playoff = true]).should.be.fulfilled;
    });

    it('check that penalties are played in playoff games', async () => {
        // this game ends up in a tie if there are no penalties:
        seedDraw = seed + 2;
        log0 = await engine.playMatch(seedDraw, now, [teamStateAll50, teamStateAll50], [tactics442, tactics1], log = [0, 0], [is2nd = false, isHomeStadium, playoff = false]).should.be.fulfilled;
        log12 = await engine.playMatch(seedDraw, now, [teamStateAll50, teamStateAll50], [tactics442, tactics1], log0, [is2nd = true, isHomeStadium,  playoff = false]).should.be.fulfilled;
        // check that the game would end 2-2
        // and that there are no penalties
        expectedResult = [0, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log12[team]);
            decodedLog[0].toNumber().should.be.equal(expectedResult[team]);
        }
        // and that there are no penalties
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log12[team]);
            for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(false);
        }
        // now play the game in 'playoff mode'
        log12 = await engine.playMatch(seedDraw,  now,[teamStateAll50, teamStateAll50], [tactics442, tactics1], log0, [is2nd = true, isHomeStadium,  playoff = true]).should.be.fulfilled;
        expected = [true, true, true, true, true, true, false]
        decodedLog = await encodingLog.decodeMatchLog(log12[team = 0]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
        expected = [true, true, true, true, true, true, true]
        decodedLog = await encodingLog.decodeMatchLog(log12[team = 1]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
    });

    it('computePenalties', async () => {
        // one team much better than the other:
        log = await precomp.computePenalties(log = [0,0], [teamStateAll50, teamStateAll1], 50, 1, seed);
        expected = [true, true, true, true, true, false, false]
        decodedLog = await encodingLog.decodeMatchLog(log[team = 0]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
        expected = [false, false, false, false, false, false, false]
        decodedLog = await encodingLog.decodeMatchLog(log[team = 1]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
        
        // both teams similar:
        log = await precomp.computePenalties(log = [0,0], [teamStateAll50, teamStateAll50], 50, 50, seed);
        expected = [false, true, true, true, true, false, false]
        decodedLog = await encodingLog.decodeMatchLog(log[team = 0]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
        expected = [true, true, true, true, true, false, false]
        decodedLog = await encodingLog.decodeMatchLog(log[team = 1]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);

        // both teams really incredible goalkeepers:
        log = await precomp.computePenalties(log = [0,0], [teamStateAll50, teamStateAll50], 5000000, 5000000, seed);
        expected = [false, false, false, false, false, false, true]
        decodedLog = await encodingLog.decodeMatchLog(log[team = 0]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
        expected = [false, false, false, false, false, false, false]
        decodedLog = await encodingLog.decodeMatchLog(log[team = 1]);
        for (i = 0; i < 7; i++) decodedLog.penalties[i].should.be.equal(expected[i]);
    });

    it('goals from 1st half are added in the 2nd half', async () => {
        seedDraw = 12;
        log0 =  await engine.playMatch(seedDraw,  now, [teamStateAll50, teamStateAll50], [tactics442, tactics1], log = [0, 0], [is2nd = false, isHomeStadium, isPlayoff]).should.be.fulfilled;
        log1 =  await engine.playMatch(seedDraw,  now, [teamStateAll50, teamStateAll50], [tactics442, tactics1], log = [0, 0], [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        log12 = await engine.playMatch(seedDraw,  now, [teamStateAll50, teamStateAll50], [tactics442, tactics1], log0, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // for this seedDraw, they all score one goal in each half
        expected = [1, 1]
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log0[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expected[team]);
        }
        // so the final result should be 2-2
        expected = [2, 2]
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log12[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expected[team]);
        }
    });
    
    it('play 2nd half with 3 changes is OK, but more than 3 is rejected, by lying in the team-states', async () => {
        // create a 2nd half using 3 players that already played in the 1st half... should work
        messi = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0).should.be.fulfilled;            
        for (p = 0; p < 3; p++) teamStateAll50[p] = messi; 
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442NoChanges, tactics1NoChanges], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // create a 2nd half using 4 players that already played in the 1st half... should fail
        teamStateAll50[5] = messi; 
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442NoChanges, tactics1NoChanges], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.rejected;
    });

    it('play 2nd half with 3 changes is OK, but more than 3 is rejected, by lying in the substitutions', async () => {
        // create a 2nd half using 1 players that already played in the 1st half, and 2 changes only... should work
        messi = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0).should.be.fulfilled;            
        teamStateAll50[p] = messi; 
        tactics442TwoChanges = await engine.encodeTactics([3,1,11], lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442TwoChanges, tactics1NoChanges], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        tactics442ThreeChanges = await engine.encodeTactics([3,1,5], lineupConsecutive, extraAttackNull, tacticId442).should.be.fulfilled;
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442ThreeChanges, tactics1NoChanges], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.rejected;
    });

    it('play with an injured / red carded / free-slot player', async () => {
        // legit works:
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442, tactics1], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.fulfilled;
        // red card fails:
        teamStateAll50[5] = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0],
            alignedLastHalf = false, redCardLastGame = true, gamesNonStopping = 0, injuryWeeksLeft = 0).should.be.fulfilled;            
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442, tactics1], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.rejected;
        // injured fails
        teamStateAll50[5] = await engine.encodePlayerSkills([50,50,50,50,50], dayOfBirth21, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0],
            alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 2).should.be.fulfilled;            
        result = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics442, tactics1], firstHalfLog, [is2nd = true, isHomeStadium, isPlayoff]).should.be.rejected;
    });

    it('computePenaltyBadPositionAndCondition for GK ', async () => {
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 232131, [potential = 1,
            forwardness = 0, leftishness = 0, aggr = 0], alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0
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
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 312321, [potential = 1,
            forwardness = 1, leftishness = 4, aggr = 0], alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0
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
            playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 1323121, [potential = 1,
                forwardness = 1, leftishness = 4, aggr = 0], alignedLastHalf = false, redCardLastGame = false, games, injuryWeeksLeft = 0
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
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 312321, [potential = 1,
            forwardness = 5, leftishness = 7, aggr = 0], alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0
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
        log = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHome = true, isPlayoff]).should.be.fulfilled;
        expected = [10, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expected[team]);
        }
    });
    
    it('play a match', async () => {
        log = await engine.playMatch(seed, now, [teamStateAll50, teamStateAll1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expected = [10, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expected[team]);
        }
    });
    
    it('manages to score with select shooter without modifiers', async () => {
        // lets put a Messi and check that it surely scores:
        teamState = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        messi = await engine.encodePlayerSkills([100,100,100,100,100], dayOfBirth21, id = 1123, [pot = 3, fwd = 3, left = 7, aggr = 0], 
            alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0).should.be.fulfilled;            
        teamState[10] = messi;
        result = await engine.selectShooter(now, teamState, playersPerZone442, extraAttackNull, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(10);
        teamThatAttacks = 0;
        log = await engine.managesToScore(now, log = [0,0], teamThatAttacks, teamState, playersPerZone442, extraAttackNull, blockShoot = 1, [kMaxRndNumHalf, kMaxRndNumHalf, kMaxRndNumHalf]).should.be.fulfilled;
        // for this case, there should be a goal, so: 1-0    
        expectedGoals       = [1, 0];
        expectedShooters    = [10, 0];
        expectedAssisters   = [10, 0];
        expectedFwd         = [3, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expectedGoals[team]);
            decodedLog.shootersIdx[0].toNumber().should.be.equal(expectedShooters[team]);
            decodedLog.assistersIdx[0].toNumber().should.be.equal(expectedAssisters[team]);
            decodedLog.shooterForwardPos[0].toNumber().should.be.equal(expectedFwd[team]);
        }
        // let's put a radically good GK, and check that it doesn't score
        log = await engine.managesToScore(now, log = [0,0], teamThatAttacks, teamState, playersPerZone442, extraAttackNull, blockShoot = 1000, [kMaxRndNumHalf, kMaxRndNumHalf, kMaxRndNumHalf]).should.be.fulfilled;
        expectedGoals       = [0, 0];
        expectedShooters    = [0, 0];
        expectedAssisters   = [0, 0];
        expectedFwd         = [0, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expectedGoals[team]);
            decodedLog.shootersIdx[0].toNumber().should.be.equal(expectedShooters[team]);
            decodedLog.assistersIdx[0].toNumber().should.be.equal(expectedAssisters[team]);
            decodedLog.shooterForwardPos[0].toNumber().should.be.equal(expectedFwd[team]);
        }
        // Finally, check that even with a super-goalkeeper, there are chances of scoring (e.g. if the rnd is super small, in this case)
        log = await engine.managesToScore(now, log = [0,0], teamThatAttacks, teamState, playersPerZone442, extraAttackNull, blockShoot = 1000, [kMaxRndNumHalf, 1, kMaxRndNumHalf]).should.be.fulfilled;
        expectedGoals       = [1, 0];
        expectedShooters    = [10, 0];
        expectedAssisters   = [10, 0];
        expectedFwd         = [3, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(log[team]);
            decodedLog.nGoals.toNumber().should.be.equal(expectedGoals[team]);
            decodedLog.shootersIdx[0].toNumber().should.be.equal(expectedShooters[team]);
            decodedLog.assistersIdx[0].toNumber().should.be.equal(expectedAssisters[team]);
            decodedLog.shooterForwardPos[0].toNumber().should.be.equal(expectedFwd[team]);
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
        messi = await engine.encodePlayerSkills([2,2,2,2,2], dayOfBirth21, id = 1323121, [pot = 3, fwd = 3, left = 7, aggr = 0],
            alignedLastHalf = false, redCardLastGame = false, gamesNonStopping = 0, injuryWeeksLeft = 0).should.be.fulfilled;            
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

    it('getLineUpAndPlayerPerZone for wrong tactics', async () => {
        tacticsWrong = await engine.encodeTactics(substitutions, lineup1, extraAttackNull, tacticIdTooLarge = 6).should.be.fulfilled;
        result = await engine.getLineUpAndPlayerPerZone(tacticsWrong, tactics1, is2ndHalf).should.be.rejected;
    });

    it('getLineUpAndPlayerPerZone', async () => {
        teamState442 = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        result = await engine.getLineUpAndPlayerPerZone(teamState442, tactics1, is2ndHalf).should.be.fulfilled;
        let {0: states, 1:fwdMods , 2: playersPerZone} = result;
        for (p = 0; p < 6; p++) playersPerZone[p].toNumber().should.be.equal(playersPerZone433[p]);
        for (p = 0; p < 11; p++) states[p].should.be.bignumber.equal(teamState442[lineup1[p]]);
    });

    it('play match with wrong tactic', async () => {
        tacticsWrong = await engine.encodeTactics(substitutions, lineup1, extraAttackNull, tacticIdTooLarge = 6);
        await engine.playMatch(seed, now, teamStateAll50, teamStateAll50, [tacticsWrong, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.rejected;
    });


    it('different team state => different result', async () => {
        matchLog = await engine.playMatch(123456, now, [teamStateAll50, teamStateAll50], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [2, 2];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(matchLog[team]);
            decodedLog[0].toNumber().should.be.equal(expectedResult[team]);
        }
        matchLog = await engine.playMatch(123456, now, [teamStateAll50, teamStateAll1], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [11, 0];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(matchLog[team]);
            decodedLog[0].toNumber().should.be.equal(expectedResult[team]);
        }
    });

    it('different seeds => different result', async () => {
        matchLog = await engine.playMatch(123456, now, [teamStateAll50, teamStateAll50], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [2, 2];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(matchLog[team]);
            decodedLog[0].toNumber().should.be.equal(expectedResult[team]);
        }
        matchLog = await engine.playMatch(654322, now, [teamStateAll50, teamStateAll50], [tactics0, tactics1], firstHalfLog, [is2ndHalf, isHomeStadium, isPlayoff]).should.be.fulfilled;
        expectedResult = [1, 1];
        for (team = 0; team < 2; team++) {
            decodedLog = await encodingLog.decodeMatchLog(matchLog[team]);
            decodedLog[0].toNumber().should.be.equal(expectedResult[team]);
        }
    });
});