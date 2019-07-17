require('chai')
    .use(require('chai-as-promised'))
    .should();

const Assets = artifacts.require('Assets');
const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');
const State = artifacts.require('LeagueState');

contract('Game', (accounts) => {
    let engine = null;
    let state = null;
    let leagues = null;
    let cronos = null;
    const initBlock = 1;
    const step = 1;
    const leagueId = 1;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 0;

    beforeEach(async () => {
        state = await State.new().should.be.fulfilled;
        assets = await Assets.new(state.address).should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, state.address).should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
    });

    const createTeam = async (name, owner) => {
        let receipt = await assets.createTeam(name, owner).should.be.fulfilled;
        const teamId = receipt.logs[0].args.id.toNumber()
        return teamId;
    }

    // we use the values in the blockchain to generate the team status
    // it will use a local DBMS in the final version
    const generateTeamState = async (id) => {
        let teamState = await state.teamStateCreate().should.be.fulfilled;
        const playersIds = await assets.getTeamPlayerIds(id).should.be.fulfilled;
        for (let i = 0; i < playersIds.length; i++) {
            const playerState = await assets.getPlayerState(playersIds[i]).should.be.fulfilled;
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

    const advanceNBlocks = async (blocks) => {
        let current = await web3.eth.getBlockNumber().should.be.fulfilled;
        await advanceToBlock(current + blocks);
    }

    const prepareMatchdayHashes = async (statesAtMatchday) => {
        let result = [];
        for (let i = 0; i < statesAtMatchday.length; i++) {
            const state = statesAtMatchday[i];
            const hash = await leagues.hashState(state).should.be.fulfilled;
            result.push(hash);
        }

        return result;
    }

    it('play a league of 6 teams', async () => {
        const barcelonaId = await createTeam("Barcelona", accounts[0]).should.be.fulfilled;
        const madridId = await createTeam("Madrid", accounts[0]).should.be.fulfilled;
        const sevillaId = await createTeam("Sevilla", accounts[0]).should.be.fulfilled;
        const bilbaoId = await createTeam("Bilbao", accounts[0]).should.be.fulfilled;
        const veniceId = await createTeam("Venice", accounts[0]).should.be.fulfilled;
        const juventusId = await createTeam("Juventus", accounts[0]).should.be.fulfilled;

        const tactics = [tactic442, tactic442, tactic442, tactic442, tactic442, tactic442];
        
        const teamIds = [barcelonaId, madridId, sevillaId, bilbaoId, veniceId, juventusId];

        await leagues.create(nTeams = 6, initBlock, step).should.be.fulfilled;
        for (var team = 0; team < 6; team++) {
            await leagues.signTeamInLeague(leagueId, teamIds[team], order, tactics[team]).should.be.fulfilled;
        }
        
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
        const initStateHash = await leagues.hashInitState(initLeagueState).should.be.fulfilled;
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
                const homeTeam = await assets.getTeamName(homeTeamId).should.be.fulfilled;
                const visitorTeam = await assets.getTeamName(visitorTeamId).should.be.fulfilled;

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
            const dayHash = await leagues.hashDayState(finalDayState).should.be.fulfilled;
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
            leagueScores,
            isLie = false
        ).should.be.fulfilled;
    });
})
