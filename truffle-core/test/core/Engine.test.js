require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let engine = null;
    const seed = 400342352351;
    const state0 = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15];
    const state1 = [11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1];
    const tactic0 = [4, 4, 3];
    const tactic1 = [4, 5, 2];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
    });

    it('play a match', async () => {
        const result = await engine.playMatch(seed, state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(2);
        result[1].toNumber().should.be.equal(0);
    });

    it('play match with less than 11 players', async () => {
        const wrongTeam = [0,1,2,3,4,5,6,7,8,9];
        await engine.playMatch(seed, wrongTeam, state1, tactic0, tactic1).should.be.rejected;
        await engine.playMatch(seed, state0, wrongTeam, tactic0, tactic1).should.be.rejected;
    });

    it('play match with wrong tactic', async () => {
        await engine.playMatch(seed, state0, state1, [4,4,2], tactic1).should.be.rejected;
        await engine.playMatch(seed, state0, state1, [4,4,4], tactic1).should.be.rejected;
        await engine.playMatch(seed, state0, state1, tactic0, [4,4,2]).should.be.rejected;
        await engine.playMatch(seed, state0, state1, tactic0, [4,4,4]).should.be.rejected;
    });

    it('different seeds => different result', async () => {
        let result = await engine.playMatch('seed', state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(0);
        result[1].toNumber().should.be.equal(1);
        result = await engine.playMatch('different seed', state0, state1, tactic0, tactic1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(0);
        result[1].toNumber().should.be.equal(0);
    })
});