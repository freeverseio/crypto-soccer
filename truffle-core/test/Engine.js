const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');

const Engine = artifacts.require('Engine');

contract('Engine', (accounts) => {
    let engine = null;

    beforeEach(async () => {
        engine = await Engine.new().should.be.fulfilled;
        encoding = engine;
        await engine.init().should.be.fulfilled;
        });

    it('check initial constants', async () =>  {
        // MATCHDAYS.toNumber().should.be.equal(14);
        // MATCHES_PER_DAY.toNumber().should.be.equal(4);
        // TEAMS_PER_LEAGUE.toNumber().should.be.equal(8);
    });


});