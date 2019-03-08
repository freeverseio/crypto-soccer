const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

const League = artifacts.require('LeagueTactics');

contract('LeaguesTactics', (accounts) => {
    beforeEach(async () => {
        league = await League.new().should.be.fulfilled;
    });
}) 