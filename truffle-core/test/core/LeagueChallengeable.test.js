const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueChallengeable');

contract('LeaguesScheduler', (accounts) => {
    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
    });
 
    it('challenge period', async () => {
        const period = await league.getChallengePeriod().should.be.fulfilled;
        period.should.be.a.bignumber.that.equals('60');
    });
})