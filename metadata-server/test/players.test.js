// Import the dependencies for testing
const chai = require('chai');
const Web3 = require('web3');
const EthCrypto = require('eth-crypto');
const contractJSON = require('../../truffle-core/build/contracts/CryptoPlayers.json')
const ganache = require("ganache-cli");
const app = require('../app');

// Configure chai
chai.use(require('chai-http'));
chai.use(require('chai-as-promised'));
chai.should();

describe('player', () => {
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
    let instance = null;

    beforeEach(async () => {
        instance = await contract.deploy({
            data: contractJSON.bytecode
        })
            .send(sendOptions)
            .on('error', error => console.log("(EE) " + error))
            // .on('transactionHash', transactionHash => console.log("(II) transactionHash: " + transactionHash))
            // .on('receipt', receipt => console.log("(II) address: ", receipt.contractAddress)) // contains the new contract address
            // .on('confirmation', (confirmationNumber, receipt) => console.log("(II) confirmation: " + confirmationNumber))
            .catch(console.error)
    });

    it('check ERC721 metadata', async () => {
        await instance.methods.mint(identity.address, "player").send(sendOptions).should.be.fulfilled;
        const id = await instance.methods.getPlayerId("player").call().should.be.fulfilled;
        chai.request(app)
            .get('/players/' + id)
            .end((err, res) => {
                res.should.have.status(200);
                res.body.should.be.a('object');
                const json = JSON.parse(res.text);
                json.name.should.be.equal("Dave Starbelly");
                // json.description.should.be.equal("Friendly OpenSea Creature that enjoys long swims in the ocean.");
                // json.image.should.be.equal("https://storage.googleapis.com/opensea-prod.appspot.com/creature/3.png");
            });
    });
});