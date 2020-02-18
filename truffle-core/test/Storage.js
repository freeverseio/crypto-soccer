const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');

const StorageProxy = artifacts.require('StorageProxy');
const Assets = artifacts.require('Assets');

contract('StorageProxy', (accounts) => {
    const FREEVERSE = accounts[0];
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const it2 = async(text, f) => {};

    const assetsId = 0;
    const marketId = 1;
    const updatesId = 2;
    
    // const assetsFunctions = [
    //     ["init()", assetsId],
    //     ["initSingleTZ(uint8)", assetsId],
    //     ["init()", assetsId],
    //     ["init()", assetsId],
    //     ["init()", assetsId],
    //     ["init()", assetsId],
    //     ["init()", assetsId],
    //     ["init()", assetsId],
    // ]
    
    beforeEach(async () => {
        sto = await StorageProxy.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
    });
    
    it('call a function inside Assets via delegate call', async () => {
        await sto.addNewContract(addr = assets.address, name = "Assets").should.be.fulfilled;
        selector = web3.eth.abi.encodeFunctionSignature('init()')
        await sto.addNewFunction(selector, contractId = 0).should.be.fulfilled;
        isInit = await assets.getIsInit().should.be.fulfilled;
        isInit.should.be.equal(false);
        await sto.sendTransaction({data: selector}).should.be.fulfilled;
        isInit = await assets.getIsInit().should.be.fulfilled;
        isInit.should.be.equal(false);
        // show that you cannot init twice:
        await sto.sendTransaction({data: selector}).should.be.rejected;
    });
return
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

    it('autorizations check', async () => {
        await sto.setStorageOwner(addr = ALICE, {from: ALICE}).should.be.rejected;
        await sto.setStorageOwner(addr = ALICE, {from: FREEVERSE}).should.be.fulfilled;
        await sto.addNewContract(addr = assets.address, name = "Assets", {from: FREEVERSE}).should.be.rejected;
        await sto.addNewContract(addr = assets.address, name = "Assets", {from: ALICE}).should.be.fulfilled;
        selector = web3.eth.abi.encodeFunctionSignature('setNewAsset(uint256,string)')
        await sto.addNewFunction(selector, contractId = 0, {from: FREEVERSE}).should.be.rejected;
        await sto.addNewFunction(selector, contractId = 0, {from: ALICE}).should.be.fulfilled;
    });


});