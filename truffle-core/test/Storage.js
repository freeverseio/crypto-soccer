const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const delegateUtils = require('../utils/delegateCallUtils.js');

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
    
    function getIdxInABI(abi, name) {
        for (i = 0; i < abi.length; i++) { 
            if (abi[i].name == name) {
                return i;
            }
        }    
    }
    
    beforeEach(async () => {
        sto = await StorageProxy.new().should.be.fulfilled;
        assets = await Assets.at(sto.address).should.be.fulfilled;
        assetsAsLib = await Assets.new().should.be.fulfilled;
    });
    
    it('deploy storage by adding Assets selectors', async () => {
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        nSelectorsPerContract = [selectors.length];
        addresses = [assetsAsLib.address];
        requiresPermission = [false];
        names = [web3.utils.utf8ToHex('Assets')];
        tx = await sto.deployNewStorageProxies(nSelectorsPerContract, selectors, addresses, requiresPermission, names).should.be.fulfilled;
        // note that contractId = 0 is the null one
        truffleAssert.eventEmitted(tx, "ContractSet", async (event) => { return event.contractId === 1 && event.names === names });
    });

    it('call init() function inside Assets via delegate call from declaring ALL selectors in Assets', async () => {
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        nSelectorsPerContract = [selectors.length];
        addresses = [assetsAsLib.address];
        requiresPermission = [false];
        names = [web3.utils.utf8ToHex('Assets')];
        await assets.init().should.be.rejected;
        tx = await sto.deployNewStorageProxies(nSelectorsPerContract, selectors, addresses, requiresPermission, names).should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);

        // I can redeploy, and, because storage is preserved, I cannot init again, but nCountries is still OK
        tx = await sto.deployNewStorageProxies(nSelectorsPerContract, selectors, addresses, requiresPermission, names).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);
        
        // If I redeploy but removing the functions in assets, even the getter fails
        tx = await sto.deployNewStorageProxies([], [], [], [], []).should.be.fulfilled;
        result = await assets.countCountries(tz = 1).should.be.rejected;
        
    });

    return
    
    it('deleteSelectors', async () => {
        await sto.addNewContract(addr = assetsAsLib.address, isSet = false, name = "Assets").should.be.fulfilled;
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        await sto.addNewSelectors(selectors, contractId = 1).should.be.fulfilled;
        await sto.deleteSelectors(selectors).should.be.fulfilled;
        await assets.init().should.be.rejected;
    });

    it('call init() function inside Assets via delegate call', async () => {
        await sto.addNewContract(addr = assetsAsLib.address, isSet = false, name = "Assets").should.be.fulfilled;
        initPosInAbi = getIdxInABI(Assets.abi, "init");
        // we first add a function different from init(), and show that we cannot call assets.init()
        selector = web3.eth.abi.encodeFunctionSignature(Assets.abi[initPosInAbi - 1])
        await sto.addNewSelectors([selector], contractId = 1).should.be.fulfilled;
        await assets.init().should.be.rejected;
        // we add the correct function and show that we can call it now:
        selector = web3.eth.abi.encodeFunctionSignature(Assets.abi[initPosInAbi])
        await sto.addNewSelectors([selector], contractId = 1).should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        // we show that storage in StorageProxy has been updated (wasInit = true), so we cannot init twice:
        await assets.init().should.be.rejected ;
        // we show that the getter in assets works well when called from storageProxy:
        wasInitPosInAbi = getIdxInABI(Assets.abi, "_wasInited");
        selector = web3.eth.abi.encodeFunctionSignature(Assets.abi[wasInitPosInAbi])
        await sto.addNewSelectors([selector], contractId = 1).should.be.fulfilled;
        isInit = await assets._wasInited().should.be.fulfilled;
        isInit.should.be.equal(true);
    });

    it('deploy correctly', async () => {
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        result = await sto.countFunctions().should.be.fulfilled;
        result.toNumber().should.be.equal(0);
    });

    it('add contract info and one function', async () => {
        await sto.addNewContract(addr = assets.address, isSet = false, name = "Assets").should.be.fulfilled;
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(2);
        for (c = 1; c < result.toNumber(); c++) { 
            let {0: ad, 1: isSetter, 2: nom} = await sto.getContractInfo(c).should.be.fulfilled;
            ad.should.be.equal(addr);
            isSetter.should.be.equal(isSet);
            nom.should.be.equal(name);
        }
        selector = web3.eth.abi.encodeFunctionSignature('setNewAsset(uint256,string)')
        await sto.addNewSelectors([selector], contractId = 1).should.be.fulfilled;
        result = await sto.countFunctions().should.be.fulfilled;
        result.toNumber().should.be.equal(1);
    });

    it('autorizations check', async () => {
        await sto.setStorageOwner(addr = ALICE, {from: ALICE}).should.be.rejected;
        await sto.setStorageOwner(addr = ALICE, {from: FREEVERSE}).should.be.fulfilled;
        await sto.addNewContract(addr = assets.address, isSet = false, name = "Assets", {from: FREEVERSE}).should.be.rejected;
        await sto.addNewContract(addr = assets.address, isSet = false, name = "Assets", {from: ALICE}).should.be.fulfilled;
        selector = web3.eth.abi.encodeFunctionSignature('setNewAsset(uint256,string)')
        await sto.addNewSelectors([selector], contractId = 1, {from: FREEVERSE}).should.be.rejected;
        await sto.addNewSelectors([selector], contractId = 1, {from: ALICE}).should.be.fulfilled;
    });


});