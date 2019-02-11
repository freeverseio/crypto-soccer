const keccak = require('keccak');
require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('LeaguesScheduler');

contract('LeaguesScheduler', (accounts) => {
    let leagues = null;
    const id = 0;

    beforeEach(async () => {
        const blocksToInit = 0;
        const step = 1
        const teamIds = [1, 2];
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
    });

    it('get days of a league', async () => {
        const days = await leagues.countLeagueDays(id).should.be.fulfilled;
        days.toNumber().should.be.equal(2);
    });

    it('get days of a wrong league', async () => {
        await leagues.countLeagueDays(1).should.be.rejected;
    })

    it('get teams for match in wrong league day', async () => {
        const day = 2; // wrong
        const matchIdx = 0; 
        await leagues.getTeamsInMatch(id, day, matchIdx).should.be.rejected;
    });

    it('get teams for match in wrong team', async () => {
        const day = 0;
        const matchIdx = 2; // wrong
        await leagues.getTeamsInMatch(id, day, matchIdx).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        let day = 0;
        const matchIdx = 0;
        let teams = await leagues.getTeamsInMatch(id, day, matchIdx).should.be.fulfilled;
        teams.length.should.be.equal(2);
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        day = 1;
        teams = await leagues.getTeamsInMatch(id, day, matchIdx).should.be.fulfilled;
        teams.length.should.be.equal(2);
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    })
});