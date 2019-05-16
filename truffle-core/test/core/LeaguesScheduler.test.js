require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('LeaguesScheduler');

contract('LeaguesScheduler', (accounts) => {
    let leagues = null;
    const id = 0;

    beforeEach(async () => {
        const initBlock = 1;
        const step = 1;
        const teamIds = [1, 2];
        const tactics = [[4, 4, 3], [4, 4, 3]];
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.fulfilled;
    });

    it('end block of unexistend league', async () => {
        const id = 1;
        await leagues.getEndBlock(id).should.be.rejected;
    })

    it('end block of a league', async () => {
        const endBlock = await leagues.getEndBlock(id).should.be.fulfilled;
        endBlock.toNumber().should.be.equal(2);
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
        teams[0].toNumber().should.be.equal(0);
        teams[1].toNumber().should.be.equal(1);
        day = 1;
        teams = await leagues.getTeamsInMatch(id, day, matchIdx).should.be.fulfilled;
        teams[0].toNumber().should.be.equal(1);
        teams[1].toNumber().should.be.equal(0);
    });

    it('get match day', async () => {
        let hash = await leagues.getMatchDayBlockHash(id, 0).should.be.fulfilled;
        hash.should.be.equal("0xb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6");
    });
});