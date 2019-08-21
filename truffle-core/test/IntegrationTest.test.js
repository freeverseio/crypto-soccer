require('chai')
    .use(require('chai-as-promised'))
    .should();
const util = require('util');

const Assets = artifacts.require('Assets');
const Leagues = artifacts.require('Leagues');
const Engine = artifacts.require('Engine');
const State = artifacts.require('LeagueState');
const Cronos = artifacts.require('Cronos');
const GameController = artifacts.require('GameController');

const UNENROLLED       = 0;
const ENROLLING        = 1;
const UNENROLLING      = 2;
const UNENROLLABLE     = 3;
const ENROLLED         = 4;
const CHALLENGE_TT     = 5;
const CHALLENGE_LI_RES = 6;
const CHALLENGE_TT_RES = 7;
const SLASHABLE        = 8;
const SLASHED          = 9

const PLAYERS_PER_TEAM = 25;
const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
const tactic442 = 0;
const tactic511 = 1;
const tactic433 = 2;


contract('IntegrationTest', (accounts) => {
    const [owner, bob, alice, carol] = accounts

    let engine = null;
    let state = null;
    let leagues = null;
    let cronos = null;
    let stake = null;
    let currentBlock = null;

    const onion0 = web3.utils.keccak256("hel24"); // finishes 2
    const onion1 = web3.utils.keccak256(onion0);  // finishes 0
    const onion2 = web3.utils.keccak256(onion1);  // finishes b
    const onion3 = web3.utils.keccak256(onion2);  // finishes

    beforeEach(async () => {
        state = await State.new().should.be.fulfilled;
        assets = await Assets.new(state.address).should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        leagues = await Leagues.new(engine.address, state.address).should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        cronos = await Cronos.new().should.be.fulfilled;
        controller = await GameController.new(leagues.address).should.be.fulfilled;
        await controller.setGameContractAddress(leagues.address).should.be.fulfilled;
        stake = await controller.REQUIRED_STAKE();

        await leagues.setStakersContract(controller.address).should.be.fulfilled;

        await controller.enroll(onion2,{from:bob, value:stake});
        await controller.enroll(onion3,{from:alice, value:stake});
        await controller.enroll(onion3,{from:carol, value:stake});
        await jumpSeconds((await controller.MINENROLL_SECS()).toNumber());
    });

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
            const hash = await leagues.hashDayState(state).should.be.fulfilled;
            result.push(hash);
        }

        return result;
    }

    // it('play a league', async () => {
    //     let receipt = await assets.createTeam("Barca", accounts[0]).should.be.fulfilled;
    //     const teamIdx1 = receipt.logs[0].args.id.toNumber();
    //     teamIdx1.should.be.equal(1);

    //     receipt = await assets.createTeam("Madrid", accounts[0]).should.be.fulfilled;
    //     const teamIdx2 = receipt.logs[0].args.id.toNumber();
    //     teamIdx2.should.be.equal(2);

    //     currentBlock = await web3.eth.getBlockNumber();

    //     let blockInit = currentBlock + 90;
    //     let blockStep = 10;

    //     const usersInitData = {
    //         teamIdxs: [teamIdx1, teamIdx2],
    //         teamOrders: [order, order],
    //         tacticsIds: [tactic442, tactic433]
    //     };

    //     let leagueIdx = 1;
        
    //     await leagues.create(nTeams = 2, blockInit, blockStep).should.be.fulfilled;
    //     for (var team = 0; team < 2; team++) {
    //         await leagues.signTeamInLeague(
    //             leagueIdx, 
    //             usersInitData.teamIdxs[team], 
    //             usersInitData.teamOrders[team], 
    //             usersInitData.tacticsIds[team]
    //         ).should.be.fulfilled;
    //     }

    //     const startBlock = await leagues.getInitBlock(leagueIdx).should.be.fulfilled;
    //     startBlock.toNumber().should.be.equal(blockInit);

    //     // Advance to matchday 2
    //     await advanceToBlock(blockInit + blockStep/2);
    //     const started = await leagues.hasStarted(leagueIdx).should.be.fulfilled;
    //     started.should.be.equal(true);
    //     let finished = await leagues.hasFinished(leagueIdx).should.be.fulfilled;
    //     finished.should.be.equal(false);

    //     // Note that we could specify only for 1 of the teams if we wanted.
    //     currentBlock = await web3.eth.getBlockNumber();
    //     usersAlongData = {
    //         teamIdxsWithinLeague: [teamIdx1, teamIdx2],
    //         tacticsIds: [tactic433, tactic442],
    //         blocks: [currentBlock, currentBlock]
    //     };

    //     // Submit data to change tacticsIds
    //     await leagues.updateUsersAlongDataHash(leagueIdx, usersAlongData.teamIdxsWithinLeague, usersAlongData.tacticsIds, usersAlongData.blocks).should.be.fulfilled;

    //     // Move beyond league end
    //     await advanceNBlocks(blockStep).should.be.fulfilled;
    //     finished = await leagues.hasFinished(leagueIdx).should.be.fulfilled;
    //     finished.should.be.equal(true);

    //     // START: The CLIENT computes the data needed to submit as an UPDATER: statesAtMatchday, scores.
    //     const leagueDays = await leagues.countLeagueDays(leagueIdx).should.be.fulfilled;
    //     leagueDays.toNumber().should.be.equal(2);

    //     // generate the init state of the teams from ERC721 because they never evolved
    //     const team1State = await generateTeamState(teamIdx1).should.be.fulfilled;
    //     const team2State = await generateTeamState(teamIdx2).should.be.fulfilled;

    //     // construct the initPlayerState data structure
    //     let initPlayerStates = await state.leagueStateCreate().should.be.fulfilled;
    //     initPlayerStates = await state.leagueStateAppend(initPlayerStates, team1State).should.be.fulfilled;
    //     initPlayerStates = await state.leagueStateAppend(initPlayerStates, team2State).should.be.fulfilled;

    //     // day 0
    //     let leagueDay = 0;
    //     const tacticsDay0 = [...usersInitData.tacticsIds];
    //     const initPlayerStatesDay0 = [...initPlayerStates];
    //     let result = await leagues.computeDay(leagueIdx, leagueDay, initPlayerStatesDay0, tacticsDay0).should.be.fulfilled;
    //     const scoresDay0 = [...result.scores];
    //     const finalPlayerStatesDay0 = [...result.finalLeagueState];
    //     // day 1
    //     leagueDay = 1;
    //     const tacticsDay1 = [...usersAlongData.tacticsIds];
    //     const initPlayerStatesDay1 = [...finalPlayerStatesDay0];
    //     result = await leagues.computeDay(leagueIdx, leagueDay, initPlayerStatesDay1, tacticsDay1).should.be.fulfilled;
    //     const scoresDay1 = [...result.scores];
    //     const finalPlayerStatesDay1 = [...result.finalLeagueState];

    //     const statesAtMatchday = [finalPlayerStatesDay0, finalPlayerStatesDay1];
    //     let scores = await leagues.scoresCreate().should.be.fulfilled;
    //     scores = await leagues.scoresConcat(scores, scoresDay0).should.be.fulfilled;
    //     scores = await leagues.scoresConcat(scores, scoresDay1).should.be.fulfilled;
    //     // END: The CLIENT computes the data needed to submit as an UPDATER: statesAtMatchday, scores.

    //     let updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(false);

    //     const initStatesHash = await leagues.hashInitState(initPlayerStates).should.be.fulfilled;
    //     const statesAtMatchdayHashes = await prepareMatchdayHashes(statesAtMatchday);

    //     let statesAtMatchdayLie = statesAtMatchday;
    //     statesAtMatchdayLie[0][0] += 1; // sinner operation!
    //     const statesAtMatchdayHashesLie = await prepareMatchdayHashes(statesAtMatchdayLie);

    //     console.log("updating the league");
    //     await leagues.updateLeague(
    //         leagueIdx,
    //         initStatesHash,
    //         statesAtMatchdayHashesLie,
    //         scores,
    //         true,
    //         {from: bob}
    //     ).should.be.fulfilled;
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(true);

    //     console.log("update done");
    //     // A CHALLENGER tries to prove that the UPDATER lied with statesAtMatchday for matchday 0
    //     // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 20).should.be.fulfilled;
    //     let verified = await leagues.isVerified(leagueIdx).should.be.fulfilled;
    //     verified.should.be.equal(false);

    //     console.log("challenging");
    //     console.log("State Bob: " + await controller.state(bob,0));
    //     console.log("State Alice: " + await controller.state(alice,0));

    //     await leagues.challengeMatchdayStates(
    //         leagueIdx,
    //         usersInitData.teamIdxs,
    //         usersInitData.tacticsIds,
    //         usersAlongData.teamIdxsWithinLeague,
    //         usersAlongData.tacticsIds,
    //         usersAlongData.blocks,
    //         selectedMatchday = 0,
    //         prevMatchdayStates = initPlayerStatesDay0,
    //         {from: alice}

    //     ).should.be.fulfilled;
        
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(false);
    //     console.log("challenged");

    //     // ...and the CLIENT, acting as an UPDATER, submits to the BC... a lie in the initStates!:
    //     // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
    //     verified = await leagues.isVerified(leagueIdx).should.be.fulfilled;
    //     verified.should.be.equal(false);
    //     let initPlayerStatesLie = [...initPlayerStates];
    //     initPlayerStatesLie[0] += 1; // the sinner instruction
    //     const initStatesHashLie = await leagues.hashInitState(initPlayerStatesLie);
    //     await leagues.updateLeague(
    //         leagueIdx,
    //         initStatesHashLie,
    //         statesAtMatchdayHashes,
    //         scores,
    //         true,
    //         {from: alice}
    //     ).should.be.fulfilled;
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(true);
    //     console.log("alice updated");

    //     // A CHALLENGER tries to prove that the UPDATER lied with the initHash
    //     // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
    //     verified = await leagues.isVerified(leagueIdx).should.be.fulfilled;
    //     verified.should.be.equal(false);

    //     // create all the state of the environment (player team, owner, previous team ... league state)
    //     const dataToChallengeInitStates = [] // TODO;

    //     await leagues.challengeInitStates(
    //         leagueIdx,
    //         usersInitData.teamIdxs,
    //         usersInitData.tacticsIds,
    //         dataToChallengeInitStates, {from: bob}
    //     ).should.be.fulfilled;
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(false);

    //     console.log("bob challenged");

    //     await controller.resolveChallenge(onion1, {from: bob}).should.be.fulfilled;
    //     console.log("bob reveal the secret");

    //     // A nicer UPDATER now tells the truth:
    //     // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
    //     await leagues.updateLeague(
    //         leagueIdx,
    //         initStatesHash,
    //         statesAtMatchdayHashes,
    //         scores,
    //         false,
    //         {from: bob}
    //     ).should.be.fulfilled;
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(true);
    //     console.log("bob updated 2");


    //     // ...and the CHALLENGER fails to prove anything
    //     // await advanceNBlocks(CHALLENGING_PERIOD_BLKS - 5);
    //     await leagues.challengeMatchdayStates(
    //         leagueIdx,
    //         usersInitData.teamIdxs,
    //         usersInitData.tacticsIds,
    //         usersAlongData.teamIdxsWithinLeague,
    //         usersAlongData.tacticsIds,
    //         usersAlongData.blocks,
    //         selectedMatchday = 0,
    //         prevMatchdayStates = initPlayerStatesDay0, {from: carol}
    //     ).should.be.fulfilled;
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(true);

    //     await leagues.challengeMatchdayStates(
    //         leagueIdx,
    //         usersInitData.teamIdxs,
    //         usersInitData.tacticsIds,
    //         usersAlongData.teamIdxsWithinLeague,
    //         usersAlongData.tacticsIds,
    //         usersAlongData.blocks,
    //         selectedMatchday = 1,
    //         prevMatchdayStates = initPlayerStatesDay1, {from: carol}
    //     ).should.be.fulfilled;
    //     updated = await leagues.isUpdated(leagueIdx).should.be.fulfilled;
    //     updated.should.be.equal(true);

    //     await jumpSeconds((await controller.MAXIDLE_SECS()).toNumber());
    //     // console.log("alice state : " + await controller.state(alice,0));
    //     await controller.slash(alice, {from: bob}).should.be.fulfilled;
    //     console.log("Alice slashed");
    // });


const jumpSeconds = function(duration) {
  const id = Date.now()
  const sendAsync = util.promisify(web3.currentProvider.send);

  return new Promise((resolve, reject) => {
    sendAsync({
      jsonrpc: '2.0',
      method: 'evm_increaseTime',
      params: [duration],
      id: id,
    }, err1 => {
      if (err1) return reject(err1)

      sendAsync({
        jsonrpc: '2.0',
        method: 'evm_mine',
        id: id+1,
      }, (err2, res) => {
        return err2 ? reject(err2) : resolve(res)
      })
    })
  })
}
})
