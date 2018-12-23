require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayersBaseMock');
const CryptoTeams = artifacts.require('CryptoTeamsBaseMock');

contract('CryptoTeamsBase', (accounts) => {
    let contract = null;
    let cryptoPlayers = null;

    beforeEach(async () => {
        cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
        contract = await CryptoTeams.new().should.be.fulfilled;
        await contract.setPlayersContract(cryptoPlayers.address);
    });

    it('no initial teams', async () => {
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(0);
    })

    it('id 0 is forbidden', async () => {
        const id = 0;
        await contract.mintWithName(accounts[0], id, "team").should.be.rejected;
    });

    it('mint team', async () => {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(id);
    });

    it('mint team with id more tha 2**22 is forbidden', async () => {
        const id = 2**22;
        await contract.mintWithName(accounts[0], id+1, "team").should.be.rejected;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
    })

    it('mint team with an existent id', async () =>  {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        await contract.mintWithName(accounts[0], id, "team").should.be.rejected;
    });

    it('team owner', async () => {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        const owner = await contract.ownerOf(id).should.be.fulfilled;
        owner.should.be.equal(accounts[0]);
    });

    it('team name', async () => {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
    });

    it('get name of unexistent team', async () => {
        await contract.getName(1).should.be.rejected;
    });

    it('get playersIds of unexistent team', async () => {
        await contract.getPlayersIds(1).should.be.rejected;
    });

    it('set playersIds of unexistent team', async () => {
        await contract.setPlayersIds(1, 0).should.be.rejected;
    });

    it('set playersIds', async () => {
        const id = 1;
        const playersIds = 31231234;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        await contract.setPlayersIds(id, playersIds).should.be.fulfilled;
        const result = await contract.getPlayersIds(id).should.be.fulfilled;
        result.toNumber().should.be.equal(playersIds);
    });

    it('get team id by name', async () => {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        const result = await contract.getTeamId("team").should.be.fulfilled;
        result.toNumber().should.be.equal(id);
        await contract.getTeamId("team1").should.be.rejected;
    });

    it('new team has no players', async () => {
        const id = 1;
        await contract.mintWithName(accounts[0], id, "team").should.be.fulfilled;
        const playerIds = await contract.getPlayersIds(id).should.be.fulfilled;
        playerIds.toNumber().should.be.equal(0);
    });

    it('add unexistent player to team', async () => {
        const teamId = 1;
        await contract.mintWithName(accounts[0], teamId, "team").should.be.fulfilled;
        const unexistentPlayerId = 1;
        const position = 0;
        await contract.addPlayer(teamId, position, unexistentPlayerId).should.be.rejected;
    });

    it('add existent player to team', async () => {
        const teamId = 1;
        const playerId = 1;
        await contract.mintWithName(accounts[0], teamId, "team").should.be.fulfilled;
        await cryptoPlayers.mintWithName(accounts[0], playerId, "player").should.be.fulfilled;
        const position = 0;
        await contract.addPlayer(teamId, position, playerId).should.be.fulfilled;
    });

    it('add player to team', async () => {
        const teamId = 1;
        const playerId = 1;
        await contract.mintWithName(accounts[0], teamId, "team").should.be.fulfilled;
        await cryptoPlayers.mintWithName(accounts[0], playerId, "player").should.be.fulfilled;
        const players = await contract.getPlayersIds(teamId).should.be.fulfilled;
        players.toNumber().should.be.equal(4);
    });

    // it('selling team changes players ownership', async () => {
    //     const playerId = 1;
    //     const teamId = 1;
    //     await contract.mintWithName(accounts[0], teamId, "team").should.be.fulfilled;
    //     await cryptoPlayers.mintWithName(accounts[0], playerId, "player").should.be.fulfilled;
    //     await cryptoPlayers.setTeam(playerId, teamId).should.be.fulfilled;
    //     await contract.safeTransferFrom(accounts[0], accounts[1], teamId).should.be.fulfilled;
    //     const teamOwner = await contract.ownerOf(teamId).should.be.fulfilled;
    //     teamOwner.should.be.equal(accounts[1]);
    //     const playerOwner = await cryptoPlayers.ownerOf(playerId).should.be.fulfilled;
    //     playerOwner.should.be.equal(accounts[1]);
    // });
});