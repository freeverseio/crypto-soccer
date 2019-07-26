const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const Engine = artifacts.require('Engine');
const States = artifacts.require('LeagueState');
const League = artifacts.require('LeagueChallengeable');
const Cronos = artifacts.require('Cronos');
const Assets = artifacts.require('Assets');
const PlayerStateLib = artifacts.require('PlayerState');

contract('LeagueChallengeable', (accounts) => {
    let leagues = null;
    let cronos = null;
    let assets = null;
    let playerStateLib = null;
    const initBlock = 1;
    const step = 1;
    const leagueId = 1;
    const PLAYERS_PER_TEAM = 25;
    const order = Array.from(new Array(PLAYERS_PER_TEAM), (x,i) => i) // [0,1,...24]
    const tactic442 = 0;
    const tactic541 = 1;
    const tacticsIds = [tactic442, tactic541];
    const teamIds = [1, 2];
    let challengePeriod = null;

    beforeEach(async () => {
        playerStateLib = await PlayerStateLib.new().should.be.fulfilled;
        assets = await Assets.new(playerStateLib.address).should.be.fulfilled;
        const engine = await Engine.new().should.be.fulfilled;
        const states = await States.new().should.be.fulfilled;
        leagues = await League.new(engine.address, states.address).should.be.fulfilled;
        await leagues.setAssetsContract(assets.address).should.be.fulfilled;
        await assets.createTeam(name = "Barca", accounts[1]).should.be.fulfilled;
        await assets.createTeam(name = "Mardid", accounts[2]).should.be.fulfilled;
        await leagues.create(nTeams = 2, initBlock, step).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 1, order, tactic442).should.be.fulfilled;
        await leagues.signTeamInLeague(leagueId, teamId = 2, order, tactic442).should.be.fulfilled;
        const result = await leagues.getChallengePeriod().should.be.fulfilled;
        challengePeriod = result.toNumber();
        cronos = await Cronos.new().should.be.fulfilled;
    });

    
    const advanceToBlock = async (block) => {
        let current = await web3.eth.getBlockNumber().should.be.fulfilled;
        while (current.toString() < block) {
            await cronos.wait().should.be.fulfilled;
            current = await web3.eth.getBlockNumber().should.be.fulfilled;
        }
        // console.log("current block: " + current);
    }
 

    it('challenge period', async () => {
        challengePeriod.should.be.equal(60);
    });

    it('last challenge block', async () => {
        await leagues.getLastChallengeBlock(leagueId).should.be.rejected;
        const result = await leagues.updateLeague(
            leagueId, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3'],
            isLie = false
        ).should.be.fulfilled;
        const lastChallengeBlock = await leagues.getLastChallengeBlock(leagueId).should.be.fulfilled;
        lastChallengeBlock.toNumber().should.be.equal(result.receipt.blockNumber + challengePeriod);
    });

    it('reset league makes invalid last challenge block', async () => {
        await leagues.updateLeague(
            leagueId, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3'],
            isLie = false
        ).should.be.fulfilled;
        await leagues.resetUpdater(leagueId).should.be.fulfilled;
        await leagues.getLastChallengeBlock(leagueId).should.be.rejected;
    });

    it('is verified', async () => {
        let verified = await leagues.isVerified(leagueId).should.be.fulfilled;
        verified.should.be.equal(false);
        let result = await leagues.updateLeague(
            leagueId, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3'],
            isLie = false
        ).should.be.fulfilled;
        verified = await leagues.isVerified(leagueId).should.be.fulfilled;
        verified.should.be.equal(false);
        const updateBlockNumber = result.receipt.blockNumber;
        await advanceToBlock(updateBlockNumber + challengePeriod - 1).should.be.fulfilled;
        verified = await leagues.isVerified(leagueId).should.be.fulfilled;
        verified.should.be.equal(false);
        await advanceToBlock(updateBlockNumber + challengePeriod).should.be.fulfilled;
        verified = await leagues.isVerified(leagueId).should.be.fulfilled;
        verified.should.be.equal(true);
    });

    it('challenge init state', async () => {
        await leagues.updateLeague(
            leagueId, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3'],
            isLie = true
        ).should.be.fulfilled;
        const receipt = await leagues.challengeInitStates(leagueId, teamIds, tacticsIds, dataToChallengeInitStates = []).should.be.fulfilled;
        receipt.logs[0].args.challengeSucceeded.should.be.equal(true);
    });


    it('challenge init state when league was correctly updated', async () => {
        await leagues.updateLeague(
            leagueId, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3'],
            isLie = false
        ).should.be.fulfilled;
        let receipt = await leagues.challengeInitStates(leagueId, [3, 4], tacticsIds, []).should.be.fulfilled;
        receipt.logs[0].args.challengeSucceeded.should.be.equal(false);
        receipt = await leagues.challengeInitStates(leagueId, teamIds, tacticsIds, []).should.be.fulfilled;
        receipt.logs[0].args.challengeSucceeded.should.be.equal(false);
    });


    // it('update tacticsIds with no new tacticsIds', async () => {
    //     const tacticsIds = await leagues.updateTacticsToBlockNum(
    //         usersInitDataTeamIds = [1],
    //         userInitDataTactics = [[4,4,2]],
    //         blockNum = [10],
    //         usersAlongDataTeamIds = [],
    //         usersAlongDataTactics = [],
    //         usersAlongDataBlocks = []
    //     ).should.be.fulfilled;
    //     tacticsIds.length.should.be.equal(3);
    //     tacticsIds[0].toNumber().should.be.equal(4);
    //     tacticsIds[1].toNumber().should.be.equal(4);
    //     tacticsIds[2].toNumber().should.be.equal(2);
    // });

    // it('update tacticsIds new tacticsIds', async () => {
    //     const tacticsIds = await leagues.updateTacticsToBlockNum(
    //         usersInitDataTeamIds = [1, 5],
    //         userInitDataTactics = [[4, 4, 2], [5, 5, 0]],
    //         blockNum = [10],
    //         usersAlongDataTeamIds = [1, 5, 2],
    //         usersAlongDataTactics = [[5, 3, 2], [4, 4, 2], [1, 8, 1]],
    //         usersAlongDataBlocks = [1, 8, 10]
    //     ).should.be.fulfilled;
    //     tacticsIds.length.should.be.equal(6);
    //     tacticsIds[0].toNumber().should.be.equal(5);
    //     tacticsIds[1].toNumber().should.be.equal(3);
    //     tacticsIds[2].toNumber().should.be.equal(2);
    //     tacticsIds[3].toNumber().should.be.equal(4);
    //     tacticsIds[4].toNumber().should.be.equal(4);
    //     tacticsIds[5].toNumber().should.be.equal(2);
    // });
})