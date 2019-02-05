require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('PlayersTeam');
const Teams = artifacts.require('TeamsPlayers');

contract('TeamsPlayers', (accounts) => {
    let contract = null;
    let players = null;

    beforeEach(async () => {
        players = await Players.new().should.be.fulfilled;
        contract = await Teams.new(players.address).should.be.fulfilled;
        await players.addTeamsContract(contract.address).should.be.fulfilled;
    });

    it('check players address', async () => {
        const result = await contract.getPlayersAddress().should.be.fulfilled;
        result.should.be.equal(players.address);
    });

    it('add unexistent player to team', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        const unexistentPlayerId = 1;
        const position = 0;
        await contract.addPlayer(teamId, position, unexistentPlayerId).should.be.rejected;
    });

    it('add existent player to team', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await players.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await players.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
    });

    it('add player to team', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await players.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await players.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
        const teamPlayers = await contract.getPlayers(teamId).should.be.fulfilled;
        teamPlayers.length.should.be.equal(1);
        teamPlayers[0].toNumber().should.be.equal(playerId.toNumber());
    });

    it('selling team changes players ownership', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await players.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await players.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
        await contract.safeTransferFrom(accounts[0], accounts[1], teamId).should.be.fulfilled;
        const teamOwner = await contract.ownerOf(teamId).should.be.fulfilled;
        teamOwner.should.be.equal(accounts[1]);
        const playerOwner = await players.ownerOf(playerId).should.be.fulfilled;
        playerOwner.should.be.equal(accounts[1]);
    });

    it('if team adds a player, player knows his team', async () => {
        await contract.mint(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        await players.mint(accounts[0], "player").should.be.fulfilled;
        const playerId = await players.getPlayerId("player").should.be.fulfilled;
        await contract.addPlayer(teamId, playerId).should.be.fulfilled;
        const team = await players.getTeam(playerId).should.be.fulfilled;
        team.toNumber().should.be.equal(teamId.toNumber());
    });
});