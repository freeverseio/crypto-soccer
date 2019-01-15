// Import the dependencies for testing
const chai = require('chai');
const chaiHttp = require('chai-http');
const app = require('../app');
const Web3 = require('web3');
const EthCrypto = require('eth-crypto');
const contractJSON = require('../../truffle-core/build/contracts/CryptoPlayers.json')
const ganache = require("ganache-cli");

// Configure chai
chai.use(chaiHttp);
chai.should();

describe('routing', () => {
    const identity = EthCrypto.createIdentity();
    console.log(identity)
    const provider = ganache.provider({
        accounts: [{secretKey: identity.privateKey, balance: Web3.utils.toWei('100', 'ether') }]
    })
    const web3 = new Web3(provider);
    const contract = new web3.eth.Contract(contractJSON.abi);
    let instance = null;

    beforeEach(async () => {
        instance = await contract.deploy({
            data: contractJSON.bytecode
        })
            .send({
                from: identity.address,
                gas: 4712388,
                gasPrice: 100000000000
            })
            .on('error', error => console.log("(EE) " + error))
            .on('transactionHash', transactionHash => console.log("(II) transactionHash: " + transactionHash))
            .on('receipt', receipt => console.log("(II) address: ", receipt.contractAddress)) // contains the new contract address
            .on('confirmation', (confirmationNumber, receipt) => console.log("(II) confirmation: " + confirmationNumber))
            .then(blackHole => {
                console.log("(II) ... done");
                return blackHole;
            })
            .catch(reason => {
                console.error(reason);
                process.exit(1);
            });
    });

    it('check deployment', () => {
        console.log(instance.address);
    });
});