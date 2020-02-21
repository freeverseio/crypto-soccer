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
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(0);
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        nSelectorsPerContract = [selectors.length];
        addresses = [assetsAsLib.address];
        requiresPermission = [false];
        names = [web3.utils.utf8ToHex('Assets')];
        tx = await sto.deployNewStorageProxies(nSelectorsPerContract, selectors, addresses, requiresPermission, names).should.be.fulfilled;
        truffleAssert.eventEmitted(tx, "ContractSet", async (event) => { return event.contractId === 1 && event.names === names });
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(1);

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
});