require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('LeaguesBaseMock');

contract('LeaguesBase', (accounts) => {
    let leagues = null;
    const initBlock = 1;
    const step = 1;
    const id = 0;
    const teamIds = [1, 2];
    const tactics = [[4,4,3], [4,5,2]];

    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
    });

    it('unexistent league', async () => {
        await leagues.getInitBlock(id).should.be.rejected;
        await leagues.getStep(id).should.be.rejected;
        await leagues.getNTeams(id).should.be.rejected;
    });

    it('create league with no team', async () => {
        const teamIds = [];
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.rejected;
    });

    it('create league with 1 team', async () => {
        const teamIds = [1];
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.rejected;
    });

    it('create league with 2 teams', async () => {
        const receipt = await leagues.create(initBlock, step, 2).should.be.fulfilled;
        const leagueId = receipt.logs[0].args.id.toNumber();
        leagueId.should.be.equal(1);
        const count = await leagues.leaguesCount().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });
    return;

    it('create leagues with odd teams', async () => {
        await leagues.create(id, initBlock, step, [1], [[4,4,3]]).should.be.rejected;
    });

    it('init block of a league', async () => {
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.fulfilled;
        const result = await leagues.getInitBlock(id).should.be.fulfilled;
        result.toNumber().should.be.equal(initBlock);
    });

    it('create 2 leagues with the same id', async () => {
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.fulfilled;
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.rejected;
    });

    it('count leagues', async () => {
        let counter = await leagues.leaguesCount().should.be.fulfilled;
        counter.toNumber().should.be.equal(0);
        await leagues.create(1, initBlock, step, teamIds, tactics).should.be.fulfilled;
        counter = await leagues.leaguesCount().should.be.fulfilled;
        counter.toNumber().should.be.equal(1);
        await leagues.create(2, initBlock, step, teamIds, tactics).should.be.fulfilled;
        counter = await leagues.leaguesCount().should.be.fulfilled;
        counter.toNumber().should.be.equal(2);
    });

    it('step == 0 is invalid', async () => {
        const step = 0;
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.rejected;
    });

    it('count teams', async () => {
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.fulfilled;
        const count = await leagues.getNTeams(id).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('hash users init data', async () => {
        const hash = await leagues.hashUsersInitData(teamIds, tactics).should.be.fulfilled;
        hash.should.be.equal('0xf8a82ba6630ed0305c4d7718ec5f87567f404ebffc7ddd22a344831368bf4537');
        await leagues.create(id, initBlock, step, teamIds, tactics).should.be.fulfilled;
        const usersInitDataHash = await leagues.getUsersInitDataHash(id).should.be.fulfilled;
        usersInitDataHash.should.be.equal(hash);
    });
});