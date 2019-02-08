require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;
    const blocksToInit = 1;
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
        await leagues.getFinalHashes(id).should.be.rejected;
        await leagues.getEndBlock(id).should.be.rejected;
        await leagues.countTeams(id).should.be.rejected;
    })

    it('default hashes values on create league', async () =>{
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const initHash = await leagues.getInitHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await leagues.getFinalHashes(id).should.be.fulfilled;
        finalHashes.length.should.be.equal(teamIds.length);
        finalHashes.forEach(hash => (hash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000')));
    })

    it('create league with no team', async () => {
        const teamIds = [];
        await leagues.create(id, blocksToInit, step, teamIds).should.be.rejected;
    });

    it('create league with 1 team', async () => {
        const teamIds = [1];
        await leagues.create(id, blocksToInit, step, teamIds).should.be.rejected;
    });

    it('create league with 2 teams', async () => {
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const result = await leagues.getTeamIds(id).should.be.fulfilled;
        result.length.should.be.equal(2);
        result[0].toNumber().should.be.equal(1);
        result[1].toNumber().should.be.equal(2);
    });

    it('init block of a league', async () => {
        const result = await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const blockNumber = result.receipt.blockNumber;
        const initBlock = await leagues.getInitBlock(id).should.be.fulfilled;
        initBlock.toNumber().should.be.equal(blockNumber + blocksToInit);
    });

    it('end block of a league', async () => {
        const result = await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const blockNumber = result.receipt.blockNumber;
        const endBlock = await leagues.getEndBlock(id).should.be.fulfilled;
        endBlock.toNumber().should.be.equal(blockNumber + blocksToInit + step);
    });

    it('create 2 leagues with the same id', async () => {
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.rejected;
    });

    it('step == 0 is invalid', async () => {
        const step = 0;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.rejected;
    });

    it('count teams', async () => {
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        const count = await leagues.countTeams(id).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    })
});