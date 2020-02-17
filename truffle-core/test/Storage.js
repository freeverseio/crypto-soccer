const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');

const Storage = artifacts.require('Storage');
const Assets = artifacts.require('Assets');

contract('Storage', (accounts) => {
    const FREEVERSE = accounts[0];
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    
    
    const it2 = async(text, f) => {};

    beforeEach(async () => {
        sto = await Storage.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
    });

    it('deploy correctly', async () => {
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        result = await sto.countFunctions().should.be.fulfilled;
        result.toNumber().should.be.equal(0);
    });

    it('add contract info and one function', async () => {
        await sto.addNewContract(addr = assets.address, name = "Assets").should.be.fulfilled;
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        for (c = 0; c < result.toNumber(); c++) { 
            let {0: ad, 1: nom} = await sto.getContractInfo(c).should.be.fulfilled;
            ad.should.be.equal(addr);
            nom.should.be.equal(name);
        }
        selector = web3.eth.abi.encodeFunctionSignature('setNewAsset(uint256,string)')
        await sto.addNewFunction(selector, contractId = 0).should.be.fulfilled;
        result = await sto.countFunctions().should.be.fulfilled;
        result.toNumber().should.be.equal(1);
    });


});