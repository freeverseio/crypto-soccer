require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('LeaguesScheduler');
const Assets = artifacts.require('Assets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('LeaguesScheduler', (accounts) => {
    let leagues = null;
    let assets = null;
    let playerStateLib = null;
    const initBlock = 1;
    const step = 1;
    const leagueId = 1;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 0;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 1, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 2, order, tactic442).should.be.fulfilled;
    });

    it('end block of unexistend league', async () => {
        await leagues.getEndBlock(thisLeagueId = 2).should.be.rejected;
    })

    it('end block of a league', async () => {
        const endBlock = await leagues.getEndBlock(leagueId).should.be.fulfilled;
        endBlock.toNumber().should.be.equal(2);
    });

    it('get days of a league', async () => {
        const days = await leagues.countLeagueDays(leagueId).should.be.fulfilled;
        days.toNumber().should.be.equal(2);
    });

    it('get days of a wrong league', async () => {
        await leagues.countLeagueDays(thisLeagueId = 2).should.be.rejected;
    })

    it('get teams for match in wrong league day', async () => {
        const day = 2; // wrong
        const matchIdx = 0; 
        await leagues.getTeamsInMatch(leagueId, day, matchIdx).should.be.rejected;
    });

    it('get teams for match in wrong team', async () => {
        const day = 0;
        const matchIdx = 2; // wrong
        await leagues.getTeamsInMatch(leagueId, day, matchIdx).should.be.rejected;
    });

    it('get teams for match in league day', async () => {
        let day = 0;
        const matchIdx = 0;
        let teams = await leagues.getTeamsInMatch(leagueId, day, matchIdx).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        day = 1;
        teams = await leagues.getTeamsInMatch(leagueId, day, matchIdx).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });

    it('get match day', async () => {
        let hash = await leagues.getMatchDayBlockHash(leagueId, 0).should.be.fulfilled;
        hash.should.be.equal("0xb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6");
    });
});