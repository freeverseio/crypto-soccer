const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    const seed = 610106;
    const lineup0 = [0, 3, 4, 5, 6, 9, 10, 11, 12, 15, 16];
    const lineup1 = [0, 3, 4, 5, 6, 9, 10, 11, 16, 17, 18];
    const tacticId0 = 0; // 442
    const tacticId1 = 2; // 433
    const playersPerZone0 = [1,2,1,2,0,2];
    const playersPerZone1 = [1,2,1,1,1,1];
    const PLAYERS_PER_TEAM_MAX = 25;

    const createTeamStateFromSinglePlayer = async (skills, engine) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1,
            potential = 3,
            forwardness = 3,
            leftishness = 2
        ).should.be.fulfilled;

        teamState = []
        for (player = 0; player < PLAYERS_PER_TEAM_MAX; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        tactics0 = await engine.encodeTactics(lineup0, tacticId0).should.be.fulfilled;
        tactics1 = await engine.encodeTactics(lineup1, tacticId1).should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine).should.be.fulfilled;
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine).should.be.fulfilled;
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

    // it('play a match to estimate cost', async () => {
    //     const result = await engine.playMatchWithCost(seed, teamStateAll50, teamStateAll1, tactic0, tactic1).should.be.fulfilled;
    // });

    it('play a match', async () => {
        const result = await engine.playMatch(seed, teamStateAll50, teamStateAll1, [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(17);
        result[1].toNumber().should.be.equal(0);
    });

    it('manages to score', async () => {
        // interface: 
        // managesToScore(uint8 nAttackers, uint[] attackersSpeed, uint[], attackersShoot, blockShoot, rndNum1,rndNum2)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        let attackersSpeed = [10,1,1];
        let attackersShoot = [10,1,1];
        let blockShoot     = 1;
        nAttackers         = attackersShoot.length;
        let result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
        blockShoot     = 1000;
        result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(false);
        // even with a super-goalkeeper, there are chances of scoring (e.g. if the rnd is super small, in this case)
        kMaxRndNumHalf = 1;
        result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('throws dice array', async () => {
        // interface: throwDiceArray(uint[] memory weights, uint rndNum)
        let kMaxRndNumHalf = 8000; // the max allowed random number is 16383, so this is about half of it
        let result = await engine.throwDiceArray([1000,1,1],kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDiceArray([1,1000,1],kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await engine.throwDiceArray([1,1,1000],kMaxRndNumHalf).should.be.fulfilled;
        result.toNumber().should.be.equal(2);
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
        result[0].should.be.bignumber.equal("6666");
        result[nRnds-1].should.be.bignumber.equal("5318");
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
        
        let result = await engine.getTeamGlobSkills(teamStateAll1, playersPerZone0, lineup0).should.be.fulfilled;
        result.attackersSpeed.length.should.be.equal(2);
        result.attackersShoot.length.should.be.equal(2);
        result.attackersSpeed[0].should.be.bignumber.equal("1");
        result.attackersSpeed[1].should.be.bignumber.equal("1");
        result.attackersShoot[0].should.be.bignumber.equal("1");
        result.attackersShoot[1].should.be.bignumber.equal("1");
        result.globSkills[0].should.be.bignumber.equal("42");
        result.globSkills[1].should.be.bignumber.equal("4");
        result.globSkills[2].should.be.bignumber.equal("8");
        result.globSkills[3].should.be.bignumber.equal("1");
        result.globSkills[4].should.be.bignumber.equal("70");
    });

    it('getLineUpAndPlayerPerZone for wrong tactics', async () => {
        tacticsWrong = await engine.encodeTactics(lineup1, tacticIdTooLarge = 6).should.be.fulfilled;
        result = await engine.getLineUpAndPlayerPerZone(tacticsWrong).should.be.rejected;
    });

    it('getLineUpAndPlayerPerZone', async () => {
        result = await engine.getLineUpAndPlayerPerZone(tactics1).should.be.fulfilled;
        let {0: line, 1: playersPerZone} = result;
        for (p = 0; p < 6; p++) playersPerZone[p].toNumber().should.be.equal(playersPerZone1[p]);
        for (p = 0; p < 11; p++) line[p].toNumber().should.be.equal(lineup1[p]);
    });

    it('play match with wrong tactic', async () => {
        tacticsWrong = await engine.encodeTactics(lineup1, tacticIdTooLarge = 6);
        await engine.playMatch(seed, teamStateAll50, teamStateAll50, [tacticsWrong, tactics1]).should.be.rejected;
    });


    it('different team state => different result', async () => {
        let result = await engine.playMatch(123456, teamStateAll50, teamStateAll50, [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(1);
        result = await engine.playMatch(123456, teamStateAll50, teamStateAll1, [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(14);
        result[1].toNumber().should.be.equal(0);
    });

    it('different seeds => different result', async () => {
        let result = await engine.playMatch(123456, teamStateAll50, teamStateAll50, [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(1);
        result = await engine.playMatch(654321, teamStateAll50, teamStateAll50, [tactics0, tactics1]).should.be.fulfilled;
        result[0].toNumber().should.be.equal(1);
        result[1].toNumber().should.be.equal(0);
    });
});