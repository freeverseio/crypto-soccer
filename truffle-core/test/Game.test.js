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
        barcelonaId = await teams.getTeamId("Barcelona").should.be.fulfilled;
        madridId = await teams.getTeamId("Madrid").should.be.fulfilled;
        sevillaId = await teams.getTeamId("Sevilla").should.be.fulfilled;
        bilbaoId = await teams.getTeamId("Bilbao").should.be.fulfilled;
    });

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

    it('play a league of 2 teams', async () => {
        const initBlock = 1;
        const step = 1;
        const leagueId = 0;
        const teamIds = [barcelonaId, madridId, sevillaId, bilbaoId];
        await leagues.create(leagueId, initBlock, step, teamIds).should.be.fulfilled;

        const barcelonaState = await generateTeamState(barcelonaId).should.be.fulfilled;
        const madridState = await generateTeamState(madridId).should.be.fulfilled;
        const sevillaState = await generateTeamState(sevillaId).should.be.fulfilled;
        const bilbaoState = await generateTeamState(bilbaoId).should.be.fulfilled;

        let leagueState = await stateLib.append(barcelonaState, madridState).should.be.fulfilled;
        leagueState = await stateLib.append(leagueState, sevillaState).should.be.fulfilled;
        leagueState = await stateLib.append(leagueState, bilbaoState).should.be.fulfilled;
        const leagueScores = await leagues.computeAllMatchdayStates(leagueId, leagueState, [[4,4,3], [4,4,3], [4,4,3], [4,4,3]]).should.be.fulfilled;
        const nDayScores = await leagues.countDaysInTournamentScores(leagueScores).should.be.fulfilled;
        nDayScores.toNumber().should.be.equal(6);
        for (i = 0; i < nDayScores.toNumber(); i++) {
            const dayScores = await leagues.getDayScores(leagueScores, i).should.be.fulfilled;
            // 2 matches per day
            dayScores.length.should.be.equal(2);
            for (j = 0; j < 2; j++) {
                // get the indexes of the teams of match j
                const teamsInMatch = await leagues.getTeamsInMatch(leagueId, i, j).should.be.fulfilled;
                const goals = await leagues.decodeScore(dayScores[0]).should.be.fulfilled;
                const homeTeam = await teams.getName(teamIds[teamsInMatch.homeIdx.toNumber()]).should.be.fulfilled;
                const visitorTeam = await teams.getName(teamIds[teamsInMatch.visitorIdx.toNumber()]).should.be.fulfilled;
                console.log("DAY " + i + ": " + homeTeam + " " + goals.home.toNumber() + " - " + goals.visitor.toNumber() + " " + visitorTeam)
            }
        }

        const initStateHash = await leagues.hashInitState(leagueState).should.be.fulfilled;
        const finalTeamsStateHashes = await leagues.hashLeagueState(leagueState).should.be.fulfilled;

        await leagues.updateLeague(
            leagueId,
            initStateHash,
            finalTeamsStateHashes,
            leagueScores
        ).should.be.fulfilled;
    });
})
