require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Horizon = artifacts.require('Horizon');
const Leagues = artifacts.require('LeaguesComputer');
const Engine = artifacts.require('Engine');

contract('Game', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;
    let engine = null;
    let leagues = null;

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

        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address).should.be.fulfilled;
    });

    it('play a league of 2 teams', async () => {
        await horizon.createTeam("Barcelona").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        const barcelonaId = await teams.getTeamId("Barcelona").should.be.fulfilled;
        const madridId = await teams.getTeamId("Madrid").should.be.fulfilled;
        const blockInitDelta = 0;
        const step = 1;
        const leagueId = 0;
        await leagues.create(leagueId, blockInitDelta, step, [barcelonaId, madridId]).should.be.fulfilled;

    });
})
