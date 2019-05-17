// Import the dependencies for testing
const chai = require('chai');
const EthCrypto = require('eth-crypto');
const { deployer, mintPlayer } = require('./environmentDeployer');
const { generateJSON } = require('../routes/players');

// Configure chai
chai.use(require('chai-as-promised'));
chai.should();

describe('player', () => {
    const identity = EthCrypto.createIdentity();
    let playersContract = null;
    let teamsContract = null;
    let playerId = null;

    before(async () => {
        const environment = await deployer(identity).should.be.fulfilled;
        playersContract = environment.playersContract;
        teamsContract = environment.teamsContract;
        playerId = environment.playerId;
    });

    it('check ERC721 metadata', async () => {
        const schema = await generateJSON({playersContract, teamsContract, playerId}).should.be.fulfilled;
        schema.name.should.be.equal("player");
        schema.description.should.not.be.undefined;
        // schema.image.should.be.equal(config.players_image_base_URL + id);
        schema.image.should.be.equal('http://metadata.busyverse.com:3000/images/' + playerId + '.svg');
    });

    it('check OpenSea metadata', async () => {
        const schema = await generateJSON({playersContract, teamsContract, playerId}).should.be.fulfilled;
        schema.external_url.should.be.equal("https://www.freeverse.io/");
        schema.attributes.length.should.be.equal(9);
        const speed = await playersContract.methods.getSpeed(playerId).call().should.be.fulfilled;
        schema.attributes[0].trait_type.should.be.equal('speed');
        schema.attributes[0].value.should.be.equal(Number(speed));
        const defence = await playersContract.methods.getDefence(playerId).call().should.be.fulfilled;
        schema.attributes[1].trait_type.should.be.equal('defence');
        schema.attributes[1].value.should.be.equal(Number(defence));
        const endurance = await playersContract.methods.getEndurance(playerId).call().should.be.fulfilled;
        schema.attributes[2].trait_type.should.be.equal('endurance');
        schema.attributes[2].value.should.be.equal(Number(endurance));
        const shoot = await playersContract.methods.getShoot(playerId).call().should.be.fulfilled;
        schema.attributes[3].trait_type.should.be.equal('shoot');
        schema.attributes[3].value.should.be.equal(Number(shoot));
        const pass = await playersContract.methods.getPass(playerId).call().should.be.fulfilled;
        schema.attributes[4].trait_type.should.be.equal('pass');
        schema.attributes[4].value.should.be.equal(Number(pass));
        schema.attributes[5].trait_type.should.be.equal('team');
        schema.attributes[5].value.should.be.equal('team');
    });

    it('team metadata when player has no team', async () => {
        const playerId = await mintPlayer(identity.address, "no team player").should.be.fulfilled;
        const schema = await generateJSON({ playersContract, teamsContract, playerId }).should.be.fulfilled;
        schema.should.not.be.undefined;
        schema.name.should.be.equal("no team player");
        schema.attributes[5].trait_type.should.be.equal('team');
        schema.attributes[5].value.should.be.equal('');
   });
});