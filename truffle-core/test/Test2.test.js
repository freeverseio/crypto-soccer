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

contract('Test2', (accounts) => {
    let horizon = null;
    let players = null;
    let teams = null;
    let engine = null;
    let state = null;
    let leagues = null;
    let cronos = null;
    let CHALLENGING_PERIOD_BLKS = null;

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
        CHALLENGING_PERIOD_BLKS = await leagues.getChallengePeriod().should.be.fulfilled;
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

    const advanceNBlocks = async (blocks) => {
        let current = await web3.eth.getBlockNumber().should.be.fulfilled;
        await advanceToBlock(current + blocks);
    }

    const prepareMatchdayHashes = async (statesAtMatchday) => {
        let result = [];
        for (let i = 0; i < statesAtMatchday.length; i++) {
            const state = statesAtMatchday[i];
            const hash = await leagues.hashDayState(state).should.be.fulfilled;
            result.push(hash);
        }

        return result;
    }

    it('test2', async () => {
        await horizon.createTeam("Barca").should.be.fulfilled;
        await horizon.createTeam("Madrid").should.be.fulfilled;
        const teamIdx1 = await teams.getTeamId("Barca").should.be.fulfilled;
        const teamIdx2 = await teams.getTeamId("Madrid").should.be.fulfilled;

        await advanceToBlock(100);

        const blockInit = 190;
        const blockStep = 10;

        const usersInitData = {
            teamIdxs: [teamIdx1, teamIdx2],
            // teamOrders: [DEFAULT_ORDER, REVERSE_ORDER],
            tactics: [[4, 4, 2], [4, 3, 3]]
        };

        leagueIdx = 0;
        await leagues.create(leagueIdx, blockInit, blockStep, usersInitData.teamIdxs, usersInitData.tactics).should.be.fulfilled;

        const startBlock = await leagues.getInitBlock(leagueIdx).should.be.fulfilled;
        startBlock.toNumber().should.be.equal(blockInit);

        // Advance to matchday 2
        await advanceToBlock(blockInit + blockStep - 5);
        const started = await leagues.hasStarted(leagueIdx).should.be.fulfilled;
        started.should.be.equal(true);
        let finished = await leagues.hasFinished(leagueIdx).should.be.fulfilled;
        finished.should.be.equal(false);

        // Note that we could specify only for 1 of the teams if we wanted.
        const currentBlock = await web3.eth.getBlockNumber();
        usersAlongData = {
            teamIdxsWithinLeague: [teamIdx1, teamIdx2],
            tactics: [[4, 3, 3], [4, 4, 2]],
            blocks: [currentBlock, currentBlock]
        };

        // Submit data to change tactics
        await leagues.updateUsersAlongDataHash(leagueIdx, usersAlongData.teamIdxsWithinLeague, usersAlongData.tactics, usersAlongData.blocks).should.be.fulfilled;

        // Move beyond league end
        await advanceNBlocks(blockStep).should.be.fulfilled;
        finished = await leagues.hasFinished(leagueIdx).should.be.fulfilled;
        finished.should.be.equal(true);

        // START: The CLIENT computes the data needed to submit as an UPDATER: statesAtMatchday, scores.
        const leagueDays = await leagues.countLeagueDays(leagueIdx).should.be.fulfilled;
        leagueDays.toNumber().should.be.equal(2);

        // generate the init state of the teams from ERC721 because they never evolved
        const team1State = await generateTeamState(teamIdx1).should.be.fulfilled;
        const team2State = await generateTeamState(teamIdx2).should.be.fulfilled;

        // construct the initPlayerState data structure
        let initPlayerStates = await state.leagueStateCreate().should.be.fulfilled;
        initPlayerStates = await state.leagueStateAppend(initPlayerStates, team1State).should.be.fulfilled;
        initPlayerStates = await state.leagueStateAppend(initPlayerStates, team2State).should.be.fulfilled;

        // day 0
        let leagueDay = 0;
        const tacticsDay0 = [...usersInitData.tactics];
        const initPlayerStatesDay0 = [...initPlayerStates];
        let result = await leagues.computeDay(leagueIdx, leagueDay, initPlayerStatesDay0, tacticsDay0).should.be.fulfilled;
        const scoresDay0 = [...result.scores];
        const finalPlayerStatesDay0 = [...result.finalLeagueState];
        // day 1
        leagueDay = 1;
        const tacticsDay1 = [...usersAlongData.tactics];
        const initPlayerStatesDay1 = [...finalPlayerStatesDay0];
        result = await leagues.computeDay(leagueIdx, leagueDay, initPlayerStatesDay1, tacticsDay1).should.be.fulfilled;
        const scoresDay1 = [...result.scores];
        const finalPlayerStatesDay1 = [...result.finalLeagueState];

        const statesAtMatchday = [finalPlayerStatesDay0, finalPlayerStatesDay1];
        let scores = await leagues.scoresCreate().should.be.fulfilled;
        scores = await leagues.scoresConcat(scores, scoresDay0).should.be.fulfilled;
        scores = await leagues.scoresConcat(scores, scoresDay1).should.be.fulfilled;
        // END: The CLIENT computes the data needed to submit as an UPDATER: statesAtMatchday, scores.

        let updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(false);

        const initStatesHash = await leagues.hashInitState(initPlayerStates).should.be.fulfilled;
        const statesAtMatchdayHashes = await prepareMatchdayHashes(statesAtMatchday);

        let statesAtMatchdayLie = statesAtMatchday;
        statesAtMatchdayLie[0][0] += 1; // sinner operation!
        const statesAtMatchdayHashesLie = await prepareMatchdayHashes(statesAtMatchdayLie);

        await leagues.updateLeague(
            leagueIdx,
            initStatesHash,
            statesAtMatchdayHashesLie,
            scores
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(true);

        // A CHALLENGER tries to prove that the UPDATER lied with statesAtMatchday for matchday 0
        // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 20).should.be.fulfilled;
        let verified = await leagues.isVerified(leagueIdx).should.be.fulfilled;
        verified.should.be.equal(false);

        await leagues.challengeMatchdayStates(
            leagueIdx,
            usersInitData.teamIdxs,
            usersInitData.tactics,
            usersAlongData.teamIdxsWithinLeague,
            usersAlongData.tactics,
            usersAlongData.blocks,
            selectedMatchday = 0,
            prevMatchdayStates = initPlayerStatesDay0
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(false);

        // ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the initStates!:
        // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
        verified = await leagues.isVerified(leagueIdx).should.be.fulfilled;
        verified.should.be.equal(false);
        let initPlayerStatesLie = [...initPlayerStates];
        initPlayerStatesLie[0] += 1; // the sinner instruction
        const initStatesHashLie = await leagues.hashInitState(initPlayerStatesLie);
        await leagues.updateLeague(
            leagueIdx,
            initStatesHashLie,
            statesAtMatchdayHashes,
            scores
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(true);

        // A CHALLENGER tries to prove that the UPDATER lied with the initHash
        // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
        verified = await leagues.isVerified(leagueIdx).should.be.fulfilled;
        verified.should.be.equal(false);

        // create all the state of the environment (player team, owner, previous team ... league state)
        const dataToChallengeInitStates = [] // TODO;

        await leagues.challengeInitStates(
            leagueIdx,
            usersInitData.teamIdxs,
            usersInitData.tactics,
            dataToChallengeInitStates
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(false);

        // A nicer UPDATER now tells the truth:
        // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
        await leagues.updateLeague(
            leagueIdx,
            initStatesHash,
            statesAtMatchdayHashes,
            scores
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(true);

        // ...and the CHALLENGER fails to prove anything
        // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
        await leagues.challengeMatchdayStates(
            leagueIdx,
            usersInitData.teamIdxs,
            usersInitData.tactics,
            usersAlongData.teamIdxsWithinLeague,
            usersAlongData.tactics,
            usersAlongData.blocks,
            selectedMatchday = 0,
            prevMatchdayStates = initPlayerStatesDay0
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(true);

        await leagues.challengeMatchdayStates(
            leagueIdx,
            usersInitData.teamIdxs,
            usersInitData.tactics,
            usersAlongData.teamIdxsWithinLeague,
            usersAlongData.tactics,
            usersAlongData.blocks,
            selectedMatchday = 1,
            prevMatchdayStates = initPlayerStatesDay1
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(true);

        return;
        await leagues.challengeInitStates(
            leagueIdx,
            usersInitData.teamIdxs,
            usersInitData.tactics,
            dataToChallengeInitStates
        ).should.be.fulfilled;
        updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
        updated.should.be.equal(true);

        // We do not wait enough and try to:
        //   create another league. It fails to do so because teams are still busy
        await advanceNBlocks(2).should.be.fulfilled;
    });
})
