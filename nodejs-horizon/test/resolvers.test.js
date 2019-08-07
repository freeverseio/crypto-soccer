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

        const { states, assets, leagues } = contracts;
        resolvers = new Resolvers({
            states,
            assets,
            leagues,
            from: identity.address
        });
    });

    describe('Mutation', () => {
        it('createLeague', async () => {
            const result = await resolvers.Mutation.createLeague(_, { nTeams: 4, initBlock: 10, step: 20 }).should.be.fulfilled;
            result.should.be.equal(true);
        });

        it('transferTeam', async () => {
            const gas = await assets.methods.createTeam("Madrid", identity.address).estimateGas().should.be.fulfilled;
            await assets.methods.createTeam("Madrid", identity.address).send({ from: identity.address, gas }).should.be.fulfilled;
            const result = await resolvers.Mutation.transferTeam(_, { teamId: "1", owner: "0x8c221609CC1C92FF648F3187fb12F8f821b46d9C"}).should.be.fulfilled;
            result.should.be.equal(true);
        })
    });
});
