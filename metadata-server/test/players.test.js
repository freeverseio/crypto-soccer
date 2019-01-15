// Import the dependencies for testing
const chai = require('chai');
const Web3 = require('web3');
const EthCrypto = require('eth-crypto');
const contractJSON = require('../../truffle-core/build/contracts/CryptoPlayers.json')
const ganache = require("ganache-cli");
const playersJSON = require('../routes/playersJSON');
const config = require('../config.json');

// Configure chai
chai.use(require('chai-as-promised'));
chai.should();

describe('player', () => {
    let instance = null;
    let id = null;

    before(async () => {
        const identity = EthCrypto.createIdentity();
        const provider = ganache.provider({
            accounts: [{
                secretKey: identity.privateKey,
                balance: Web3.utils.toWei('100', 'ether')
            }]
        })
        const web3 = new Web3(provider);
        const contract = new web3.eth.Contract(contractJSON.abi);
        const sendOptions = {
            from: identity.address,
            gas: 4712388,
            gasPrice: provider.gasPrice
        };

        instance = await contract.deploy({
            data: contractJSON.bytecode
        })
            .send(sendOptions)
            .on('error', console.error)
            .catch(console.error);
        await instance.methods.mint(identity.address, "player").send(sendOptions).should.be.fulfilled;
        id = await instance.methods.getPlayerId("player").call().should.be.fulfilled;
    });

    it('check ERC721 metadata', async () => {
        const schema = await playersJSON(instance, id).should.be.fulfilled;
        schema.name.should.be.equal("player");
        schema.description.should.be.equal("put a description");
        schema.image.should.be.equal(config.players_image_base_URL + id);
    });

    it('check OpenSea metadata', async () => {
        const schema = await playersJSON(instance, id).should.be.fulfilled;
        schema.attributes.length.should.be.equal(5);
        const speed = await instance.methods.getSpeed(id).call().should.be.fulfilled;
        schema.attributes[0].trait_type.should.be.equal('speed');
        schema.attributes[0].value.should.be.equal(speed);
        const defence = await instance.methods.getDefence(id).call().should.be.fulfilled;
        schema.attributes[1].trait_type.should.be.equal('defence');
        schema.attributes[1].value.should.be.equal(defence);
        const endurance = await instance.methods.getEndurance(id).call().should.be.fulfilled;
        schema.attributes[2].trait_type.should.be.equal('endurance');
        schema.attributes[2].value.should.be.equal(endurance);
        const shoot = await instance.methods.getShoot(id).call().should.be.fulfilled;
        schema.attributes[3].trait_type.should.be.equal('shoot');
        schema.attributes[3].value.should.be.equal(shoot);
        const pass = await instance.methods.getPass(id).call().should.be.fulfilled;
        schema.attributes[4].trait_type.should.be.equal('pass');
        schema.attributes[4].value.should.be.equal(pass);
    });
});