require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('LeagueUpdatableMock');

contract('LeaguesUpdatable', (accounts) => {
    let leagues = null;
    const id = 0;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.create(
            id,
            blocksToInit = 1,
            step = 1,
            teamIds = [1, 2],
            tactics = [[4, 4, 3], [4, 4, 3]]
        ).should.be.fulfilled;
    });

    it('unexistent league', async () => {
        const id = 3;
        await leagues.getDayStateHashes(id).should.be.rejected;
        await leagues.getInitStateHash(id).should.be.rejected;
    })
    
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
    
    it('default hashes values on create league', async () => {
        const initHash = await leagues.getInitStateHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await leagues.getDayStateHashes(id).should.be.fulfilled;
        finalHashes.length.should.be.equal(0);
    });

    it('is updated', async () => {
        let result = await leagues.isUpdated(id).should.be.fulfilled;
        result.should.be.equal(false);
        const initStateHash = '0x54564';
        const dayStateHashes = ['0x24353', '0x5434432'];
        const scores = ['0x12', '0x3'];
        await leagues.updateLeague(id, initStateHash, dayStateHashes, scores).should.be.fulfilled;
        result = await leagues.isUpdated(id).should.be.fulfilled;
        result.should.be.equal(true);
    });
 
    it('updateBlock and updater', async () => {
        const initStateHash = '0x54564';
        const dayStateHashes = ['0x24353', '0x5434432'];
        const scores = ['0x12', '0x3'];
        const result = await leagues.updateLeague(id, initStateHash, dayStateHashes, scores).should.be.fulfilled;
        const updateBlock = await leagues.getUpdateBlock(id).should.be.fulfilled;
        updateBlock.toNumber().should.be.equal(result.receipt.blockNumber);
        const updater = await leagues.getUpdater(id).should.be.fulfilled;
        updater.should.be.equal(accounts[0]);
    });
})