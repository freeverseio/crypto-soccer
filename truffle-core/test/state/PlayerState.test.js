require('chai')
    .use(require('chai-as-promised'))
    .should();

const PlayerState = artifacts.require('PlayerState');

contract('PlayerState', (accounts) => {
    let instance = null;
    let TEAMSTATEDIVIDER = null;
    let LEAGUESTATEDIVIDER = null;

    beforeEach(async () => {
        instance = await PlayerState.new().should.be.fulfilled;
        TEAMSTATEDIVIDER = await instance.TEAMSTATEDIVIDER().should.be.fulfilled;
        LEAGUESTATEDIVIDER = await instance.LEAGUESTATEDIVIDER().should.be.fulfilled;
    });

    it('create player state', async () => {
        const defence = 3;
        const speed = 23;
        const pass = 2;
        const shoot = 21;
        const endurance = 10;
        const state = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
        (state.toNumber() & 0xff).should.be.equal(endurance);
        (state.toNumber() >> 8 & 0xff).should.be.equal(shoot);
        (state.toNumber() >> 8*2 & 0xff).should.be.equal(pass);
        (state.toNumber() >> 8*3 & 0xff).should.be.equal(speed);
        (state.toNumber() >> 8*4 & 0xff).should.be.equal(endurance);
    });

    it('is valid player state', async () => {
        let result = await instance.isValidPlayerState(TEAMSTATEDIVIDER).should.be.fulfilled;
        result.should.be.equal(false);
        result = await instance.isValidPlayerState(LEAGUESTATEDIVIDER).should.be.fulfilled;
        result.should.be.equal(false);
    });

    it('skills getters from state player', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        const playerState = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
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
    });

    it('player state evolve', async () => {
        const defence = 3;
        const speed = 4;
        const pass = 6;
        const shoot = 11;
        const endurance = 9;
        const playerState = await instance.playerStateCreate(defence, speed, pass, shoot, endurance).should.be.fulfilled;
        const delta = 3;
        const updatedState = await instance.playerStateEvolve(playerState, delta).should.be.fulfilled;
        updatedState.toNumber().should.not.be.equal(playerState.toNumber());
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
});