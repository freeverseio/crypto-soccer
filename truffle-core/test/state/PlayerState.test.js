const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const PlayerState = artifacts.require('PlayerState');

/// TODO: evaluate to extract the skills part
contract('PlayerState', (accounts) => {
    let instance = null;

    beforeEach(async () => {
        instance = await PlayerState.new().should.be.fulfilled;
    });

    it('create player state', async () => {
        const state = await instance.playerStateCreate(
            defence = '16383',
            speed = '13',
            pass = '4',
            shoot = '56',
            endurance = '456',
            0, 
            playerId = '1',
            0, 0, 0, 0, 0, 0
        ).should.be.fulfilled;
        let result = await instance.getDefence(state).should.be.fulfilled;
        result.should.be.bignumber.that.equals(defence);
        result = await instance.getSpeed(state).should.be.fulfilled;
        result.should.be.bignumber.that.equals(speed);
        result = await instance.getPass(state).should.be.fulfilled;
        result.should.be.bignumber.that.equals(pass);
        result = await instance.getShoot(state).should.be.fulfilled;
        result.should.be.bignumber.that.equals(shoot);
        result = await instance.getEndurance(state).should.be.fulfilled;
        result.should.be.bignumber.that.equals(endurance);
    });

    it('player with all skills 0 is valid', async () => {
        const playerState = await instance.playerStateCreate(0,0,0,0,0,0,playerId = '1' ,0,0,0,0,0,1).should.be.fulfilled;
        const valid = await instance.isValidPlayerState(playerState).should.be.fulfilled;
        valid.should.be.equal(true);
    })

    it('is valid player state', async () => {
        let result = await instance.isValidPlayerState(0).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('skills getters from state player', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        const playerState = await instance.playerStateCreate(
            defence,
            speed,
            pass,
            shoot,
            endurance,
            monthOfBirthInUnixTime = 40,
            playerId = 41,
            currentTeamId = 42,
            currentShirtNum = 3,
            prevLeagueId = 44,
            prevTeamPosInLeague = 45,
            prevShirtNumInLeague = 6,
            lastSaleBlock = 47).should.be.fulfilled;
        let result = await instance.getDefence(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(defence);
        result = await instance.getSpeed(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(speed);
        result = await instance.getPass(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(pass);
        result = await instance.getShoot(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(shoot);
        result = await instance.getEndurance(playerState).should.be.fulfilled;
        result.toNumber().should.be.equal(endurance);
        result = await instance.getCurrentTeamId(playerState).should.be.fulfilled;
        result.should.be.bignumber.equal('42');
    });

    it('player state evolve', async () => {
        const playerState = await instance.playerStateCreate(
            defence = 3, 
            speed = 4, 
            pass = 6, 
            shoot = 11, 
            endurance = 9, 
            0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        const delta = 3;
        const updatedState = await instance.playerStateEvolve(playerState, delta).should.be.fulfilled;
        updatedState.should.be.bignumber.that.not.equals(playerState);
        let skill = await instance.getDefence(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(defence + delta);
        skill = await instance.getSpeed(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(speed + delta);
        skill = await instance.getPass(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(pass + delta);
        skill = await instance.getShoot(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(shoot + delta);
        skill = await instance.getEndurance(updatedState).should.be.fulfilled;
        skill.toNumber().should.be.equal(endurance + delta);
    });

    it('get skills', async () => {
        const playerState = await instance.playerStateCreate(
            defence = 0, 
            speed = 0, 
            pass = 0, 
            shoot = 0, 
            endurance = 1, 
            0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        const skills = await instance.getSkills(playerState).should.be.fulfilled;
        skills.toNumber().should.be.equal(1);
    });
});