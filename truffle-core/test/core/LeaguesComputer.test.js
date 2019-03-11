require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const States = artifacts.require('LeagueState');
const Leagues = artifacts.require('LeaguesComputerMock');

contract('LeaguesComputer', (accounts) => {
    let engine = null;
    let states = null;
    let leagues = null;

    const id = 0;
    const teamState0 = [1, 2, 3, 4, 5, 6, 7, 10, 9, 10, 11];
    const teamState1 = [11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1];
    const tactics = [
        [4, 4, 3],  // Team 0
        [5, 4, 2]   // Team 1
    ];
    let leagueState = null;

    beforeEach(async () => {
        const blocksToInit = 1;
        const step = 1
        const teamIds = [1, 2];
        engine = await Engine.new().should.be.fulfilled;
        states = await States.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, states.address).should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        leagueState = await states.leagueStateCreate().should.be.fulfilled;
        leagueState = await states.leagueStateAppend(leagueState, teamState0).should.be.fulfilled;
        leagueState = await states.leagueStateAppend(leagueState, teamState1).should.be.fulfilled;
    });

    it('Engine contract', async () => {
        const address = await leagues.getEngineContract().should.be.fulfilled;
        address.should.be.equal(engine.address);
    });

    it('evolve team', async () => {
        const result = await leagues.evolveTeams(teamState0, teamState1, 0, 2).should.be.fulfilled;
        let valid = await states.isValidTeamState(result.updatedHomeTeamState).should.be.fulfilled;
        valid.should.be.equal(true);
        valid = await states.isValidTeamState(result.updatedVisitorTeamState).should.be.fulfilled;
        valid.should.be.equal(true);
    })

    it('compute match', async () => {
        const result = await leagues.computeMatch(teamState0, teamState1, tactics, '0x0').should.be.fulfilled;
        const scores = await leagues.decodeScore(result.score).should.be.fulfilled;
        scores.home.toString().should.be.equal('0');
        scores.visitor.toString().should.be.equal('2');
        let valid = await states.isValidTeamState(result.newHomeState).should.be.fulfilled;
        valid.should.be.equal(true);
        valid = await states.isValidTeamState(result.newVisitorState).should.be.fulfilled;
        valid.should.be.equal(true);
    });

    it('calculate a day in a league', async () => {
        let day = 0;
        let result = await leagues.computeDayWithSeed(id, day, leagueState, tactics, '0x0').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(0);
        scores.visitor.toNumber().should.be.equal(2);
        day = 1;
        result = await leagues.computeDayWithSeed(id, day, leagueState, tactics, '0x4354646451').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(1);
        scores.visitor.toNumber().should.be.equal(3);
    });

    it('result of a day in league is deterministic', async () => {
        const day = 1;
        let result = await leagues.computeDayWithSeed(id, day, leagueState, tactics, '0x123456').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(2);
        scores.visitor.toNumber().should.be.equal(1);
        result = await leagues.computeDayWithSeed(id, day, leagueState, tactics, '0x123456').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(2);
        scores.visitor.toNumber().should.be.equal(1);
    });

    it('different seed => different results', async () => {
        const day = 1;
        const result0 = await leagues.computeDayWithSeed(id, day, leagueState, tactics, '0x1234').should.be.fulfilled;
        const result1 = await leagues.computeDayWithSeed(id, day, leagueState, tactics, '0x4321').should.be.fulfilled;
        result0.scores.length.should.be.equal(1);
        result1.scores.length.should.be.equal(1);
        result0.scores[0].toNumber().should.not.be.equal(result1.scores[0].toNumber());
    });

    it('compute points same rating', async () => {
        let points = await leagues.computePoints(teamState0, teamState0, 2, 2).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(0);
        points.visitorPoints.toNumber().should.be.equal(0);
        points = await leagues.computePoints(teamState0, teamState0, 2, 1).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(5);
        points.visitorPoints.toNumber().should.be.equal(0);
        points = await leagues.computePoints(teamState0, teamState0, 1, 2).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(0);
        points.visitorPoints.toNumber().should.be.equal(5);
    });

    // it('compute points home rating higher than visitor', async () => {
    //     const homeTeamRating = await states.computeTeamRating(teamState0).should.be.fulfilled;
    //     const visitorTeamRating = await states.computeTeamRating(teamState1).should.be.fulfilled;
    //     homeTeamRating.toNumber().should.be.equal(68);
    //     visitorTeamRating.toNumber().should.be.equal(66);
    //     let points = await leagues.computePoints(teamState0, teamState1, 2, 2).should.be.fulfilled;
    //     points.homePoints.toNumber().should.be.equal(0);
    //     points.visitorPoints.toNumber().should.be.equal(0);
    //     points = await leagues.computePoints(teamState0, teamState1, 2, 1).should.be.fulfilled;
    //     points.homePoints.toNumber().should.be.equal(2);
    //     points.visitorPoints.toNumber().should.be.equal(0);
    //     points = await leagues.computePoints(teamState0, teamState1, 1, 2).should.be.fulfilled;
    //     points.homePoints.toNumber().should.be.equal(0);
    //     points.visitorPoints.toNumber().should.be.equal(8);
    // });

    // it('compute points home rating lower than visitor', async () => {
    //     const homeTeamRating = await states.computeTeamRating(teamState1).should.be.fulfilled;
    //     const visitorTeamRating = await states.computeTeamRating(teamState0).should.be.fulfilled;
    //     homeTeamRating.toNumber().should.be.equal(66);
    //     visitorTeamRating.toNumber().should.be.equal(68);
    //     let points = await leagues.computePoints(teamState1, teamState0, 2, 2).should.be.fulfilled;
    //     points.homePoints.toNumber().should.be.equal(0);
    //     points.visitorPoints.toNumber().should.be.equal(0);
    //     points = await leagues.computePoints(teamState1, teamState0, 2, 1).should.be.fulfilled;
    //     points.homePoints.toNumber().should.be.equal(8);
    //     points.visitorPoints.toNumber().should.be.equal(0);
    //     points = await leagues.computePoints(teamState1, teamState0, 1, 2).should.be.fulfilled;
    //     points.homePoints.toNumber().should.be.equal(0);
    //     points.visitorPoints.toNumber().should.be.equal(2);
    // });
});