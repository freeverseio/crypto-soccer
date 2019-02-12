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
        1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, // Team 0
        0, // divider
        10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 11, 12, // Team 1
        0
    ];
    const tactics = [
        [4, 4, 3],  // Team 0
        [5, 4, 2]   // Team 1
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
        const hash0 = await leagues.calculateFinalHash([[0, 1]]).should.be.fulfilled;
        const hash1 = await leagues.calculateFinalHash([[0, 1], [2, 1]]).should.be.fulfilled;
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
        const hash0 = await leagues.hashTactics([[4, 4, 2]]).should.be.fulfilled;
        const hash1 = await leagues.hashTactics([[4, 4, 2]]).should.be.fulfilled;
        hash1.should.be.equal(hash0);
        const hash2 = await leagues.hashTactics([[3, 4, 2]]).should.be.fulfilled;
        hash2.should.be.not.equal(hash0);
        const hash3 = await leagues.hashTactics([[4, 5, 2]]).should.be.fulfilled;
        hash3.should.be.not.equal(hash0);
        const hash4 = await leagues.hashTactics([[4, 4, 3]]).should.be.fulfilled;
        hash4.should.be.not.equal(hash0);
        const hash5 = await leagues.hashTactics([[4, 4, 2], [4, 4, 2]]).should.be.fulfilled;
        hash5.should.be.not.equal(hash0);
    });

    it('hash team state', async () => {
        const state = [324, 435, 5];
        const hash0 = await leagues.hashTeamState(state).should.be.fulfilled;
        hash0.should.be.equal('0x2dcd8f033162e070f623608ce1f1a913bc979d6d070221b812a25fa27b78f86b');
        const hash1 = await leagues.hashTeamState([2, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('hash init state', async () => {
        const state = [3, 435, 5];
        const hash0 = await leagues.hashTeamState(state).should.be.fulfilled;
        hash0.should.be.equal('0x8c95de1a9b22dd1419122bfe86a58534751f629fd72d98bb03da9c4f1b24d420');
        const hash1 = await leagues.hashTeamState([2, 3, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('hash league state', async () => {
        const state = [3,5,2,0,4,56,6,0];
        const hashes = await leagues.hashLeagueState(state).should.be.fulfilled;
        hashes.length.should.be.equal(2);
        hashes[0].should.be.equal('0xc6951eb9cd3a570943a21ee0c2156cc75258037b093e3bb1690d6d92c9af8c29');
        hashes[1].should.be.equal('0xa4efe975f12e6f4dce1ee13b2355a8b2057cbd37cad109f6fb027c0a645ee92a');
    })

    it('update league state', async () => {
        const initStateHash = '0x435a354320000000000000000000000000000000000000000000000000000000';
        const finalTeamStateHashes = [
            '0x245964674ab00000000000000000000000000000000000000000000000000000',
            '0xaaaaa00000000000000000000000000000000000000000000000000000000000'
        ];
        const scores = [[1, 2], [2, 2]];
        await leagues.updateLeague(id, initStateHash, finalTeamStateHashes, scores).should.be.fulfilled;
        let result = await leagues.getInitHash(id).should.be.fulfilled;
        result.should.be.equal(initStateHash);
        result = await leagues.getFinalTeamStateHashes(id).should.be.fulfilled;
        result.length.should.be.equal(finalTeamStateHashes.length);
        result[0].should.be.equal(finalTeamStateHashes[0]);
        result[1].should.be.equal(finalTeamStateHashes[1]);
        result = await leagues.getScores(id).should.be.fulfilled;
        result.length.should.be.equal(scores.length);
        result[0][0].toNumber().should.be.equal(scores[0][0]);
        result[0][1].toNumber().should.be.equal(scores[0][1]);
        result[1][0].toNumber().should.be.equal(scores[1][0]);
        result[1][1].toNumber().should.be.equal(scores[1][1]);
    });

    it('count empty teams status', async () => {
        await leagues.countTeamsStatus([]).should.be.rejected;
    });

    it('teams status should not init with invalid status', async () => {
        await leagues.countTeamsStatus([0, 3, 4, 0]).should.be.rejected;
    });

    it('teams status should not end with invalid status', async () => {
        await leagues.countTeamsStatus([3, 4]).should.be.rejected;
    });

    it('count teams status', async () => {
        const result = await leagues.countTeamsStatus([3, 4, 0, 454, 0, 6, 0]).should.be.fulfilled;
        result.toNumber().should.be.equal(3);
    });

});