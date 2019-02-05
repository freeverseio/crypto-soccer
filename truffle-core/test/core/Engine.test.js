const keccak = require('keccak');
require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let leagues = null;
    let engine = null;
    const id = 0;
    const initPlayerState = [
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
    ];

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        engine = await Engine.new(leagues.address).should.be.fulfilled;
    });

    it('Leagues contract', async () => {
        const address = await engine.getLeaguesContract().should.be.fulfilled;
        address.should.be.equal(leagues.address);
    });

    it('compute unexistent league', async () => {
        await engine.computeLeagueFinalState(id, initPlayerState).should.be.rejected;
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

    // it('compute league', async () => {
    //     const blocksToInit = 1;
    //     const step = 1;
    //     const teamIds = [1, 2];
    //     await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
    //     const scores = await engine.computeLeagueFinalState(id, initPlayerState).should.be.fulfilled;
    // });
});