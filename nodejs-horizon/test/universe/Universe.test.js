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
        universe = new Universe({
            provider,
            from: identity.address
        });
        universe.web3.currentProvider.setMaxListeners(0);
        await universe.genesis();
    });

    it('count teams', async () => {
        const count = await universe.countTeams().should.be.fulfilled;
        count.should.be.equal('0');
    });

    it('create team', async () => {
        await universe.createTeam("Barca", identity.address).should.be.fulfilled;
        const count = await universe.countTeams().should.be.fulfilled;
        count.should.be.equal('1');
    });

    it('get team name', async () => {
        await universe.getTeamName(1).should.be.rejected;
        await universe.createTeam("Barca", identity.address).should.be.fulfilled;
        const name = await universe.getTeamName(1).should.be.fulfilled;
        name.should.be.equal("Barca");
    });

    it('get team player ids', async () => {
        await universe.getTeamPlayerIds(1).should.be.rejected;
        await universe.createTeam("Barca", identity.address).should.be.fulfilled;
        const players = await universe.getTeamPlayerIds(1).should.be.fulfilled;
        players.should.be.eql([ '1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11' ]);
    })
});       