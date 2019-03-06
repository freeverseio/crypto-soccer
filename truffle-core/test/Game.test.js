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
        let leagueState = await state.leagueStateCreate().should.be.fulfilled;
        leagueState = await state.leagueStateAppend(leagueState, barcelonaState).should.be.fulfilled;
        leagueState = await state.leagueStateAppend(leagueState, madridState).should.be.fulfilled;
        leagueState = await state.leagueStateAppend(leagueState, sevillaState).should.be.fulfilled;
        leagueState = await state.leagueStateAppend(leagueState, bilbaoState).should.be.fulfilled;
        leagueState = await state.leagueStateAppend(leagueState, veniceState).should.be.fulfilled;
        leagueState = await state.leagueStateAppend(leagueState, juventusState).should.be.fulfilled;

        // generate init league state hash
        const initStateHash = await leagues.hashState(leagueState).should.be.fulfilled;
        let dayStateHashes = [];

        let leagueScores = await leagues.scoresCreate().should.be.fulfilled;

        // get days in a league
        const leagueDays = await leagues.countLeagueDays(leagueId).should.be.fulfilled;
        leagueDays.toNumber().should.be.equal(10);
        const matchesPerDay = await leagues.getMatchPerDay(leagueId).should.be.fulfilled;
        matchesPerDay.toNumber().should.be.equal(3);

        // compute result for each league day
        for (leagueDay = 0; leagueDay < leagueDays.toNumber(); leagueDay++) {
            // compute result for league day
            const result = await leagues.computeDay(leagueId, leagueDay, leagueState, tactics).should.be.fulfilled;
            const dayScores = result.scores;

            // concat day scores to league scores
            leagueScores = await leagues.scoresConcat(leagueScores, dayScores).should.be.fulfilled;

            // update the league state with the updated league state
            leagueState = result.finalLeagueState;

            // for each match of the day
            for (match = 0; match < matchesPerDay; match++) {
                 // get the indexes of the teams of match 
                 const teamsInMatch = await leagues.getTeamsInMatch(leagueId, leagueDay, match).should.be.fulfilled;

                 // get ids of teams in match
                 const homeTeamId = teamIds[teamsInMatch.homeIdx.toNumber()];
                 const visitorTeamId = teamIds[teamsInMatch.visitorIdx.toNumber()];

                 // get the names of the teams in the match
                 const homeTeam = await teams.getName(homeTeamId).should.be.fulfilled;
                 const visitorTeam = await teams.getName(visitorTeamId).should.be.fulfilled;

                 // getting the goal of match 
                 const goals = await leagues.decodeScore(dayScores[match]).should.be.fulfilled;

                 // get the state of teams 
                 const homeTeamState = await state.leagueStateAt(leagueState, teamsInMatch.homeIdx).should.be.fulfilled;
                 const visitorTeamState = await state.leagueStateAt(leagueState, teamsInMatch.visitorIdx).should.be.fulfilled;

                 // calculate rating of teams
                 const homeTeamRating = await state.computeTeamRating(homeTeamState).should.be.fulfilled;
                 const visitorTeamRating = await state.computeTeamRating(visitorTeamState).should.be.fulfilled;

                 console.log(
                     "DAY " + leagueDay + ": " 
                     + homeTeam + "(" + homeTeamRating.toNumber() + ") " + goals.home.toNumber() 
                     + " - " 
                     + goals.visitor.toNumber() + " " + visitorTeam + "(" + visitorTeamRating + ")");
            }

            // hash of the day state
            const dayHash = await leagues.hashState(leagueState).should.be.fulfilled;

            // append the day state hash
            dayStateHashes.push(dayHash);
        }

        dayStateHashes.length.should.be.equal(leagueDays.toNumber());

        // updating the league
        await leagues.updateLeague(
            leagueId,
            initStateHash,
            dayStateHashes,
            leagueScores
        ).should.be.fulfilled;
    });
})
