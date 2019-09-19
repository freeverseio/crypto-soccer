const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Engine = artifacts.require('Engine');
const Assets = artifacts.require('Assets');

contract('Engine', (accounts) => {
    // const seed = 610106;
    const seed = web3.utils.toBN(web3.utils.keccak256("32123"));
    const lineup0 = [0, 3, 4, 5, 6, 9, 10, 11, 12, 15, 16];
    const lineup1 = [0, 3, 4, 5, 6, 9, 10, 11, 16, 17, 18];
    const lineupConsecutive =  Array.from(new Array(11), (x,i) => i);
    const tacticId442 = 0; // 442
    const tacticId433 = 2; // 433
    const playersPerZone442 = [1,2,1,1,2,1,0,2,0];
    const playersPerZone433 = [1,2,1,1,1,1,1,1,1];
    const PLAYERS_PER_TEAM_MAX = 25;
    const IDX_R = 1;
    const IDX_C = 2;
    const IDX_CR = 3;
    const IDX_L = 4;
    const IDX_LR = 5;
    const IDX_LC = 6;
    const IDX_LCR = 7;

    
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
                playerId = 1,
                potential,
                forwardness,
                leftishness
            ).should.be.fulfilled;            
            teamState.push(playerSkillsTemp)
        }        
        return teamState;
    };

    const createTeamState442 = async (engine, forceSkills) => {
        teamState = [];
        month = 0;
        playerId = 1;
        pot = 3;
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 0, left = 0).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 1, left = IDX_L).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 1, left = IDX_C).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 1, left = IDX_C).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 1, left = IDX_R).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 2, left = IDX_L).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 2, left = IDX_C).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 2, left = IDX_C).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 2, left = IDX_R).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 3, left = IDX_C).should.be.fulfilled 
        teamState.push(pSkills)
        pSkills = await engine.encodePlayerSkills(forceSkills, month, playerId, pot, fwd = 3, left = IDX_C).should.be.fulfilled 
        teamState.push(pSkills)
        for (p = 11; p < PLAYERS_PER_TEAM_MAX; p++) {
            teamState.push(pSkills)
        }        
        return teamState;
    };


    const createTeamStateFromSinglePlayer = async (skills, engine, forwardness = 3, leftishness = 2) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1,
            potential = 3,
            forwardness,
            leftishness
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
        tactics0 = await engine.encodeTactics(lineup0, tacticId442).should.be.fulfilled;
        tactics1 = await engine.encodeTactics(lineup1, tacticId433).should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2).should.be.fulfilled;
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2).should.be.fulfilled;
        MAX_PENALTY = await engine.MAX_PENALTY().should.be.fulfilled;
        MAX_PENALTY = MAX_PENALTY.toNumber();
    });

    // it('play a match to estimate cost', async () => {
    //     const result = await engine.playMatchWithCost(seed, teamStateAll50, teamStateAll1, [tactics0, tactics1]).should.be.fulfilled;
    // });
    // return;

    it('computePenalty for GK ', async () => {
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 1, potential = 1,
            forwardness = 0, leftishness = 0
        ).should.be.fulfilled;            
        expected = Array.from(new Array(11), (x,i) => MAX_PENALTY);
        expected[0] = 0;
        for (p=0; p < 11; p++) {
            penalty = await engine.computePenalty(p, playersPerZone442, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected[p]);
        }
    });

    it('computePenalty for DL ', async () => {
            // for a DL:
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 1, potential = 1,
            forwardness = 1, leftishness = 4
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
            penalty = await engine.computePenalty(p, playersPerZone442, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected442[p]);
            penalty = await engine.computePenalty(p, playersPerZone433, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected433[p]);
        }
    });

    it('computePenalty for MFLCR ', async () => {
        // for a DL:
        playerSkills= await engine.encodePlayerSkills(skills = [1,1,1,1,1], monthOfBirth = 0,  playerId = 1, potential = 1,
            forwardness = 5, leftishness = 7
        ).should.be.fulfilled;            
        expected442 = [MAX_PENALTY, 
            1000, 1000, 1000, 1000, 
            0, 0, 0, 0, 
            0, 0 
        ];
        expected433 = expected442;
        for (p=0; p < 11; p++) {
            penalty = await engine.computePenalty(p, playersPerZone442, playerSkills).should.be.fulfilled;
            penalty.toNumber().should.be.equal(10000 - expected442[p]);
            penalty = await engine.computePenalty(p, playersPerZone433, playerSkills).should.be.fulfilled;
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
    

    it('play a match', async () => {
        const result = await engine.playMatch(seed, [teamStateAll50, teamStateAll1], [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(17);
        result[1].toNumber().should.be.equal(0);
    });

    it('select shooter and manages to score', async () => {
        // interface: 
        // managesToScore(teamState, playersPerZone, lineup, blockShoot, rndNum1,rndNum2)
        teamState = await createTeamState442(engine, forceSkills= [1,1,1,1,1]).should.be.fulfilled;
        messi = await engine.encodePlayerSkills([100,100,100,100,100], month = 0, id = 1, pot = 3, fwd = 3, left = 7).should.be.fulfilled;            
        teamState[10] = messi;
        kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        result = await engine.selectShooter(teamState, playersPerZone442, lineupConsecutive, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(10);
        result = await engine.managesToScore(teamState, playersPerZone442, lineupConsecutive, blockShoot = 1, kMaxRndNumHalf, kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
        result = await engine.managesToScore(teamState, playersPerZone442, lineupConsecutive, blockShoot = 1000, kMaxRndNumHalf, kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        // even with a super-goalkeeper, there are chances of scoring (e.g. if the rnd is super small, in this case)
        kMaxRndNumHalf = 1;
        result = await engine.managesToScore(teamState, playersPerZone442, lineupConsecutive, blockShoot = 1000, 8000, 1).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('throws dice array11', async () => {
        // interface: throwDiceArray(uint[11] memory weights, uint rndNum)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        weights = Array.from(new Array(11), (x,i) => 1);
        weights[8] = 1000;
        let result = await engine.throwDiceArray11(weights, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(8);
        weights[8] = 1;
        weights[9] = 1000;
        result = await engine.throwDiceArray11(weights, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(9);
        weights[9] = 1;
        weights[10] = 1000;
        result = await engine.throwDiceArray11(weights, kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(10);
    });

    it('manages to shoot', async () => {
        // interface: managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
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
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
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
        nRnds = 235;
        const result = await engine.getNRandsFromSeed(nRnds, seed).should.be.fulfilled;
        result.length.should.be.equal(nRnds);
        result[0].should.be.bignumber.equal("2433");
        result[nRnds-1].should.be.bignumber.equal("7039");
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
        result = await engine.getTeamGlobSkills(teamState442, playersPerZone442, lineupConsecutive).should.be.fulfilled;
        result.attackersSpeed.length.should.be.equal(2);
        result.attackersShoot.length.should.be.equal(2);
        result.attackersSpeed[0].should.be.bignumber.equal("1");
        result.attackersSpeed[1].should.be.bignumber.equal("1");
        result.attackersShoot[0].should.be.bignumber.equal("1");
        result.attackersShoot[1].should.be.bignumber.equal("1");
        expectedGlob = [42, 4, 8, 1, 70];
        for (g = 0; g < 5; g++) result.globSkills[g].toNumber().should.be.equal(expectedGlob[g]);
    });
    
    it('getLineUpAndPlayerPerZone for wrong tactics', async () => {
        tacticsWrong = await engine.encodeTactics(lineup1, tacticIdTooLarge = 6).should.be.fulfilled;
        result = await engine.getLineUpAndPlayerPerZone(tacticsWrong).should.be.rejected;
    });

    it('getLineUpAndPlayerPerZone', async () => {
        result = await engine.getLineUpAndPlayerPerZone(tactics1).should.be.fulfilled;
        let {0: line, 1: playersPerZone} = result;
        for (p = 0; p < 6; p++) playersPerZone[p].toNumber().should.be.equal(playersPerZone433[p]);
        for (p = 0; p < 11; p++) line[p].toNumber().should.be.equal(lineup1[p]);
    });

    it('play match with wrong tactic', async () => {
        tacticsWrong = await engine.encodeTactics(lineup1, tacticIdTooLarge = 6);
        await engine.playMatch(seed, teamStateAll50, teamStateAll50, [tacticsWrong, tactics1]).should.be.rejected;
    });


    it('different team state => different result', async () => {
        let result = await engine.playMatch(123456, [teamStateAll50, teamStateAll50], [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(2);
        result = await engine.playMatch(123456, [teamStateAll50, teamStateAll1], [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(12);
        result[1].toNumber().should.be.equal(0);
    });

    it('different seeds => different result', async () => {
        let result = await engine.playMatch(123456, [teamStateAll50, teamStateAll50], [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(2);
        result = await engine.playMatch(654321, [teamStateAll50, teamStateAll50], [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(1);
    });
});