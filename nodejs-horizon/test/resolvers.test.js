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
        universe = new Universe({
            provider,
            from: identity.address
        });
        universe.web3.currentProvider.setMaxListeners(0);
        await universe.genesis();
        resolvers = new Resolvers(universe);
    });

    it('countTeams', async () => {
        let count = await resolvers.Query.countTeams().should.be.fulfilled;
        count.should.be.equal('0');
    });

    it('create team', async () => {
        await resolvers.Mutation.createTeam(_, { name: "Barca", owner: identity.address }).should.be.fulfilled;
        let count = await resolvers.Query.countTeams().should.be.fulfilled;
        count.should.be.equal('1');
    });

    it('get player', async () => {
        const id = 3;
        const player = await resolvers.Query.getPlayer(_, { id });
        player.should.be.equal(id);
    });

    it('get all teams', async () => {
        const teams = await resolvers.Query.allTeams().should.be.fulfilled;
    });

    describe('Player', () => {
        it('id', () => {
            resolvers.Player.id(3).should.be.equal(3);
        });

        it('name', () => {
            resolvers.Player.name(3).should.be.equal('player_3');
        });

        it('defence', async () => {
            const id = 3;
            await resolvers.Player.defence(id).should.be.rejected;
            await resolvers.Mutation.createTeam(_, { name: "Barca", owner: identity.address }).should.be.fulfilled;
            const skill = await resolvers.Player.defence(id).should.be.fulfilled;
            skill.should.be.equal('50');
        });

        it('speed', async () => {
            const id = 3;
            await resolvers.Player.speed(id).should.be.rejected;
            await resolvers.Mutation.createTeam(_, { name: "Barca", owner: identity.address }).should.be.fulfilled;
            const skill = await resolvers.Player.speed(id).should.be.fulfilled;
            skill.should.be.equal('62');
        }); 
        
        it('pass', async () => {
            const id = 3;
            await resolvers.Player.pass(id).should.be.rejected;
            await resolvers.Mutation.createTeam(_, { name: "Barca", owner: identity.address }).should.be.fulfilled;
            const skill = await resolvers.Player.pass(id).should.be.fulfilled;
            skill.should.be.equal('47');
        }); 
        
        it('shoot', async () => {
            const id = 3;
            await resolvers.Player.shoot(id).should.be.rejected;
            await resolvers.Mutation.createTeam(_, { name: "Barca", owner: identity.address }).should.be.fulfilled;
            const skill = await resolvers.Player.shoot(id).should.be.fulfilled;
            skill.should.be.equal('27');
        }); 
        
        it('endurance', async () => {
            const id = 3;
            await resolvers.Player.endurance(id).should.be.rejected;
            await resolvers.Mutation.createTeam(_, { name: "Barca", owner: identity.address }).should.be.fulfilled;
            const skill = await resolvers.Player.endurance(id).should.be.fulfilled;
            skill.should.be.equal('64');
        }); 
    });
});
