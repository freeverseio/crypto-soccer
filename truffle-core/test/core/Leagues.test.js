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
        const teamIds = [1, 2];
        await leagues.create(0, teamIds).should.be.rejected;
    });

    it('update league with state 0', async () => {
        await leagues.update(0).should.be.rejected;
    });

    it('update not created league', async () => {
        await leagues.update(2).should.be.rejected;
    });

    it('create league with no team', async () => {
        const initState = 1;
        const teamIds = [];
        await leagues.create(initState, teamIds).should.be.rejected;
    });

    it('create league with 1 team', async () => {
        const initState = 1;
        const teamIds = [1];
        await leagues.create(initState, teamIds).should.be.rejected;
    });

    it('create league with 2 teams', async () => {
        const initState = 1;
        const teamIds = [1, 2];
        await leagues.create(initState, teamIds).should.be.fulfilled;
        const result = await leagues.getTeamIds().should.be.fulfilled;
        result.length.should.be.equal(2);
        result[0].toNumber().should.be.equal(1);
        result[1].toNumber().should.be.equal(2);
    });

    it('create league with state 1', async () => {
        const teamIds = [1, 2];
        await leagues.create(1, teamIds).should.be.fulfilled;
        const init = await leagues.getInit().should.be.fulfilled;
        init.toNumber().should.be.equal(1);
    });

    it('update league with state 1', async () => {
        const teamIds = [1, 2];
        await leagues.create(1, teamIds).should.be.fulfilled;
        await leagues.update(2).should.be.fulfilled;
        const status = await leagues.getFinal().should.be.fulfilled;
        status.toNumber().should.be.equal(2);
    });
});