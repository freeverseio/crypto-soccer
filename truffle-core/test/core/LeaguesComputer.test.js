require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const LeagueState = artifacts.require('LeagueState');
const Leagues = artifacts.require('LeaguesComputer');

contract('LeaguesComputer', (accounts) => {
    let leagues = null;
    let engine = null;

    const id = 0;
    const initPlayerState = [
        1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, // Team 0
        0, // divider
        10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 11, 12 // Team 1
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
        const leagueState = await LeagueState.new().should.be.fulfilled;
        Leagues.link("LeagueState", leagueState.address);
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
        scores0.length.should.be.equal(2);
        scores1.length.should.be.equal(2);
        scores0[0][0].toNumber().should.be.equal(scores1[0][0].toNumber());
        scores0[0][1].toNumber().should.be.equal(scores1[0][1].toNumber());
        scores0[1][0].toNumber().should.be.equal(scores1[1][0].toNumber());
        scores0[1][1].toNumber().should.be.equal(scores1[1][1].toNumber());
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
        hash0.should.be.equal('0xd2bd5ff7bf1c27008d66d436546cb0454443b772b9cb33429b3d1f9ee9848a10');
        const hash1 = await leagues.hashTeamState([2, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('hash init state', async () => {
        const state = [3, 435, 5];
        const hash0 = await leagues.hashTeamState(state).should.be.fulfilled;
        hash0.should.be.equal('0xed292958eade9f3b6aab7fe70037eb3115b8472dfceec6031c86cb1c4198dcf9');
        const hash1 = await leagues.hashTeamState([2, 3, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('hash league state', async () => {
        const state = [3,5,2,0,4,56,6];
        const hashes = await leagues.hashLeagueState(state).should.be.fulfilled;
        hashes.length.should.be.equal(2);
        hashes[0].should.be.equal('0x7c1a4da9ae4219e9f58d1c3b3efacea1c4e962cf3330297d46eceb26a8500221');
        hashes[1].should.be.equal('0x29b5c96991b3957cb253e235c95e56369e43542d3d1273bc916229afb773b205');
    });

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
});