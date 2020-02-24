const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const delegateUtils = require('../utils/delegateCallUtils.js');

const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');

contract('Proxy', (accounts) => {
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
        proxy = await Proxy.new().should.be.fulfilled;
        assets = await Assets.at(proxy.address).should.be.fulfilled;
        assetsAsLib = await Assets.new().should.be.fulfilled;
    });
    
    it('permissions check on full deploy: everyone can, currently, until we set restrictions inside Assets contract', async () => {
        depl = await delegateUtils.deployDelegate(proxy, Assets, Market, Updates);
        assets = depl[0]
        await assets.init({from: ALICE}).should.be.fulfilled;
        await assets.getNCountriesInTZ(tz = 1, {from: ALICE}).should.be.fulfilled;
        tz = 1;
        countryIdxInTZ = 0;
        teamIdxInCountry = 0;
        teamId = await assets.encodeTZCountryAndVal(tz, countryIdxInTZ, teamIdxInCountry);
        await assets.transferFirstBotToAddr(tz, countryIdxInTZ, ALICE, {from: FREEVERSE}).should.be.fulfilled;
    });

    it('deploy storage by adding Assets selectors', async () => {
        // contact[0] is the NULL contract
        result = await proxy.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(1);
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        tx0 = await proxy.addContract(contractId = 0, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.rejected;
        tx0 = await proxy.addContract(contractId = 2, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.rejected;
        contractId = 1;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        truffleAssert.eventEmitted(tx0, "ContractAdded", async (event) => { return event.contractId === contractId && event.name === name});
        
        tx1 = await proxy.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        truffleAssert.eventEmitted(tx1, "ContractsActivated", async (event) => { return event.contractId === contractId });

        result = await proxy.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(2);
    });

    it('call init() function inside Assets via delegate call from declaring ALL selectors in Assets', async () => {
        await assets.init().should.be.rejected;

        // add function (still not enough to call assets):
        selectors = delegateUtils.extractSelectorsFromAbi(Assets.abi);
        contractId = 1;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        await assets.init().should.be.rejected;
        // activate function, now, enough to call assets:
        tx1 = await proxy.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);

        // test that deleteContracts destroys all calls to assets functions
        tx1 = await proxy.deleteContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.rejected;

        // I can re-activate, and, because storage is preserved, I cannot init again, but nCountries is still OK
        contractId = 2;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        tx1 = await proxy.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);

        // I can do the same thing in one atomic TX:
        contractId = 3;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        tx1 = await proxy.deleteAndActivateContracts(deactivate = [2], activate = [3]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);
    });
    
    
});