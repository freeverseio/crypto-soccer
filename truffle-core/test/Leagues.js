const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');

contract('Leagues', (accounts) => {
    let leagues = null;
    let engine = null;
    let TEAMS_PER_LEAGUE = null;
    let MATCHDAYS = null;
    let MATCHES_PER_DAY = null;
    let teamStateAll50 = null;
    let teamStateAll1 = null;

    const createTeamStateFromSinglePlayer = async (skills, engine) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1
        ).should.be.fulfilled;

        teamState = []
        for (player = 0; player < 11; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };
    
    beforeEach(async () => {
        leagues = await Leagues.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        await leagues.init().should.be.fulfilled;
        await leagues.setEngineAdress(engine.address).should.be.fulfilled;
        TEAMS_PER_LEAGUE = await leagues.TEAMS_PER_LEAGUE().should.be.fulfilled;
        MATCHDAYS = await leagues.MATCHDAYS().should.be.fulfilled;
        MATCHES_PER_DAY = await leagues.MATCHES_PER_DAY().should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine);
    });

    // it('check initial constants', async () =>  {
    //     engine = 0;
    //     MATCHDAYS.toNumber().should.be.equal(14);
    //     MATCHES_PER_DAY.toNumber().should.be.equal(4);
    //     TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    // });

    // it('get teams for match in wrong day', async () => {
    //     matchIdxInDay = 0; 
    //     day = MATCHDAYS-1; 
    //     await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
    //     day = MATCHDAYS; 
    //     await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    // });

    // it('get teams for match in wrong match in day', async () => {
    //     day = 0;
    //     matchIdxInDay = MATCHES_PER_DAY-1;
    //     await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
    //     matchIdxInDay = MATCHES_PER_DAY;
    //     await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.rejected;
    // });

    // it('get teams for match in league day', async () => {
    //     day = 0;
    //     matchIdxInDay = 0;
    //     teams = await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
    //     teams[0].toNumber().should.be.equal(0);
    //     teams[1].toNumber().should.be.equal(1);
    //     day = Math.floor(MATCHDAYS/2);
    //     teams = await leagues.getTeamsInMatch(day, matchIdxInDay).should.be.fulfilled;
    //     teams[0].toNumber().should.be.equal(1);
    //     teams[1].toNumber().should.be.equal(0);
    // });
    
    it('compute points same rating', async () => {
        let points = await leagues.computeEvolutionPoints(teamStateAll50, teamStateAll50, 2, 2).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(0);
        points.visitorPoints.toNumber().should.be.equal(0);
        points = await leagues.computeEvolutionPoints(teamStateAll50, teamStateAll50, 2, 1).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(5);
        points.visitorPoints.toNumber().should.be.equal(0);
        points = await leagues.computeEvolutionPoints(teamStateAll50, teamStateAll50, 1, 2).should.be.fulfilled;
        points.homePoints.toNumber().should.be.equal(0);
        points.visitorPoints.toNumber().should.be.equal(5);
    });


});