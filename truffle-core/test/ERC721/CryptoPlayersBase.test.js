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

    it('when players is sold, he has no team', async () => {
        const playerId = 1;
        const teamId = 1;
        await contract.mintWithName(accounts[0], playerId, "player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], playerId).should.be.fulfilled;
        const team = await contract.getTeam(playerId).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });
});
