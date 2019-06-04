const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Engine = artifacts.require('Engine');
const TeamStateLib = artifacts.require('TeamState');

contract('Engine', (accounts) => {
    let engine = null;
    let teamStateLib = null;
    let teamState = null;
    const seed = 610106;
    const state0 = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    const state1 = state0;
    const tactic0 = [4, 4, 2];
    const tactic1 = [4, 5, 1];
    const nPlayers = 11;

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        teamStateLib = await TeamStateLib.new().should.be.fulfilled;
        const playerState = await teamStateLib.playerStateCreate(
            defence = '1',
            speed = '1',
            pass = '1',
            shoot = '1',
            endurance = '1',
            0, 
            playerId = '1',
            0, 0, 0, 0, 0, 0
        ).should.be.fulfilled;
        teamState = await teamStateLib.teamStateCreate().should.be.fulfilled;
        for (var i = 0; i < nPlayers; i++) {
            teamState = await teamStateLib.teamStateAppend(teamState, playerState).should.be.fulfilled;
        }
    });

    it('manages to score', async () => {
        // interface: 
        // managesToScore(uint8 nAttackers, uint[] attackersSpeed, uint[], attackersShoot, blockShoot, rndNum1,rndNum2, uint maxRndNum)
        const kMaxRndNum = 16383; // 16383 = 2^kBitsPerRndNum-1 
        let kMaxRndNumHalf = 8000;
        let attackersSpeed = [10,1,1];
        let attackersShoot = [10,1,1];
        let blockShoot     = 1;
        nAttackers         = attackersShoot.length;
        let result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(true);
        blockShoot     = 1000;
        result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(false);
        // even with a super-goalkeeper, there are chances of scoring (e.g. if the rnd is super small, in this case)
        kMaxRndNumHalf = 1;
        result = await engine.managesToScore(nAttackers,attackersSpeed,attackersShoot,blockShoot,kMaxRndNumHalf,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(true);
    });


    it('throws dice array', async () => {
        // interface: throwDiceArray(uint[] memory weights, uint rndNum, uint maxRndNum)
        const kMaxRndNum = 16383; // 16383 = 2^kBitsPerRndNum-1 
        const kMaxRndNumHalf = 8000;
        let result = await engine.throwDiceArray([1000,1,1],kMaxRndNumHalf, kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDiceArray([1,1000,1],kMaxRndNumHalf, kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await engine.throwDiceArray([1,1,1000],kMaxRndNumHalf, kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(2);
    });

    it('manages to shoot', async () => {
        // interface: managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum, uint kMaxRndNum)
        const kMaxRndNum = 16383; // 16383 = 2^kBitsPerRndNum-1 
        const kMaxRndNumHalf = 8000;
        let globSkills = [[100,100,100,100,100], [1,1,1,1,1]];
        let result = await engine.managesToShoot(0,globSkills,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(true);
        result = await engine.managesToShoot(1,globSkills,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(false);
        globSkills = [[1,1,1,1,1], [100,100,100,100,100]];
        result = await engine.managesToShoot(0,globSkills,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(false);
        result = await engine.managesToShoot(1,globSkills,kMaxRndNumHalf,kMaxRndNum).should.be.fulfilled;
        result.should.be.equal(true);
    });

    it('throws dice', async () => {
        // interface: throwDice(uint weight1, uint weight2, uint rndNum, uint maxRndNum)
        const kMaxRndNum = 16383; // 16383 = 2^kBitsPerRndNum-1 
        let result = await engine.throwDice(0,10,2,kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await engine.throwDice(10,0,2,kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDice(10,10,2,kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await engine.throwDice(10,10,16000,kMaxRndNum).should.be.fulfilled;
        result.toNumber().should.be.equal(1);
    });

    return;
    it('play a match', async () => {
        const result = await engine.playMatch(seed, teamState, teamState, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(0);
        result[1].toNumber().should.be.equal(0);
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
        const playerState = await teamStateLib.playerStateCreate(
            defence = '1',
            speed = '1',
            pass = '1',
            shoot = '1',
            endurance = '1',
            0, 
            playerId = '1',
            0, 0, 0, 0, 0, 0
        ).should.be.fulfilled;
        const nPlayers = 11;
        let teamState = await teamStateLib.teamStateCreate().should.be.fulfilled;
        for (var i = 0; i < nPlayers; i++) {
            teamState = await teamStateLib.teamStateAppend(teamState, playerState).should.be.fulfilled;
        }
        let result = await engine.getTeamGlobSkills(teamState, [4,4,2]).should.be.fulfilled;
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
        console.log(result);
        // result.home.toNumber().should.be.equal(2);
    });

    it('play match with less than 11 players', async () => {
        const wrongTeam = [0,1,2,3,4,5,6,7,8,9];
        await engine.playMatch(seed, wrongTeam, state1, tactic0, tactic1).should.be.rejected;
        await engine.playMatch(seed, state0, wrongTeam, tactic0, tactic1).should.be.rejected;
    });

    it('play match with wrong tactic', async () => {
        await engine.playMatch(seed, state0, state1, [4,4,1], tactic1).should.be.rejected;
        await engine.playMatch(seed, state0, state1, [4,4,3], tactic1).should.be.rejected;
        await engine.playMatch(seed, state0, state1, tactic0, [4,4,1]).should.be.rejected;
        await engine.playMatch(seed, state0, state1, tactic0, [4,4,3]).should.be.rejected;
    });

    it('different seeds => different result', async () => {
        let result = await engine.playMatch(123456, state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(2);
        result = await engine.playMatch(654321, state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(3);
        result[1].toNumber().should.be.equal(3);
    });

    it('different team state => different result', async () => {
        const state = [11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1];
        let result = await engine.playMatch(123456, state, state, tactic0, tactic1).should.be.fulfilled;
        result.home.toNumber().should.be.equal(1);
        result.visitor.toNumber().should.be.equal(1);
        const state0 = [44, 12, 13, 44, 3, 66, 5, 5, 3, 2, 1];
        result = await engine.playMatch(123456, state, state0, tactic0, tactic1).should.be.fulfilled;
        result.home.toNumber().should.be.equal(1);
        result.visitor.toNumber().should.be.equal(0);
    });

});