// Import the dependencies for testing
const chai = require('chai');
const Web3 = require('web3');
const EthCrypto = require('eth-crypto');
const contractJSON = require('../../truffle-core/build/contracts/CryptoPlayers.json')
const ganache = require("ganache-cli");
const playersJSON = require('../routes/playersJSON');

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
        const schema = await playersJSON(instance, id).should.be.fulfilled;
        schema.name.should.be.equal("player");
    });
});