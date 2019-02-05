const keccak = require('keccak');
require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let engine = null;
    const initPlayerState = [
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
    ];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
    });

    it('play a match', async () => {
        const seed = keccak('keccak256').update('Hello World!').digest('hex');
        const stateTeam0 = initPlayerState[0];
        const stateTeam1 = initPlayerState[1];
        const tacticTeam0 = [4, 4, 2];
        const tacticTeam1 = [4, 3, 3];
        const result = await engine.playMatch(seed, stateTeam0, stateTeam1, tacticTeam0, tacticTeam1).should.be.fulfilled;
        result[0].toNumber().should.be.equal(3);
        result[1].toNumber().should.be.equal(2);
    });
});