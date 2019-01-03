require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersTeam');

contract('CryptoPlayersTeam', (accounts) => {
    let contract = null;

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
    });

    it('default team', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        const team = await contract.getTeam(id).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });
     
    it('when players is sold, he has no team', async () => {
        const teamId = 1;
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], playerId).should.be.fulfilled;
        const team = await contract.getTeam(playerId).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });

    it('set team not possible by any account', async () =>{
        const teamId = 1;
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId, {from: accounts[1]}).should.be.rejected;
    });

    it('set team possible by team account', async () => {
        const teamId = 1;
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.fulfilled;
    })
});
