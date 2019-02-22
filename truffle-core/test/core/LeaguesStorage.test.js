require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('LeaguesStorageMock');

contract('LeaguesStorage', (accounts) => {
    let leagues = null;
    const initBlock = 1;
    const step = 1;
    const id = 0;
    const teamIds = [1, 2];

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('unexistent league', async () => {
        await leagues.getInitBlock(id).should.be.rejected;
        await leagues.getStep(id).should.be.rejected;
        await leagues.getTeamIds(id).should.be.rejected;
        await leagues.getInitHash(id).should.be.rejected;
        await leagues.getFinalTeamStateHashes(id).should.be.rejected;
        await leagues.getEndBlock(id).should.be.rejected;
        await leagues.countTeams(id).should.be.rejected;
        await leagues.getInitStateHash(id).should.be.rejected;
    })

    it('default hashes values on create league', async () =>{
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const initHash = await leagues.getInitHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await leagues.getFinalTeamStateHashes(id).should.be.fulfilled;
        finalHashes.length.should.be.equal(0);
    })

    it('create league with no team', async () => {
        const teamIds = [];
        await leagues.create(id, initBlock, step, teamIds).should.be.rejected;
    });

    it('create league with 1 team', async () => {
        const teamIds = [1];
        await leagues.create(id, initBlock, step, teamIds).should.be.rejected;
    });

    it('create league with 2 teams', async () => {
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const result = await leagues.getTeamIds(id).should.be.fulfilled;
        result.length.should.be.equal(2);
        result[0].toNumber().should.be.equal(1);
        result[1].toNumber().should.be.equal(2);
    });

    it('create leagues with odd teams', async () => {
        await leagues.create(id, initBlock, step, [1, 2, 3]).should.be.rejected;
        await leagues.create(id, initBlock, step, [1, 2, 3, 4, 5]).should.be.rejected;
        await leagues.create(id, initBlock, step, [1, 2, 3, 4, 5, 6, 7]).should.be.rejected;
    });

    it('init block of a league', async () => {
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const result = await leagues.getInitBlock(id).should.be.fulfilled;
        result.toNumber().should.be.equal(initBlock);
    });

    it('end block of a league', async () => {
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const endBlock = await leagues.getEndBlock(id).should.be.fulfilled;
        endBlock.toNumber().should.be.equal(initBlock + step);
    });

    it('create 2 leagues with the same id', async () => {
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        await leagues.create(id, initBlock, step, teamIds).should.be.rejected;
    });

    it('step == 0 is invalid', async () => {
        const step = 0;
        await leagues.create(id, initBlock, step, teamIds).should.be.rejected;
    });

    it('count teams', async () => {
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const count = await leagues.countTeams(id).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('set scores', async () => {
        await leagues.create(id, initBlock, step, teamIds).should.be.fulfilled;
        await leagues.setScores(id, [[0, 1], [2, 2]]).should.be.fulfilled;
        const scores = await leagues.getScores(id).should.be.fulfilled;
        // TODO: finish it
    })
});