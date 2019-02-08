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
    const tactics = [
        [4,4,3],  // Team 0
        [5,4,2]   // Team 1
    ];
    const blocksToInit = 3;
    const step = 1
    const teamIds = [1, 2];

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address).should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
    });

    it('Engine contract', async () => {
        const address = await leagues.getEngineContract().should.be.fulfilled;
        address.should.be.equal(engine.address);
    });

    it('compute unexistent league', async () => {
        const id = 532;
        await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.rejected;
    });

    it('compute league', async () => {
        const scores = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;
        scores.length.should.be.equal(teamIds.length * (teamIds.length - 1));
    });

    it('check hashing of the result', async () => {
        const score = [
            [3,1]
        ];
        const finalHash = await leagues.calculateFinalHash(score).should.be.fulfilled;
    });

    it('compute league 2 times gives the same result', async () => {
        const scores0 = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;       
        const scores1 = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;   
        const finalHash0 = await leagues.calculateFinalHash(scores0).should.be.fulfilled;
        const finalHash1 = await leagues.calculateFinalHash(scores1).should.be.fulfilled;
        finalHash0.should.be.equal(finalHash1);
    });

    it('compute league and update changed final hash of the league', async () => {
        const before = await leagues.getHash(id).should.be.fulfilled;
        await leagues.computeLeagueAndUpdate(id, initPlayerState, tactics).should.be.fulfilled;
        const after = await leagues.getHash(id).should.be.fulfilled;
        after.should.not.be.equal(before);
    });

    it('external league final hash is equal to bc calculated one', async () => {
        const scores = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;       
        const finalHash = await leagues.calculateFinalHash(scores).should.be.fulfilled;
        await leagues.computeLeagueAndUpdate(id, initPlayerState, tactics).should.be.fulfilled;
        const bcFinalHash = await leagues.getHash(id).should.be.fulfilled;
        bcFinalHash.should.be.equal(finalHash);
    });
});