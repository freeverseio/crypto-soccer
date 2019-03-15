require('chai')
    .use(require('chai-as-promised'))
    .should();

const Players = artifacts.require('Players');
const Teams = artifacts.require('Teams');
const Horizon = artifacts.require('Horizon');
const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');
const State = artifacts.require('LeagueState');
const Cronos = artifacts.require('Cronos');

contract('Game', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;
    let engine = null;
    let state = null;
    let leagues = null;
    let cronos = null;

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
        cronos = await Cronos.new().should.be.fulfilled;
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
                endurance,
                0,
                1, // TODO: make it properly
                0, 0, 0, 0, 0, 0
            );
            teamState = await state.teamStateAppend(teamState, playerState).should.be.fulfilled;
        }
        return teamState;
    }

    const advanceToBlock = async (block) => {
        let current = await web3.eth.getBlockNumber().should.be.fulfilled;
        while (current.toString() < block) {
            await cronos.wait().should.be.fulfilled;
            current = await web3.eth.getBlockNumber().should.be.fulfilled;
        }
    }

    it('test2', async () => {
        await horizon.createTeam("Barca").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        await horizon.createTeam("Milan").should.be.fulfilled;
        await horizon.createTeam("PSG").should.be.fulfilled;
        const teamIdx1 = await teams.getTeamId("Barca").should.be.fulfilled;
        const teamIdx2 = await teams.getTeamId("Madrid").should.be.fulfilled;
        const teamIdx3 = await teams.getTeamId("Milan").should.be.fulfilled;
        const teamIdx4 = await teams.getTeamId("PSG").should.be.fulfilled;

        await advanceToBlock(100);

        const blockInit = 190;
        const blockStep = 10;

        const usersInitData = {
            teamIdxs: [teamIdx1, teamIdx2],
            // teamOrders: [DEFAULT_ORDER, REVERSE_ORDER],
            tactics: [[4,4,2], [4,3,3]] 
        };

        leagueIdx = 0;
        await leagues.create(leagueIdx, blockInit, blockStep, usersInitData.teamIdxs, usersInitData.tactics).should.be.fulfilled;

        const startBlock = await leagues.getInitBlock(leagueIdx).should.be.fulfilled;
        startBlock.toNumber().should.be.equal(blockInit);

        // Advance to matchday 2
        await advanceToBlock(blockInit + blockStep - 5);
        const started = await leagues.hasStarted(leagueIdx).should.be.fulfilled;
        started.should.be.equal(true);
        const finished = await leagues.hasFinished(leagueIdx).should.be.fulfilled;
        finished.should.be.equal(false);

        // Note that we could specify only for 1 of the teams if we wanted.
        usersAlongData = {
            teamIdxsWithinLeague: [teamIdx1, teamIdx2],
            tactics: [[4, 3, 3], [4, 4, 2]],
        };

        // Submit data to change tactics
        await leagues.updateUsersAlongDataHash(leagueIdx, usersAlongData.teamIdxsWithinLeague[0], usersAlongData.tactics[0]).should.be.fulfilled;
        await leagues.updateUsersAlongDataHash(leagueIdx, usersAlongData.teamIdxsWithinLeague[1], usersAlongData.tactics[1]).should.be.fulfilled;
    });

    return;

    it('play a league of 6 teams', async () => {
        await horizon.createTeam("Barcelona").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        await horizon.createTeam("Sevilla").should.be.fulfilled;
        await horizon.createTeam("Bilbao").should.be.fulfilled;
        await horizon.createTeam("Venice").should.be.fulfilled;
        await horizon.createTeam("Juventus").should.be.fulfilled;
        const barcelonaId = await teams.getTeamId("Barcelona").should.be.fulfilled;
        const madridId = await teams.getTeamId("Madrid").should.be.fulfilled;
        const sevillaId = await teams.getTeamId("Sevilla").should.be.fulfilled;
        const bilbaoId = await teams.getTeamId("Bilbao").should.be.fulfilled;
        const veniceId = await teams.getTeamId("Venice").should.be.fulfilled;
        const juventusId = await teams.getTeamId("Juventus").should.be.fulfilled;


        const initBlock = 1;
        const step = 1;
        const leagueId = 0;
        const teamIds = [barcelonaId, madridId, sevillaId, bilbaoId, veniceId, juventusId];
        const tactics = [[4, 4, 3], [4, 4, 3], [4, 4, 3], [4, 4, 3], [4, 4, 3], [4, 4, 3]];
        await leagues.create(leagueId, initBlock, step, teamIds).should.be.fulfilled;

        const barcelonaState = await generateTeamState(barcelonaId).should.be.fulfilled;
        const madridState = await generateTeamState(madridId).should.be.fulfilled;
        const sevillaState = await generateTeamState(sevillaId).should.be.fulfilled;
        const bilbaoState = await generateTeamState(bilbaoId).should.be.fulfilled;
        const veniceState = await generateTeamState(veniceId).should.be.fulfilled;
        const juventusState = await generateTeamState(juventusId).should.be.fulfilled;

        // we build the league state
        let initLeagueState = await state.leagueStateCreate().should.be.fulfilled;
        initLeagueState = await state.leagueStateAppend(initLeagueState, barcelonaState).should.be.fulfilled;
        initLeagueState = await state.leagueStateAppend(initLeagueState, madridState).should.be.fulfilled;
        initLeagueState = await state.leagueStateAppend(initLeagueState, sevillaState).should.be.fulfilled;
        initLeagueState = await state.leagueStateAppend(initLeagueState, bilbaoState).should.be.fulfilled;
        initLeagueState = await state.leagueStateAppend(initLeagueState, veniceState).should.be.fulfilled;
        initLeagueState = await state.leagueStateAppend(initLeagueState, juventusState).should.be.fulfilled;

        // generate init league state hash
        const initStateHash = await leagues.hashState(initLeagueState).should.be.fulfilled;
        let dayStateHashes = [];

        let leagueScores = await leagues.scoresCreate().should.be.fulfilled;

        // get days in a league
        const leagueDays = await leagues.countLeagueDays(leagueId).should.be.fulfilled;
        leagueDays.toNumber().should.be.equal(10);
        const matchesPerDay = await leagues.getMatchPerDay(leagueId).should.be.fulfilled;
        matchesPerDay.toNumber().should.be.equal(3);

        // compute result for each league day
        let initDayState = initLeagueState;
        for (leagueDay = 0; leagueDay < leagueDays.toNumber(); leagueDay++) {
            // compute result for league day
            const result = await leagues.computeDay(leagueId, leagueDay, initDayState, tactics).should.be.fulfilled;
            const dayScores = result.scores;
            const finalDayState = result.finalLeagueState;

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
                const homeTeamState = await state.leagueStateAt(finalDayState, teamsInMatch.homeIdx).should.be.fulfilled;
                const visitorTeamState = await state.leagueStateAt(finalDayState, teamsInMatch.visitorIdx).should.be.fulfilled;

                // calculate rating of teams
                const homeTeamRating = await state.computeTeamRating(homeTeamState).should.be.fulfilled;
                const visitorTeamRating = await state.computeTeamRating(visitorTeamState).should.be.fulfilled;

                console.log(
                    "DAY " + leagueDay + ": "
                    + homeTeam + "(" + homeTeamRating.toNumber() + ") " + goals.home.toNumber()
                    + " - "
                    + goals.visitor.toNumber() + " " + visitorTeam + "(" + visitorTeamRating + ")");
            }
            // concat day scores to league scores
            leagueScores = await leagues.scoresConcat(leagueScores, dayScores).should.be.fulfilled;

            // hash of the day state
            const dayHash = await leagues.hashState(finalDayState).should.be.fulfilled;
            // append the day state hash
            dayStateHashes.push(dayHash);

            initDayState = finalDayState;
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
