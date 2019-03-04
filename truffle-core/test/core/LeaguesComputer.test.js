require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const States = artifacts.require('DayState');
const Leagues = artifacts.require('LeaguesComputer');

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

    it('hash tactics', async () => {
        const hash0 = await leagues.hashTactics([[4, 4, 2]]).should.be.fulfilled;
        const hash1 = await leagues.hashTactics([[4, 4, 2]]).should.be.fulfilled;
        hash1.should.be.equal(hash0);
        const hash2 = await leagues.hashTactics([[3, 4, 2]]).should.be.fulfilled;
        hash2.should.be.not.equal(hash0);
        const hash3 = await leagues.hashTactics([[4, 5, 2]]).should.be.fulfilled;
        hash3.should.be.not.equal(hash0);
        const hash4 = await leagues.hashTactics([[4, 4, 3]]).should.be.fulfilled;
        hash4.should.be.not.equal(hash0);
        const hash5 = await leagues.hashTactics([[4, 4, 2], [4, 4, 2]]).should.be.fulfilled;
        hash5.should.be.not.equal(hash0);
    });

    it('hash team state', async () => {
        const state = [324, 435, 5];
        const hash0 = await leagues.hashTeamState(state).should.be.fulfilled;
        hash0.should.be.equal('0xd2bd5ff7bf1c27008d66d436546cb0454443b772b9cb33429b3d1f9ee9848a10');
        const hash1 = await leagues.hashTeamState([2, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('hash init state', async () => {
        const state = [3, 435, 5];
        const hash0 = await leagues.hashTeamState(state).should.be.fulfilled;
        hash0.should.be.equal('0xed292958eade9f3b6aab7fe70037eb3115b8472dfceec6031c86cb1c4198dcf9');
        const hash1 = await leagues.hashTeamState([2, 3, 3]).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
    });

    it('hash day state', async () => {
        const state = [3,5,2,0,4,56,6,0];
        const hashes = await leagues.hashDayState(state).should.be.fulfilled;
        hashes.length.should.be.equal(2);
        hashes[0].should.be.equal('0x7c1a4da9ae4219e9f58d1c3b3efacea1c4e962cf3330297d46eceb26a8500221');
        hashes[1].should.be.equal('0x29b5c96991b3957cb253e235c95e56369e43542d3d1273bc916229afb773b205');
    });

    it('update league state', async () => {
        const initStateHash = '0x435a354320000000000000000000000000000000000000000000000000000000';
        const finalTeamStateHashes = [
            '0x245964674ab00000000000000000000000000000000000000000000000000000',
            '0xaaaaa00000000000000000000000000000000000000000000000000000000000'
        ];
        const scores = [1, 2, 2, 2];
        await leagues.updateLeague(id, initStateHash, finalTeamStateHashes, scores).should.be.fulfilled;
        let result = await leagues.getInitStateHash(id).should.be.fulfilled;
        result.should.be.equal(initStateHash);
        result = await leagues.getFinalTeamStateHashes(id).should.be.fulfilled;
        result.length.should.be.equal(finalTeamStateHashes.length);
        result[0].should.be.equal(finalTeamStateHashes[0]);
        result[1].should.be.equal(finalTeamStateHashes[1]);
        result = await leagues.getScores(id).should.be.fulfilled;
        result.length.should.be.equal(scores.length);
        for (i = 0 ; i < result.length ; i++)
            result[i].toNumber().should.be.equal(scores[i])
    });

    it('calculate a day in a league', async () => {
        let day = 0;
        let result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x0').should.be.fulfilled;
        result.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(0);
        scores.visitor.toNumber().should.be.equal(2);
        day = 1;
        result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x4354646451').should.be.fulfilled;
        result.length.should.be.equal(1);
        scores = await leagues.decodeScore(result[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(1);
        scores.visitor.toNumber().should.be.equal(3);
    });

    it('result of a day in league is deterministic', async () => {
        const day = 1;
        let result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x123456').should.be.fulfilled;
        result.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(2);
        scores.visitor.toNumber().should.be.equal(1);
        result = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x123456').should.be.fulfilled;
        result.length.should.be.equal(1);
        scores = await leagues.decodeScore(result[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(2);
        scores.visitor.toNumber().should.be.equal(1);
    });

    it('different seed => different results', async () => {
        const day = 1;
        const result0 = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x1234').should.be.fulfilled;
        const result1 = await leagues.computeStatesAtMatchday(id, day, dayState, tactics, '0x4321').should.be.fulfilled;
        result0.length.should.be.equal(1);
        result1.length.should.be.equal(1);
        result0[0].toNumber().should.not.be.equal(result1[0].toNumber());
    });

    it('calculate all a league', async () => {
        const leagueScores = await leagues.computeAllMatchdayStates(id, dayState, tactics).should.be.fulfilled;
        const nDayScores = await leagues.countDaysInTournamentScores(leagueScores).should.be.fulfilled;
        nDayScores.toNumber().should.be.equal(2);
        let dayScores = await leagues.getDayScores(leagueScores, 0).should.be.fulfilled;
        dayScores.length.should.be.equal(1);
        dayScores = await leagues.getDayScores(leagueScores, 1).should.be.fulfilled;
        dayScores.length.should.be.equal(1);
    });

    it('compute points for winner', async () => {
        let points = await leagues.computePointsWon(teamState0, teamState0, 2, 2).should.be.fulfilled;
        points.toNumber().should.be.equal(5);
        points = await leagues.computePointsWon(teamState0, teamState0, 2, 1).should.be.fulfilled;
        points.toNumber().should.be.equal(5);
        points = await leagues.computePointsWon(teamState0, teamState0, 1, 2).should.be.fulfilled;
        points.toNumber().should.be.equal(5);
    });
});