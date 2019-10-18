const ganache = require('ganache-cli');
const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const genesis = require('./Genesis');
const Resolvers = require('../src/resolvers');

const identity = {
    address: '0x3Abf1775944E2B2C15c05D044632831f0Dfc9130',
    privateKey: '0x0a69684608770d018143dd70dc5dc5b6beadc366b87e45fcb567fc09407e7fe5'
};

// we preset the balance of our identities to 100 ether
const provider = ganache.provider({
    accounts: [{ secretKey: identity.privateKey, balance: '100000000000000000000000' }]
});

describe('assets resolvers', () => {
    let resolvers = null;

    beforeEach(async () => {
        contracts = await genesis(provider, identity.address);

        const { market } = contracts;
        resolvers = new Resolvers({
            market,
            from: identity.address
        });
    });

    describe('Mutation', () => {
        it('transferFirstBotToAddr', async () => {
            await resolvers.Mutation.transferFirstBotToAddr(_, { timezone: 0, countryIdxInTimezone: 0, address: "0x8c221609CC1C92FF648F3187fb12F8f821b46d9C" }).should.be.rejected;
        })
    });
});
