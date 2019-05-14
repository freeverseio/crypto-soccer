require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Horizon = artifacts.require('Horizon');

contract('Horizon', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;

    beforeEach(async () => {
        players = await Players.new().should.be.fulfilled;
        teams = await Teams.new(players.address).should.be.fulfilled;
        horizon = await Horizon.new(teams.address).should.be.fulfilled;
        await players.addMinter(horizon.address).should.be.fulfilled;
        await players.renounceMinter().should.be.fulfilled;
        await teams.addMinter(horizon.address).should.be.fulfilled;
        await teams.renounceMinter().should.be.fulfilled;
        await players.addTeamsContract(teams.address).should.be.fulfilled;
        await players.renounceTeamsContract().should.be.fulfilled;
    });

    it('create Team', async () => {
        await horizon.createTeam("team").should.be.fulfilled;
        const teamCount = await teams.totalSupply().should.be.fulfilled;
        teamCount.toNumber().should.be.equal(1);
        const playerCount = await players.totalSupply().should.be.fulfilled;
        playerCount.toNumber().should.be.equal(11);

        // TODO way to much for a single test
        const teamId = await teams.getTeamId("team").should.be.fulfilled;
        const teamPlayers = await teams.getPlayers(teamId).should.be.fulfilled;
        for(var i=0 ; i < teamPlayers.length ; i++){
            const playerId = teamPlayers[i];
            const result = await players.getTeam(playerId).should.be.fulfilled;
            result.toString().should.be.equal(teamId.toString());
        }
    })
});