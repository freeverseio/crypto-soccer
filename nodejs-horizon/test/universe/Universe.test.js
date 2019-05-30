const ganache = require('ganache-cli');
require('chai')
    .use(require('chai-as-promised'))
    .should();
const Universe = require('../../src/universe/Universe');

const identity = {
    address: '0x3Abf1775944E2B2C15c05D044632831f0Dfc9130',
    privateKey: '0x0a69684608770d018143dd70dc5dc5b6beadc366b87e45fcb567fc09407e7fe5'
};

// we preset the balance of our identities to 100 ether
const provider = ganache.provider({
    accounts: [{ secretKey: identity.privateKey, balance: '100000000000000000000000' }]
});

describe('Universe', () => {
    let universe = null;

    beforeEach(async () => {
        universe = new Universe(provider, null, identity.address);
        await universe.genesis();
    });

    it('count teams', async () => {
        const count = await universe.countTeams().should.be.fulfilled;
        count.should.be.equal('0');
    });
});       