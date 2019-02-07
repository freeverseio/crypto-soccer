const keccak = require('keccak');
require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('LeaguesComputer');

contract('LeaguesComputer', (accounts) => {
    let leagues = null;
    let engine = null;
    const id = 0;
    const initPlayerState = [
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
    ];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address).should.be.fulfilled;
    });

    it('Engine contract', async () => {
        const address = await leagues.getEngineContract().should.be.fulfilled;
        address.should.be.equal(engine.address);
    });

    it('compute unexistent league', async () => {
        await leagues.computeLeagueFinalState(id, initPlayerState).should.be.rejected;
    });

    // it('compute league', async () => {
    //     const blocksToInit = 1;
    //     const step = 1;
    //     const teamIds = [1, 2];
    //     await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
    //     const scores = await engine.computeLeagueFinalState(id, initPlayerState).should.be.fulfilled;
    // });
});