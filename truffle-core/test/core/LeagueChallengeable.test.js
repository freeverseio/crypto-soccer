const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueChallengeable');

contract('LeagueChallengeable', (accounts) => {
    const id = 0;
    const teamIds = [1, 2];
    const tactics = [[4, 4, 3], [4, 5, 2]];

    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
        await league.create(
            id,
            blocksToInit = 1,
            step = 1,
            teamIds,
            tactics
        ).should.be.fulfilled;
    });
 
    it('challenge period', async () => {
        const period = await league.getChallengePeriod().should.be.fulfilled;
        period.should.be.a.bignumber.that.equals('60');
    });

    it('challenge init state', async () => {
        await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        await league.challengeInitStates(id, teamIds, tactics).should.be.fulfilled;
    });

    it('challenge init state with wrong user init data', async () => {
        await league.updateLeague(
            id, 
            initStateHash = '0x54564', 
            dayStateHashes = ['0x24353', '0x5434432'],
            scores = ['0x12', '0x3']
        ).should.be.fulfilled;
        await league.challengeInitStates(id, [3, 4], tactics).should.be.rejected;
        await league.challengeInitStates(id, teamIds, [[4, 4, 2], [4, 4, 2]]).should.be.rejected;
    });
})