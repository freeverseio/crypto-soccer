// Import the dependencies for testing
const chai = require('chai');
const EthCrypto = require('eth-crypto');
const { deployer, mintPlayer, addPlayer } = require('./environmentDeployer');
const teamsJSON = require('../routes/teamsJSON');

// Configure chai
chai.use(require('chai-as-promised'));
chai.should();

describe('teams', () => {
    const identity = EthCrypto.createIdentity();
    let playersContract = null;
    let teamsContract = null;
    let teamId = null;

    before(async () => {
        const environment = await deployer(identity).should.be.fulfilled;
        playersContract = environment.playersContract;
        teamsContract = environment.teamsContract;
        teamId = environment.teamId;
    });

    it('check ERC721 metadata', async () => {
        const schema = await teamsJSON({playersContract, teamsContract, teamId}).should.be.fulfilled;
        schema.name.should.be.equal("team");
        schema.description.should.not.be.undefined;
        // schema.image.should.be.equal(config.players_image_base_URL + id);
        schema.image.should.be.equal('http://www.monkers.net/fotos/lucasartsDead.jpg');
    });

    it('check the players of the team', async () => {
        for (let i = 0; i < 5; i++) {
            const playerId = await mintPlayer(identity.address, "player" + i).should.be.fulfilled;
            await addPlayer(identity.address, teamId, playerId).should.be.fulfilled;
        }
        const schema = await teamsJSON({ playersContract, teamsContract, teamId }).should.be.fulfilled;
        schema.attributes.length.should.be.equal(6);
    });
});