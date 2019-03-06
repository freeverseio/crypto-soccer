require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('LeagueUpdater');

contract('LeaguesUpdater', (accounts) => {
    let leagues = null;
    let states = null;
    const id = 0;

    beforeEach(async () => {
        const blocksToInit = 1;
        const step = 1
        const teamIds = [1, 2];
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
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

    it('hash state', async () => {
        const state = [324, 435, 5];
        const hash0 = await leagues.hashState(state).should.be.fulfilled;
        hash0.should.be.equal('0xd2bd5ff7bf1c27008d66d436546cb0454443b772b9cb33429b3d1f9ee9848a10');
        const hash1 = await leagues.hashState([2, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('update league state', async () => {
        const initStateHash = '0x435a354320000000000000000000000000000000000000000000000000000000';
        const finalTeamStateHashes = [
            '0x245964674ab00000000000000000000000000000000000000000000000000000',
            '0xaaaaa00000000000000000000000000000000000000000000000000000000000'
        ];
        const scores = [1, 2, 2, 2];
        await leagues.updateLeague(id, initStateHash, finalTeamStateHashes, scores).should.be.fulfilled;
        let result = await leagues.getInitStateHash(id).should.be.fulfilled;
        result.should.be.equal(initStateHash);
        result = await leagues.getFinalTeamStateHashes(id).should.be.fulfilled;
        result.length.should.be.equal(finalTeamStateHashes.length);
        result[0].should.be.equal(finalTeamStateHashes[0]);
        result[1].should.be.equal(finalTeamStateHashes[1]);
        result = await leagues.getScores(id).should.be.fulfilled;
        result.length.should.be.equal(scores.length);
        for (i = 0 ; i < result.length ; i++)
            result[i].toNumber().should.be.equal(scores[i])
    });
})