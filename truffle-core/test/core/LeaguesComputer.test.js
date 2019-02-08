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
    
    beforeEach(async () => {
        const blocksToInit = 0;
        const step = 1
        const teamIds = [1, 2];
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
        const nTeams = await leagues.countTeams(id).should.be.fulfilled;
        scores.length.should.be.equal(nTeams * (nTeams - 1));
    });

    it('compute league 2 times gives the same result', async () => {
        const scores0 = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;       
        const scores1 = await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.fulfilled;   
        const finalHash0 = await leagues.calculateFinalHash(scores0).should.be.fulfilled;
        const finalHash1 = await leagues.calculateFinalHash(scores1).should.be.fulfilled;
        finalHash0.should.be.equal(finalHash1);
    });

    it('hash differents results => different hashes', async () => {
        const hash0 = await leagues.calculateFinalHash([[0,1]]).should.be.fulfilled;
        const hash1 = await leagues.calculateFinalHash([[0,1],[2,1]]).should.be.fulfilled;
        hash0.should.be.not.equal(hash1);
    });

    it('try calculare a uninitied league', async () => {
        const id = 100;
        const blocksToInit = 10000; // future block
        const step = 1
        const teamIds = [1, 2];
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        await leagues.computeLeagueFinalState(id, initPlayerState, tactics).should.be.rejected;       
    });

    it('hash tactics', async () => {
        const hash0 = await leagues.hashTactics([[4,4,2]]).should.be.fulfilled;
        const hash1 = await leagues.hashTactics([[4,4,2]]).should.be.fulfilled;
        hash1.should.be.equal(hash0);
        const hash2 = await leagues.hashTactics([[3,4,2]]).should.be.fulfilled;
        hash2.should.be.not.equal(hash0);
        const hash3 = await leagues.hashTactics([[4,5,2]]).should.be.fulfilled;
        hash3.should.be.not.equal(hash0);
        const hash4 = await leagues.hashTactics([[4,4,3]]).should.be.fulfilled;
        hash4.should.be.not.equal(hash0);
        const hash5 = await leagues.hashTactics([[4,4,2],[4,4,2]]).should.be.fulfilled;
        hash5.should.be.not.equal(hash0);
    })
});