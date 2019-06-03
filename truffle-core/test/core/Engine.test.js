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
    const seed = '0x610106';
    const state0 = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    const state1 = state0;
    const tactic0 = [4, 4, 2];
    const tactic1 = [4, 5, 1];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        teamStateLib = await TeamStateLib.new().should.be.fulfilled;
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
    return;
    it('play a match', async () => {
        const result = await engine.playMatch(seed, state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(0);
        result[1].toNumber().should.be.equal(0);
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
        let result = await engine.playMatch('0x123456', state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(0);
        result[1].toNumber().should.be.equal(0);
        result = await engine.playMatch('0x654321', state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(1);
        result[1].toNumber().should.be.equal(1);
    });

    it('different team state => different result', async () => {
        const state = [11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1];
        let result = await engine.playMatch('0x123456', state, state, tactic0, tactic1).should.be.fulfilled;
        result.home.toNumber().should.be.equal(2);
        result.visitor.toNumber().should.be.equal(2);
        const state0 = [44, 12, 13, 44, 3, 66, 5, 5, 3, 2, 1];
        result = await engine.playMatch('0x123456', state, state0, tactic0, tactic1).should.be.fulfilled;
        result.home.toNumber().should.be.equal(2);
        result.visitor.toNumber().should.be.equal(2);
    });

});