require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Horizon = artifacts.require('Horizon');
const LeagueState = artifacts.require('LeagueState');
const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');



contract('Game', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;
    let engine = null;
    let stateLib = null;
    let leagues = null;
    let barcelonaId = null;
    let madridId = null;
    let sevillaId = null;
    let bilbaoId = null;
    let veniceId = null;
    let juventusId = null;

    beforeEach(async () => {
        players = await Players.new().should.be.fulfilled;
        teams = await Teams.new(players.address).should.be.fulfilled;
        horizon = await Horizon.new(teams.address).should.be.fulfilled;
        await players.addMinter(horizon.address).should.be.fulfilled;
        await players.renounceMinter().should.be.fulfilled;
        await teams.addMinter(horizon.address).should.be.fulfilled;
        await teams.renounceMinter().should.be.fulfilled;
        await players.addTeamsContract(teams.address).should.be.fulfilled;
        await players.renounceTeamsContract().should.be.fulfilled;

        stateLib = await LeagueState.new().should.be.fulfilled;
        Leagues.link("LeagueState", stateLib.address);
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address).should.be.fulfilled;

        await horizon.createTeam("Barcelona").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        await horizon.createTeam("Sevilla").should.be.fulfilled;
        await horizon.createTeam("Bilbao").should.be.fulfilled;
        await horizon.createTeam("Venice").should.be.fulfilled;
        await horizon.createTeam("Juventus").should.be.fulfilled;
        barcelonaId = await teams.getTeamId("Barcelona").should.be.fulfilled;
        madridId = await teams.getTeamId("Madrid").should.be.fulfilled;
        sevillaId = await teams.getTeamId("Sevilla").should.be.fulfilled;
        bilbaoId = await teams.getTeamId("Bilbao").should.be.fulfilled;
        veniceId = await teams.getTeamId("Venice").should.be.fulfilled;
        juventusId = await teams.getTeamId("Juventus").should.be.fulfilled;
    });

    // we use the values in the blockchain to generate the team status
    // it will use a local DBMS in the final version
    const generateTeamState = async (id) => {
        let state = [];
        const playersIds = await teams.getPlayers(id).should.be.fulfilled;
        for (let i = 0; i < playersIds.length; i++) {
            const playerId = playersIds[i];
            const genome = await players.getGenome(playerId).should.be.fulfilled;
            state.push(genome);
        }
        return state;
    }

    it('play a league of 6 teams', async () => {
        const initBlock = 1;
        const step = 1;
        const leagueId = 0;
        const teamIds = [barcelonaId, madridId, sevillaId, bilbaoId, veniceId, juventusId];
        const tactics = [[4,4,3], [4,4,3], [4,4,3], [4,4,3], [4,4,3], [4,4,3]];
        await leagues.create(leagueId, initBlock, step, teamIds).should.be.fulfilled;

        const barcelonaState = await generateTeamState(barcelonaId).should.be.fulfilled;
        const madridState = await generateTeamState(madridId).should.be.fulfilled;
        const sevillaState = await generateTeamState(sevillaId).should.be.fulfilled;
        const bilbaoState = await generateTeamState(bilbaoId).should.be.fulfilled;
        const veniceState = await generateTeamState(veniceId).should.be.fulfilled;
        const juventusState = await generateTeamState(juventusId).should.be.fulfilled;

        // we build the league state
        let leagueState = await stateLib.append(barcelonaState, madridState).should.be.fulfilled;
        leagueState = await stateLib.append(leagueState, sevillaState).should.be.fulfilled;
        leagueState = await stateLib.append(leagueState, bilbaoState).should.be.fulfilled;
        leagueState = await stateLib.append(leagueState, veniceState).should.be.fulfilled;
        leagueState = await stateLib.append(leagueState, juventusState).should.be.fulfilled;

        // calculate the league
        const leagueScores = await leagues.computeAllMatchdayStates(leagueId, leagueState, tactics).should.be.fulfilled;

        // get the number of days in the league
        const nDayScores = await leagues.countLeagueDays(leagueId).should.be.fulfilled;
        nDayScores.toNumber().should.be.equal(10);

        // for each day we get the scores
        for (i = 0; i < nDayScores.toNumber(); i++) {
            const dayScores = await leagues.getDayScores(leagueScores, i).should.be.fulfilled;
            // 2 matches per day
            dayScores.length.should.be.equal(3);
            for (j = 0; j < dayScores.length; j++) {
                // get the indexes of the teams of match j
                const teamsInMatch = await leagues.getTeamsInMatch(leagueId, i, j).should.be.fulfilled;

                // get the names of the teams in the match
                const homeTeam = await teams.getName(teamIds[teamsInMatch.homeIdx.toNumber()]).should.be.fulfilled;
                const visitorTeam = await teams.getName(teamIds[teamsInMatch.visitorIdx.toNumber()]).should.be.fulfilled;

                // getting the goal of match j
                const goals = await leagues.decodeScore(dayScores[j]).should.be.fulfilled;

                console.log("DAY " + i + ": " + homeTeam + " " + goals.home.toNumber() + " - " + goals.visitor.toNumber() + " " + visitorTeam);
            }
        }

        // generate init league state hash
        const initStateHash = await leagues.hashInitState(leagueState).should.be.fulfilled;

        // generate the final state for each team
        // n.b. we use init league state cause the game doesn't evolve the teams yet
        const finalTeamsStateHashes = await leagues.hashLeagueState(leagueState).should.be.fulfilled;

        // updating the league
        await leagues.updateLeague(
            leagueId,
            initStateHash,
            finalTeamsStateHashes,
            leagueScores
        ).should.be.fulfilled;
    });
})
