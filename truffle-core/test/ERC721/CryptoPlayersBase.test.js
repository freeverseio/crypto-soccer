require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersBaseMock');

contract('CryptoPlayersBase', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    });

    it('no initial players', async () => {
        const count = await contract.totalSupply().should.be.fulfilled;
        count.toNumber().should.be.equal(0);
    });

    it('mint 2 player with same name', async () => {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "player").should.be.fulfilled;
        await contract.mintWithName(accounts[0], id, "player").should.be.rejected;
    });

    it('after transfer team is 0', async () => {
        const id = 1;
        const teamId = 1;
        await contract.mintWithName(accounts[0], id, "player").should.be.fulfilled;
        await contract.setTeam(id, teamId).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], id).should.be.fulfilled;
        const result = await contract.getTeam(id).should.be.fulfilled;
        result.toNumber().should.be.equal(0);
    });
});
