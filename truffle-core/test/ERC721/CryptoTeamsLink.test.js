require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersLink');
const CryptoTeams = artifacts.require('CryptoTeamsLink');

contract('CryptoTeamsLink', (accounts) => {
    let contract = null;
    let cryptoPlayers = null;

    beforeEach(async () => {
        cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
        contract = await CryptoTeams.new().should.be.fulfilled;
        await contract.setPlayersContract(cryptoPlayers.address).should.be.fulfilled;
        await cryptoPlayers.setTeamsContract(contract.address).should.be.fulfilled;
    });

    it('add unexistent player to team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        const unexistentPlayerId = 1;
        const position = 0;
        await contract.addPlayer(teamId, position, unexistentPlayerId).should.be.rejected;
    });

    it('add existent player to team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await cryptoPlayers.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await cryptoPlayers.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
    });

    it('add player to team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await cryptoPlayers.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await cryptoPlayers.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
        const players = await contract.getPlayers(teamId).should.be.fulfilled;
        players.length.should.be.equal(1);
        players[0].toNumber().should.be.equal(playerId.toNumber());
    });

    it('selling team changes players ownership', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await cryptoPlayers.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await cryptoPlayers.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], teamId).should.be.fulfilled;
        const teamOwner = await contract.ownerOf(teamId).should.be.fulfilled;
        teamOwner.should.be.equal(accounts[1]);
        const playerOwner = await cryptoPlayers.ownerOf(playerId).should.be.fulfilled;
        playerOwner.should.be.equal(accounts[1]);
    });

    it('if team adds a player, player knows his team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await cryptoPlayers.mintWithName(accounts[0], "player").should.be.fulfilled;
        const playerId = await cryptoPlayers.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
        const team = await cryptoPlayers.getTeam(playerId).should.be.fulfilled;
        team.toNumber().should.be.equal(teamId.toNumber());
    });
});