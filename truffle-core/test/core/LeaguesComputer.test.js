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
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, // Team 0
        10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0  // Team 1
    ];
    const blocksToInit = 3;
    const step = 1
    const teamIds = [1, 2];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address).should.be.fulfilled;
    });

    it('Engine contract', async () => {
        const address = await leagues.getEngineContract().should.be.fulfilled;
        address.should.be.equal(engine.address);
    });

    it('compute unexistent league', async () => {
        const id = 532;
        await leagues.computeLeagueFinalState(id).should.be.rejected;
    });

    it('compute league', async () => {
        const tactics = [[4,4,3], [5,4,2]];
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const scores = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;
        scores.length.should.be.equal(teamIds.length * (teamIds.length - 1));
        scores[0][0].toNumber().should.be.equal(1);
        scores[0][1].toNumber().should.be.equal(0);
        scores[1][0].toNumber().should.be.equal(1);
        scores[1][1].toNumber().should.be.equal(0);
    });
});