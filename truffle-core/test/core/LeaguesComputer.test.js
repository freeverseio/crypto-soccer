require('chai')
    .use(require('chai-as-promised'))
    .should();

const Engine = artifacts.require('Engine');
const States = artifacts.require('LeagueState');
const Leagues = artifacts.require('LeaguesComputerMock');
const Assets = artifacts.require('Assets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('LeaguesComputer', (accounts) => {
    let engine = null;
    let states = null;
    let leagues = null;
    let assets = null;
    let playerStateLib = null;
    const initBlock = 1;
    const step = 1;
    const leagueId = 1;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 0;
    const tactic541 = 1;
    const tacticsIds = [tactic442, tactic541];
    let teamStateAll1 = null;
    let teamStateAll50 = null;
    let leagueState = null;

    const createTeamStateFromSinglePlayer = async (defence, speed, pass, shoot, endurance, teamStateLib) => {
        const playerStateTemp = await teamStateLib.playerStateCreate(
            defence,
            speed,
            pass,
            shoot,
            endurance,
            0, 
            playerId = '1',
            0, 0, 0, 0, 0, 0
        ).should.be.fulfilled;
        teamStateTemp = await teamStateLib.teamStateCreate().should.be.fulfilled;
        for (var i = 0; i < PLAYERS_PER_TEAM; i++) {
            teamStateTemp = await teamStateLib.teamStateAppend(teamStateTemp, playerStateTemp).should.be.fulfilled;
        }
        return teamStateTemp;
    };
    

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        states = await States.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, states.address).should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        teamStateAll1 = await createTeamStateFromSinglePlayer(1,1,1,1,1,states);
        teamStateAll50 = await createTeamStateFromSinglePlayer(50,50,50,50,50,states);
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 1, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 2, order, tactic442).should.be.fulfilled;
        
        leagueState = await states.leagueStateCreate().should.be.fulfilled;
        leagueState = await states.leagueStateAppend(leagueState, teamStateAll1).should.be.fulfilled;
        leagueState = await states.leagueStateAppend(leagueState, teamStateAll50).should.be.fulfilled;
    });

    it('compute match', async () => {
        const result = await leagues.computeMatch(teamStateAll1, teamStateAll50, tacticsIds, 0).should.be.fulfilled;
        const scores = await leagues.decodeScore(result.score).should.be.fulfilled;
        scores.home.toString().should.be.equal('0');
        scores.visitor.toString().should.be.equal('16');
        let valid = await states.isValidTeamState(result.newHomeState).should.be.fulfilled;
        valid.should.be.equal(true);
        valid = await states.isValidTeamState(result.newVisitorState).should.be.fulfilled;
        valid.should.be.equal(true);
    });
    
    it('Engine contract', async () => {
        const address = await leagues.getEngineContract().should.be.fulfilled;
        address.should.be.equal(engine.address);
    });

    it('evolve team', async () => {
        const result = await leagues.evolveTeams(teamStateAll1, teamStateAll50, 0, 2).should.be.fulfilled;
        let valid = await states.isValidTeamState(result.updatedHomeTeamState).should.be.fulfilled;
        valid.should.be.equal(true);
        valid = await states.isValidTeamState(result.updatedVisitorTeamState).should.be.fulfilled;
        valid.should.be.equal(true);
    })

    it('calculate a day in a league', async () => {
        let day = 0;
        let result = await leagues.computeDayWithSeed(leagueId, day, leagueState, tacticsIds, 0).should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(0);
        scores.visitor.toNumber().should.be.equal(16);
        day = 1;
        result = await leagues.computeDayWithSeed(leagueId, day, leagueState, tacticsIds, 4354646451).should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(18);
        scores.visitor.toNumber().should.be.equal(0);
    });

    it('result of a day in league is deterministic', async () => {
        const day = 1;
        let result = await leagues.computeDayWithSeed(leagueId, day, leagueState, tacticsIds, 123456).should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        let scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(14);
        scores.visitor.toNumber().should.be.equal(0);
        result = await leagues.computeDayWithSeed(leagueId, day, leagueState, tacticsIds, 123456).should.be.fulfilled;
        result.scores.length.should.be.equal(1);
        scores = await leagues.decodeScore(result.scores[0]).should.be.fulfilled;
        scores.home.toNumber().should.be.equal(14);
        scores.visitor.toNumber().should.be.equal(0);
    });


    it('different seed => different results', async () => {
        const day = 1;
        const result0 = await leagues.computeDayWithSeed(leagueId, day, leagueState, tacticsIds, 1234).should.be.fulfilled;
        const result1 = await leagues.computeDayWithSeed(leagueId, day, leagueState, tacticsIds, 4354646451).should.be.fulfilled;
        result0.scores.length.should.be.equal(1);
        result1.scores.length.should.be.equal(1);
        result0.scores[0].toNumber().should.not.be.equal(result1.scores[0].toNumber());
    });

   it('compute points same rating', async () => {
        let points = await leagues.computePoints(teamStateAll50, teamStateAll50, 2, 2).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(0);
        points.visitorPoints.toNumber().should.be.equal(0);
        points = await leagues.computePoints(teamStateAll50, teamStateAll50, 2, 1).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(5);
        points.visitorPoints.toNumber().should.be.equal(0);
        points = await leagues.computePoints(teamStateAll50, teamStateAll50, 1, 2).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(0);
        points.visitorPoints.toNumber().should.be.equal(5);
    });

    it('hash init state', async () => {
        const playerState = await states.playerStateCreate(1, 2, 3, 4, 5, 0, playerId = 1, 0, 0, 0, 0, 0, 0).should.be.fulfilled;
        let teamState = await states.teamStateCreate().should.be.fulfilled;
        teamState = await states.teamStateAppend(teamState, playerState).should.be.fulfilled;
        let state = await states.leagueStateCreate().should.be.fulfilled;
        state = await states.leagueStateAppend(state, teamState).should.be.fulfilled;
        const hash0 = await leagues.hashInitState(state).should.be.fulfilled;
        hash0.should.be.equal('0x3314e33d03ecb18ceb0a2f84d6a7186c252a35f2bb7ce025002927dcfcd9b21d');
        state = await states.leagueStateAppend(state, teamState).should.be.fulfilled;
        const hash1 = await leagues.hashInitState(state).should.be.fulfilled;
        hash1.should.be.not.equal(hash0);
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