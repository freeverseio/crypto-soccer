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
const Market = artifacts.require('Market');
const AssetsView = artifacts.require('AssetsView');
const MarketView = artifacts.require('MarketView');

contract('StorageProxy', (accounts) => {
    const FREEVERSE = accounts[0];
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const it2 = async(text, f) => {};

    function toBytes32(name) { return web3.utils.utf8ToHex(name); }

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
    
    it('permissions check on full deploy', async () => {
        depl = await delegateUtils.deployDelegate(
            StorageProxy, 
            Assets, 
            AssetsView, 
            Market, 
            MarketView
        );
        assets = depl[0]
        // ALICE can execute any view function (in AssetsView)
        await assets.init({from: ALICE}).should.be.rejected;
        await assets.init({from: FREEVERSE}).should.be.fulfilled;
        await assets.getNCountriesInTZ(tz = 1, {from: ALICE}).should.be.fulfilled;
        // ALICE cannot execute any write function (in Assets), FREEVERSE can
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0;
        teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE, {from: BOB}).should.be.rejected;
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE, {from: FREEVERSE}).should.be.fulfilled;
    });

    it('deploy storage by adding Assets selectors', async () => {
        // contact[0] is the NULL contract
        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        tx0 = await sto.addContract(contractId = 0, assetsAsLib.address, requiresPermission = false, selectors, name = toBytes32("Assets")).should.be.rejected;
        tx0 = await sto.addContract(contractId = 2, assetsAsLib.address, requiresPermission = false, selectors, name = toBytes32("Assets")).should.be.rejected;
        contractId = 1;
        tx0 = await sto.addContract(contractId, assetsAsLib.address, requiresPermission = false, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        truffleAssert.eventEmitted(tx0, "ContractAdded", async (event) => { return event.contractId === contractId && event.name === name});
        
        tx1 = await sto.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        truffleAssert.eventEmitted(tx1, "ContractsActivated", async (event) => { return event.contractId === contractId });

        result = await sto.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(2);
    });

    it('call init() function inside Assets via delegate call from declaring ALL selectors in Assets', async () => {
        await assets.init().should.be.rejected;

        // add function (still not enough to call assets):
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        contractId = 1;
        tx0 = await sto.addContract(contractId, assetsAsLib.address, requiresPermission = false, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        await assets.init().should.be.rejected;
        // activate function, now, enough to call assets:
        tx1 = await sto.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);

        // test that deleteContracts destroys all calls to assets functions
        tx1 = await sto.deleteContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.rejected;

        // I can re-activate, and, because storage is preserved, I cannot init again, but nCountries is still OK
        contractId = 2;
        tx0 = await sto.addContract(contractId, assetsAsLib.address, requiresPermission = false, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        tx1 = await sto.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);

        // I can do the same thing in one atomic TX:
        contractId = 3;
        tx0 = await sto.addContract(contractId, assetsAsLib.address, requiresPermission = false, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        tx1 = await sto.deleteAndActivateContracts(deactivate = [2], activate = [3]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);
    });
    
    
});