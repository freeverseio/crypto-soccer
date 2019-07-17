require('chai')
    .use(require('chai-as-promised'))
    .should();

const Leagues = artifacts.require('LeaguesBaseMock');
const Assets = artifacts.require('Assets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('LeaguesBase', (accounts) => {
    let leagues = null;
    let assets = null;
    let playerStateLib = null;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const reverseOrder = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => PLAYERS_PER_TEAM-i-1) // [24,23,...0]
    const tactic442 = 1;
    const tactic433 = 2;

    const initBlock = 1;
    const step = 1;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        leagues = await Leagues.new().should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
    });

    
    it('add a created team to non-created league', async () => {
        await leagues.signTeamInLeague(
            leagueId = 1,
            teamId = 1,
            order,
            tactic442
        ).should.be.rejected;
    });

    it('add a non-created team to a created league', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(
            leagueId = 1,
            teamId = 3,
            order,
            tactic442
        ).should.be.rejected;
    });

    
    it('add created team to created league', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(
            leagueId = 1,
            teamId = 1,
            order,
            tactic442
        ).should.be.fulfilled;
    });

    it('add created team twice to created league', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(
            leagueId = 1,
            teamId = 1,
            order,
            tactic442
        ).should.be.fulfilled;
        await leagues.signTeamInLeague(
            leagueId = 1,
            teamId = 1,
            order,
            tactic442
        ).should.be.rejected;
    });

    it('add too many teams to league', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId = 1, teamId = 1, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId = 1, teamId = 2, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId = 1, teamId = 3, order, tactic442).should.be.rejected;
    });
    

    it('unexistent league', async () => {
        // leagueId = 0 is dummy
        await leagues.getInitBlock(id = 1).should.be.rejected;
        await leagues.getStep(id = 1).should.be.rejected;
        await leagues.getNTeams(id = 1).should.be.rejected;
    });

    it('create league with no team', async () => {
        await leagues.create(nTeams = 0, initBlock, step).should.be.rejected;
    });

    it('create league with odd teams', async () => {
        await leagues.create(nTeams = 1, initBlock, step).should.be.rejected;
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.create(nTeams = 7, initBlock, step).should.be.rejected;
    });

    it('check leagueId and LeagueCount for league with 2 teams', async () => {
        const receipt = await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        const leagueId = receipt.logs[0].args.leagueId.toNumber();
        leagueId.should.be.equal(1);
        const count = await leagues.leaguesCount().should.be.fulfilled;
        count.toNumber().should.be.equal(1);
    });

    it('init block of a league', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        const result = await leagues.getInitBlock(id = 1).should.be.fulfilled;
        result.toNumber().should.be.equal(initBlock);
    });

    it('count leagues', async () => {
        let counter = await leagues.leaguesCount().should.be.fulfilled;
        counter.toNumber().should.be.equal(0);
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        counter = await leagues.leaguesCount().should.be.fulfilled;
        counter.toNumber().should.be.equal(1);
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        counter = await leagues.leaguesCount().should.be.fulfilled;
        counter.toNumber().should.be.equal(2);
    });

    it('step == 0 is invalid', async () => {
        await leagues.create(nTeams = 2, initBlock, thisStep = 0).should.be.rejected;
    });
    
    it('count teams', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        const count = await leagues.getNTeams(id = 1).should.be.fulfilled;
        count.toNumber().should.be.equal(2);
    });

    it('hash users init data', async () => {
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        let usersInitDataHash = await leagues.getUsersInitDataHash(leagueId = 1).should.be.fulfilled;
        usersInitDataHash = web3.utils.hexToNumber(usersInitDataHash)
        usersInitDataHash.should.be.equal(0);
        await leagues.signTeamInLeague(leagueId = 1, teamId = 1, order, tactic442).should.be.fulfilled;
        usersInitDataHash = await leagues.getUsersInitDataHash(leagueId = 1).should.be.fulfilled;
        usersInitDataHash.should.be.equal('0x7e22cf54171452bfaf39fbc0c4cbf6d4adf7cb4c955d799207e3e7056d187921');
    });
});