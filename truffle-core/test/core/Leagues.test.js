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

    it('default values', async () =>{
        const block = await leagues.getInitBlock(id).should.be.fulfilled;
        block.toNumber().should.be.equal(0);
        const step = await leagues.getStep(id).should.be.fulfilled;
        step.toNumber().should.be.equal(0);
        const teamIds = await leagues.getTeamIds(id).should.be.fulfilled;
        teamIds.length.should.be.equal(0);
        const initHash = await leagues.getInitHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const hash = await leagues.getHash(id).should.be.fulfilled;
        hash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const end = await leagues.getEndBlock(id).should.be.fulfilled;
        end.toNumber().should.be.equal(0);
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
    })
});