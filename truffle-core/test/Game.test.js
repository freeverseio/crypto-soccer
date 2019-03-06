require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Horizon = artifacts.require('Horizon');
const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');
const State = artifacts.require('LeagueState');



contract('Game', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;
    let engine = null;
    let state = null;
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

        engine = await Engine.new().should.be.fulfilled;
        state = await State.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, state.address).should.be.fulfilled;

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
        let teamState = await state.teamStateCreate().should.be.fulfilled;
        const playersIds = await teams.getPlayers(id).should.be.fulfilled;
        for (let i = 0; i < playersIds.length; i++) {
            const playerId = playersIds[i];
            const defence = await players.getDefence(playerId).should.be.fulfilled;
            const speed = await players.getSpeed(playerId).should.be.fulfilled;
            const pass = await players.getPass(playerId).should.be.fulfilled;
            const shoot = await players.getShoot(playerId).should.be.fulfilled;
            const endurance = await players.getEndurance(playerId).should.be.fulfilled;
            const playerState = await state.playerStateCreate(
                defence,
                speed,
                pass,
                shoot,
                endurance
            );
            teamState = await state.teamStateAppend(teamState, playerState).should.be.fulfilled;
        }
        return teamState;
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
        let dayState = await state.dayStateCreate().should.be.fulfilled;
        dayState = await state.dayStateAppend(dayState, barcelonaState).should.be.fulfilled;
        dayState = await state.dayStateAppend(dayState, madridState).should.be.fulfilled;
        dayState = await state.dayStateAppend(dayState, sevillaState).should.be.fulfilled;
        dayState = await state.dayStateAppend(dayState, bilbaoState).should.be.fulfilled;
        dayState = await state.dayStateAppend(dayState, veniceState).should.be.fulfilled;
        dayState = await state.dayStateAppend(dayState, juventusState).should.be.fulfilled;

        // get days in a league
        const leagueDays = await leagues.countLeagueDays(leagueId).should.be.fulfilled;
        leagueDays.toNumber().should.be.equal(10);
        const matchesPerDay = await leagues.getMatchPerDay(leagueId).should.be.fulfilled;
        matchesPerDay.toNumber().should.be.equal(3);

        // compute result for each league day
        for (leagueDay = 0; leagueDay < leagueDays.toNumber(); leagueDay++) {
            // compute result for league day
            const result = await leagues.computeDay(leagueId, leagueDay, dayState, tactics).should.be.fulfilled;
            const dayScores = result.scores;
            for (match = 0; match < matchesPerDay; match++) {
                 // get the indexes of the teams of match j
                 const teamsInMatch = await leagues.getTeamsInMatch(leagueId, leagueDay, match).should.be.fulfilled;

                 // get the names of the teams in the match
                 const homeTeam = await teams.getName(teamIds[teamsInMatch.homeIdx.toNumber()]).should.be.fulfilled;
                 const visitorTeam = await teams.getName(teamIds[teamsInMatch.visitorIdx.toNumber()]).should.be.fulfilled;

                 // getting the goal of match 
                 const goals = await leagues.decodeScore(dayScores[match]).should.be.fulfilled;
                 const homeTeamState = await state.dayStateAt(dayState, teamsInMatch.homeIdx).should.be.fulfilled;
                 const homeTeamRating = await state.computeTeamRating(homeTeamState).should.be.fulfilled;
                 const visitorTeamState = await state.dayStateAt(dayState, teamsInMatch.visitorIdx).should.be.fulfilled;
                 const visitorTeamRating = await state.computeTeamRating(visitorTeamState).should.be.fulfilled;

                 console.log(
                     "DAY " + leagueDay + ": " 
                     + homeTeam + "(" + homeTeamRating.toNumber() + ") " + goals.home.toNumber() 
                     + " - " 
                     + goals.visitor.toNumber() + " " + visitorTeam + "(" + visitorTeamRating + ")");
            }
            dayState = result.finalDayState;
        }
        return;
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
