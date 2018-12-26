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
        await contract.mintWithName(accounts[0], "player").should.be.fulfilled;
        await contract.mintWithName(accounts[0], "player").should.be.rejected;
    });

    it('get name', async () => {
        await contract.mintWithName(accounts[0], "player").should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("player");
    });

    it('get player id of existing player', async () => {
        await contract.mintWithName(accounts[0], "player").should.be.fulfilled;
        await contract.getPlayerId("player").should.be.fulfilled;
    });

    it('get player id of unexisting player', async () => {
        await contract.getPlayerId("player").should.be.rejected;
    });

    it('when players is sold, he has no team', async () => {
        const teamId = 1;
        await contract.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], playerId).should.be.fulfilled;
        const team = await contract.getTeam(playerId).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });
});
