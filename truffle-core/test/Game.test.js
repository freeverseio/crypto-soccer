require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Horizon = artifacts.require('Horizon');
const Leagues = artifacts.require('LeaguesComputer');
const Engine = artifacts.require('Engine');

contract('Game', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;
    let engine = null;
    let leagues = null;

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
        leagues = await Leagues.new(engine.address).should.be.fulfilled;
    });

    it('play a league of 2 teams', async () => {
        await horizon.createTeam("Barcelona").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        const barcelonaId = await teams.getTeamId("Barcelona").should.be.fulfilled;
        const madridId = await teams.getTeamId("Madrid").should.be.fulfilled;
        const blockInitDelta = 0;
        const step = 1;
        const leagueId = 0;
        await leagues.create(leagueId, blockInitDelta, step, [barcelonaId, madridId]).should.be.fulfilled;
        const barcelonaPlayersIds = await teams.getPlayers(barcelonaId).should.be.fulfilled;
        const madridPlayersIds = await teams.getPlayers(madridId).should.be.fulfilled;

        let barcelonaState = [];
        for (let i = 0; i < barcelonaPlayersIds.length ; i++){
            const playerId = barcelonaPlayersIds[i];
            const genome = await players.getGenome(playerId).should.be.fulfilled;
            barcelonaState.push(genome);
        }

        let madridState = [];
        for (let i = 0; i < madridPlayersIds.length ; i++){
            const playerId = madridPlayersIds[i];
            const genome = await players.getGenome(playerId).should.be.fulfilled;
            madridState.push(genome);
        }

        let leagueState = await leagues.appendTeamToLeagueState([], barcelonaState).should.be.fulfilled;
        leagueState = await leagues.appendTeamToLeagueState(leagueState, madridState).should.be.fulfilled;
        const scores = await leagues.computeLeagueFinalState(leagueId, leagueState, [[4,4,3], [4,4,3]]).should.be.fulfilled;
        console.log("Barcelona - Madrid: " + scores[0][0].toNumber() + " - " + scores[0][1].toNumber());
        console.log("Madrid - Barcelona: " + scores[1][0].toNumber() + " - " + scores[1][1].toNumber());

        const initStateHash = await leagues.hashInitState(leagueState).should.be.fulfilled;
        const finalTeamsStateHashes = await leagues.hashLeagueState(leagueState).should.be.fulfilled;

        await leagues.updateLeague(
            leagueId,
            initStateHash,
            finalTeamsStateHashes,
            scores
        ).should.be.fulfilled;

        const recordedInitStateHash = await leagues.getInitStateHash(leagueId).should.be.fulfilled;
        recordedInitStateHash.should.be.equal(initStateHash);
        const recordedFinalTeamStateHashes = await leagues.getFinalTeamStateHashes(leagueId).should.be.fulfilled;
        for (let i = 0; i < finalTeamsStateHashes.length; i++) {
            recordedFinalTeamStateHashes[i].should.be.equal(finalTeamsStateHashes[i]);
        }
    });
})
