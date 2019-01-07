require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersTeam');

contract('CryptoPlayersTeam', (accounts) => {
    let contract = null;
    const coach = accounts[1];

    beforeEach(async () => {
        contract = await CryptoPlayers.new().should.be.fulfilled;
        await contract.addCoach(coach).should.be.fulfilled;
        await contract.renounceCoach().should.be.fulfilled;
    });

    it('player has no team by default', async () => {
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const id = await contract.getPlayerId("player").should.be.fulfilled;
        const team = await contract.getTeam(id).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });
     
    it('no coach account can\'t set the team of the player', async () =>{
        const teamId = 1;
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId).should.be.rejected;
    });

    it('coach can set the team', async () => {
        const teamId = 1;
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId, {from: coach}).should.be.fulfilled;
    });

    it('when players is sold, he has no team', async () => {
        const teamId = 1;
        await contract.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await contract.getPlayerId("player").should.be.fulfilled;
        await contract.setTeam(playerId, teamId, {from: coach}).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], playerId).should.be.fulfilled;
        const team = await contract.getTeam(playerId).should.be.fulfilled;
        team.toNumber().should.be.equal(0);
    });
});
