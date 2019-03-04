require('chai')
    .use(require('chai-as-promised'))
    .should();

const LeaguesProof = artifacts.require('LeaguesProof');

contract('LeaguesProof', (accounts) => {
    let instance = null;
    const initBlock = 1;
    const step = 1;
    const id = 0;
    const teamIds = [1, 2];

    beforeEach(async () => {
        instance = await LeaguesProof.new().should.be.fulfilled;
    });

    it('unexistent league', async () => {
        await instance.getFinalTeamStateHashes(id).should.be.rejected;
        await instance.getInitStateHash(id).should.be.rejected;
    });

    it('default hashes values on create league', async () =>{
        await instance.create(id, initBlock, step, teamIds).should.be.fulfilled;
        const initHash = await instance.getInitStateHash(id).should.be.fulfilled;
        initHash.should.be.equal('0x0000000000000000000000000000000000000000000000000000000000000000');
        const finalHashes = await instance.getFinalTeamStateHashes(id).should.be.fulfilled;
        finalHashes.length.should.be.equal(0);
    });
});