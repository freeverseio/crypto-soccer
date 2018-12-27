require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoTeams = artifacts.require('CryptoTeamsBaseMock');

contract('CryptoTeamsBase', (accounts) => {
    let contract = null;
    let cryptoPlayers = null;

    beforeEach(async () => {
        contract = await CryptoTeams.new().should.be.fulfilled;
    });

    it('no initial teams', async () => {
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(0);
    })

    it('mint team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
        const total = await contract.totalSupply().should.be.fulfilled;
        total.toNumber().should.be.equal(1);
    });

    it('mint team with same name is forbidden', async () =>  {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        await contract.mintWithName(accounts[0], "team").should.be.rejected;
    });

    it('team owner', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const owner = await contract.ownerOf(id).should.be.fulfilled;
        owner.should.be.equal(accounts[0]);
    });

    it('team name', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const name = await contract.getName(id).should.be.fulfilled;
        name.should.be.equal("team");
    });

    it('get name of unexistent team', async () => {
        await contract.getName(1).should.be.rejected;
    });

    it('getTeam team id', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
    });

    it('getTeam team id of unexistent tema', async () => {
        const id = await contract.getTeamId("team").should.be.rejected;
    });

    it('new team has no players', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const id = await contract.getTeamId("team").should.be.fulfilled;
        const players = await contract.getPlayers(id).should.be.fulfilled;
        players.length.should.be.equal(0);
    });

    it('add unexistent player to team', async () => {
        await contract.mintWithName(accounts[0], "team").should.be.fulfilled;
        const teamId = await contract.getTeamId("team").should.be.fulfilled;
        const unexistentPlayerId = 1;
        const position = 0;
        await contract.addPlayer(teamId, position, unexistentPlayerId).should.be.rejected;
    });
});