require('chai')
    .use(require('chai-as-promised'))
    .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');
const CryptoTeams = artifacts.require('CryptoTeams');
const Horizon = artifacts.require('Horizon');

contract('Horizon', (accounts) => {
    let horizon = null;
    let cryptoPlayers = null;
    let cryptoTeams = null;

    beforeEach(async () => {
        cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
        cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
        horizon = await Horizon.new(cryptoPlayers.address, cryptoTeams.address).should.be.fulfilled;
        await cryptoPlayers.addMinter(horizon.address).should.be.fulfilled;
        await cryptoPlayers.renounceMinter().should.be.fulfilled;
        await cryptoTeams.addMinter(horizon.address).should.be.fulfilled;
        await cryptoTeams.renounceMinter().should.be.fulfilled;
        await cryptoPlayers.addTeamsContract(cryptoTeams.address).should.be.fulfilled;
        await cryptoPlayers.renounceTeamsContract().should.be.fulfilled;
        await cryptoTeams.setPlayersContract(cryptoPlayers.address).should.be.fulfilled;
    });

    it('create Team', async () => {
        await horizon.createTeam("team").should.be.fulfilled;
        const teamCount = await cryptoTeams.totalSupply().should.be.fulfilled;
        teamCount.toNumber().should.be.equal(1);
        const playerCount = await cryptoPlayers.totalSupply().should.be.fulfilled;
        playerCount.toNumber().should.be.equal(11);

        // TODO way to much for a single test
        const teamId = await cryptoTeams.getTeamId("team").should.be.fulfilled;
        const players = await cryptoTeams.getPlayers(teamId).should.be.fulfilled;
        for(var i=0 ; i < players.length ; i++){
            const playerId = players[i];
            const result = await cryptoPlayers.getTeam(playerId).should.be.fulfilled;
            result.toNumber().should.be.equal(teamId.toNumber());
        }
    })
});