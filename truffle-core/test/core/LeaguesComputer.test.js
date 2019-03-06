require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const States = artifacts.require('DayState');
const Leagues = artifacts.require('LeaguesComputerMock');

contract('LeaguesComputer', (accounts) => {
    let engine = null;
    let states = null;
    let leagues = null;

    const id = 0;
    const teamState0 = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11];
    const teamState1 = [11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1];
    const tactics = [
        [4, 4, 3],  // Team 0
        [5, 4, 2]   // Team 1
    ];
    let dayState = null;

    beforeEach(async () => {
        const blocksToInit = 1;
        const step = 1
        const teamIds = [1, 2];
        engine = await Engine.new().should.be.fulfilled;
        states = await States.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, states.address).should.be.fulfilled;
        await leagues.create(id, blocksToInit, step, teamIds).should.be.fulfilled;
        dayState = await states.dayStateCreate().should.be.fulfilled;
        dayState = await states.dayStateAppend(dayState, teamState0).should.be.fulfilled;
        dayState = await states.dayStateAppend(dayState, teamState1).should.be.fulfilled;
    });

    it('Engine contract', async () => {
        const address = await leagues.getEngineContract().should.be.fulfilled;
        address.should.be.equal(engine.address);
    });

    it('estimate gas cost in calculate a day', async () => {
        const day = 0;
        let cost = await leagues.computeStatesAtMatchday.estimateGas(id, day, dayState, tactics, '0x0').should.be.fulfilled;
        cost.should.be.equal(196687);
    })

    it('calculate a day in a league', async () => {
        let day = 0;
        let result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x0').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(0);
        scores.visitor.toNumber().should.be.equal(2);
        day = 1;
        result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x4354646451').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(1);
        scores.visitor.toNumber().should.be.equal(3);
    });

    it('result of a day in league is deterministic', async () => {
        const day = 1;
        let result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x123456').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(2);
        scores.visitor.toNumber().should.be.equal(1);
        result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x123456').should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(2);
        scores.visitor.toNumber().should.be.equal(1);
    });

    it('different seed => different results', async () => {
        const day = 1;
        const result0 = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x1234').should.be.fulfilled;
        const result1 = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x4321').should.be.fulfilled;
        result0.scores.length.should.be.equal(1);
        result1.scores.length.should.be.equal(1);
        result0.scores[0].toNumber().should.not.be.equal(result1.scores[0].toNumber());
    });

    // it('calculate all a league', async () => {
    //     const leagueScores = await leagues.computeAllMatchdayStates(id, dayState, tactics).should.be.fulfilled;
    //     const nDayScores = await leagues.countDaysInTournamentScores(leagueScores).should.be.fulfilled;
    //     nDayScores.toNumber().should.be.equal(2);
    //     let dayScores = await leagues.getDayScores(leagueScores, 0).should.be.fulfilled;
    //     dayScores.length.should.be.equal(1);
    //     dayScores = await leagues.getDayScores(leagueScores, 1).should.be.fulfilled;
    //     dayScores.length.should.be.equal(1);
    // });

    it('compute points for winner', async () => {
        let points = await leagues.computePointsWon(teamState0, teamState0, 2, 2).should.be.fulfilled;
        points.toNumber().should.be.equal(5);
        points = await leagues.computePointsWon(teamState0, teamState0, 2, 1).should.be.fulfilled;
        points.toNumber().should.be.equal(5);
        points = await leagues.computePointsWon(teamState0, teamState0, 1, 2).should.be.fulfilled;
        points.toNumber().should.be.equal(5);
    });
});