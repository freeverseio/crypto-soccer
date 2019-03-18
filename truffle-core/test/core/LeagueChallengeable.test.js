const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueChallengeable');
const Cronos = artifacts.require('Cronos');

contract('LeagueChallengeable', (accounts) => {
    let league = null;
    let crosos = null;
    const id = 0;
    const teamIds = [1, 2];
    const tactics = [[4, 4, 3], [4, 5, 2]];
    let challengePeriod = null;

    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
        await league.create(
            id,
            blocksToInit = 1,
            step = 1,
            teamIds,
            tactics
        ).should.be.fulfilled;
        const result = await league.getChallengePeriod().should.be.fulfilled;
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
        await league.getLastChallengeBlock(id).should.be.rejected;
        const result = await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        const lastChallengeBlock = await league.getLastChallengeBlock(id).should.be.fulfilled;
        lastChallengeBlock.toNumber().should.be.equal(result.receipt.blockNumber + challengePeriod);
    });

    it('reset league makes invalid last challenge block', async () => {
        await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        await league.resetUpdater(id).should.be.fulfilled;
        await league.getLastChallengeBlock(id).should.be.rejected;
    });

    it('is verified', async () => {
        let verified = await league.isVerified(id).should.be.fulfilled;
        verified.should.be.equal(false);
        let result = await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        verified = await league.isVerified(id).should.be.fulfilled;
        verified.should.be.equal(false);
        const updateBlockNumber = result.receipt.blockNumber;
        await advanceToBlock(updateBlockNumber + challengePeriod - 1).should.be.fulfilled;
        verified = await league.isVerified(id).should.be.fulfilled;
        verified.should.be.equal(false);
        await advanceToBlock(updateBlockNumber + challengePeriod).should.be.fulfilled;
        verified = await league.isVerified(id).should.be.fulfilled;
        verified.should.be.equal(true);
    });

    it('challenge init state', async () => {
        await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        await league.challengeInitStates(id, teamIds, tactics, dataToChallengeInitStates = []).should.be.fulfilled;
    });

    it('challenge init state with wrong user init data', async () => {
        await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        await league.challengeInitStates(id, [3, 4], tactics, []).should.be.rejected;
        await league.challengeInitStates(id, teamIds, [[4, 4, 2], [4, 4, 2]], []).should.be.rejected;
    });
})