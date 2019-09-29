const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Cups = artifacts.require('Cups');
const Engine = artifacts.require('Engine');

contract('Cups', (accounts) => {
    const tactic442 = 0;
    const tactic433 = 1;

    const createTeamStateFromSinglePlayer = async (skills, engine) => {
        const playerStateTemp = await engine.encodePlayerSkills(
            skills, 
            monthOfBirth = 0, 
            playerId = 1321312,
            [potential = 3,
            forwardness = 3,
            leftishness = 2,
            aggressiveness = 0],
            alignedLastHalf = false,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0
        ).should.be.fulfilled;

        teamState = []
        for (player = 0; player < PLAYERS_PER_TEAM_MAX; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };
    
    const createLeagueStateFromSinglePlayer = async (skills, engine) => {
        const teamState = await createTeamStateFromSinglePlayer(skills, engine).should.be.fulfilled;
        leagueState = []
        for (team = 0; team < TEAMS_PER_LEAGUE; team++) {
            leagueState.push(teamState)
        }
        return leagueState;
    };
    

    beforeEach(async () => {
        cups = await Cups.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        await cups.setEngineAdress(engine.address).should.be.fulfilled;
        TEAMS_PER_LEAGUE = await cups.TEAMS_PER_LEAGUE().should.be.fulfilled;
        PLAYERS_PER_TEAM_MAX = await cups.PLAYERS_PER_TEAM_MAX().should.be.fulfilled;
        MATCHDAYS = await cups.MATCHDAYS().should.be.fulfilled;
        MATCHES_PER_DAY = await cups.MATCHES_PER_DAY().should.be.fulfilled;
        teamStateAll50 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine);
        teamStateAll1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine);
    });

    it('get teams for groups', async () => {
        teams = await cups.getTeamsInGroup(groupIdx = 0).should.be.fulfilled;
        teamsExpected = [ 0, 8, 16, 24, 32, 40, 48, 56 ]
        for (t = 0; t < teams.length; t++) {
            teamsExpected.push(teams[t].toNumber())
            teams[t].toNumber().should.be.equal(teamsExpected[t])
        }
        console.log(teamsExpected)
    });

});