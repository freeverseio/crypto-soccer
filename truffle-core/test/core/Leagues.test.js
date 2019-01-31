require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('Leagues');

contract('Leagues', (accounts) => {
    let leagues = null;

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('default init hash is 0', async () =>{
        const init = await leagues.getInit().should.be.fulfilled;
        init.toNumber().should.be.equal(0);
    });

    it('default final hash is 0', async () =>{
        const final = await leagues.getFinal().should.be.fulfilled;
        final.toNumber().should.be.equal(0);
    });

    it('default team ids is empty', async () => {
        const teamIds = await leagues.getTeamIds().should.be.fulfilled;
        teamIds.length.should.be.equal(0);
    })

    it('create league with state 0', async () => {
        await leagues.createLeague(0, []).should.be.rejected;
    });

    it('update league with state 0', async () => {
        await leagues.updateLeague(0).should.be.rejected;
    });

    it('update not created league', async () => {
        await leagues.updateLeague(2).should.be.rejected;
    });


    it('create league with 2 teams', async () => {
        const initState = 1;
        const teamIds = [1, 2];
        await leagues.createLeague(initState, teamIds).should.be.fulfilled;
    })

    it('create league with state 1', async () => {
        await leagues.createLeague(1, []).should.be.fulfilled;
        const init = await leagues.getInit().should.be.fulfilled;
        init.toNumber().should.be.equal(1);
    });

    it('update league with state 1', async () => {
        await leagues.createLeague(1, []).should.be.fulfilled;
        await leagues.updateLeague(2).should.be.fulfilled;
        const status = await leagues.getFinal().should.be.fulfilled;
        status.toNumber().should.be.equal(2);
    });
});