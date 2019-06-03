require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let engine = null;
    const seed = '0x610106';
    const state0 = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    const state1 = state0;
    console.log(state1);
    const tactic0 = [4, 4, 2];
    const tactic1 = [4, 5, 1];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
    });

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

    it('computes team global skills by aggregating across all players in team', async () => {
        let result = await engine.getTeamGlobSkills(state0, tactic0).should.be.fulfilled;
        // result.home.toNumber().should.be.equal(2);
    });
});