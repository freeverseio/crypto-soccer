const ganache = require('ganache-cli');
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const Universe = require('../src/universe/Universe');
const Resolvers = require('../src/resolvers');

const identity = {
    address: '0x3Abf1775944E2B2C15c05D044632831f0Dfc9130',
    privateKey: '0x0a69684608770d018143dd70dc5dc5b6beadc366b87e45fcb567fc09407e7fe5'
};

// we preset the balance of our identities to 100 ether
const provider = ganache.provider({
    accounts: [{ secretKey: identity.privateKey, balance: '100000000000000000000000' }]
})


describe('assets resolvers', () => {
    let resolvers = null;

    beforeEach(async () => {
        universe = new Universe(provider, null, identity.address);
        universe.web3.currentProvider.setMaxListeners(0);
        await universe.genesis();
        resolvers = new Resolvers(universe);
    });

    it('countTeams', async () => {
        let count = await resolvers.Query.countTeams().should.be.fulfilled;
        count.should.be.equal('0');
    });
});
