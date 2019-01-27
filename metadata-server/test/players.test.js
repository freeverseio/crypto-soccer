// Import the dependencies for testing
const chai = require('chai');
const Web3 = require('web3');
const EthCrypto = require('eth-crypto');
const ganache = require("ganache-cli");
const playersJSON = require('../routes/playersJSON');
const cryptoPlayersDeployer = require('./cryptoPlayersDeployer');
const cryptoTeamsDeployer = require('./cryptoTeamsDeployer');

// Configure chai
chai.use(require('chai-as-promised'));
chai.should();

describe('player', () => {
    let playersContract = null;
    let teamsContract = null;
    let playerId = null;

    before(async () => {
        const identity = EthCrypto.createIdentity();
        const provider = ganache.provider({
            accounts: [{
                secretKey: identity.privateKey,
                balance: Web3.utils.toWei('100', 'ether')
            }]
        });

        playersContract = await cryptoPlayersDeployer({ provider, sender: identity.address });
        teamsContract = await cryptoTeamsDeployer({ provider, playersContract, sender: identity.address });

        await playersContract.methods.mint(identity.address, "player").send({
            from: identity.address,
            gas: 4712388,
            gasPrice: provider.gasPrice
        }).should.be.fulfilled;
        playerId = await playersContract.methods.getPlayerId("player").call().should.be.fulfilled;

        await teamsContract.methods.mint(identity.address, "team").send({
            from: identity.address,
            gas: 4712388,
            gasPrice: provider.gasPrice
        }).should.be.fulfilled;
        const teamId = await teamsContract.methods.getTeamId("team").call().should.be.fulfilled;

        await teamsContract.methods.addPlayer(teamId, playerId).send({
            from: identity.address,
            gas: 4712388,
            gasPrice: provider.gasPrice
        }).should.be.fulfilled;
    });

    it('check ERC721 metadata', async () => {
        const schema = await playersJSON({playersContract, teamsContract, playerId}).should.be.fulfilled;
        schema.name.should.be.equal("player");
        schema.description.should.not.be.undefined;
        // schema.image.should.be.equal(config.players_image_base_URL + id);
        schema.image.should.be.equal('https://srv.latostadora.com/designall.dll/guybrush_threepwood--i:1413852880551413850;w:520;m:1;b:FFFFFF.jpg');
    });

    it('check OpenSea metadata', async () => {
        const schema = await playersJSON({playersContract, teamsContract, playerId}).should.be.fulfilled;
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
    });
});