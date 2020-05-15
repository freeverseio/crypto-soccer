const BN = require('bn.js');
var fs = require('fs');

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const deployUtils = require('../utils/deployUtils.js');

const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Challenges = artifacts.require('Challenges');

contract('Proxy', (accounts) => {
    const FREEVERSE = accounts[0];
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const it2 = async(text, f) => {};

    function toBytes32(name) { return web3.utils.utf8ToHex(name); }
    function fromBytes32(name) { return web3.utils.hexToUtf8(name); }

    
    function getIdxInABI(abi, name) {
        for (i = 0; i < abi.length; i++) { 
            if (abi[i].name == name) {
                return i;
            }
        }    
    }
    
    beforeEach(async () => {
        defaultSetup = deployUtils.getDefaultSetup(accounts);
        // let { singleTimezone, owners, requiredStake } = deployUtils.getDefaultSetup(accounts);
        proxy = await Proxy.new(defaultSetup.owners.company, defaultSetup.owners.superuser, deployUtils.extractSelectorsFromAbi(Proxy.abi)).should.be.fulfilled;
        assets = await Assets.at(proxy.address).should.be.fulfilled;
        assetsAsLib = await Assets.new(defaultSetup.owners.market).should.be.fulfilled;
    });

    it('fails when adding a contract to an address without contract', async () => {
        await proxy.addContract(contractId = 1, '0x0', selectors, name = toBytes32("Assets")).should.be.rejected;
        await proxy.addContract(contractId = 1, '0x32132', selectors, name = toBytes32("Assets")).should.be.rejected;
        await proxy.addContract(contractId = 1, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
    });

    it('permissions check to change owner of proxy', async () => {
        await proxy.proposeSuperUser(BOB, {from: ALICE}).should.be.rejected;
        await proxy.proposeSuperUser(BOB, {from: FREEVERSE}).should.be.fulfilled;
        await proxy.proposeSuperUser(ALICE, {from: ALICE}).should.be.rejected;
        await proxy.acceptSuperUser({from: ALICE}).should.be.rejected;
        await proxy.acceptSuperUser({from: BOB}).should.be.fulfilled;
        await proxy.proposeSuperUser(ALICE, {from: FREEVERSE}).should.be.rejected;
        await proxy.proposeSuperUser(ALICE, {from: BOB}).should.be.fulfilled;
    });

    it('full deploy should work', async () => {
        const {0: prox, 1: ass, 2: mkt, 3: updt, 4: chll} = await deployUtils.deploy(versionNumber = 0, defaultSetup.owners, Proxy, proxyAddress = '0x0', Assets, Market, Updates, Challenges);
    });
    
    it('permissions check on full deploy: everyone can call delegates, currently, until we set restrictions inside Assets contract', async () => {
        depl = await deployUtils.deploy(versionNumber = 0, defaultSetup.owners, Proxy, proxyAddress = '0x0', Assets, Market, Updates, Challenges);
        assets = depl[1]
        await assets.init({from: ALICE}).should.be.fulfilled;
        await assets.countCountries(tz = 1, {from: ALICE}).should.be.fulfilled;
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
        selectors = deployUtils.extractSelectorsFromAbi(Assets.abi);
        tx0 = await proxy.addContract(contractId = 0, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.rejected;
        tx0 = await proxy.addContract(contractId = 2, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.rejected;
        contractId = 1;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;

        truffleAssert.eventEmitted(tx0, "ContractAdded", (event) => {
            ok = true;
            for (s = 0; s < selectors.length; s++) {
                ok = ok && (event.selectors[s] == selectors[s]);
            }
            return ok && event.contractId.toNumber().should.be.equal(contractId) && fromBytes32(event.name).should.be.equal("Assets");
        });


        var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxy.getContractInfo(contractId).should.be.fulfilled;
        isActive.should.be.equal(false);
        addr.should.be.equal(assetsAsLib.address);

        
        tx1 = await proxy.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        truffleAssert.eventEmitted(tx1, "ContractsActivated", (event) => { 
            return event.contractIds[0].toNumber().should.be.equal(contractId)
        });
        var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxy.getContractInfo(contractId).should.be.fulfilled;
        isActive.should.be.equal(true);

        result = await proxy.countContracts().should.be.fulfilled;
        result.toNumber().should.be.equal(2);
        
    });

    it('call init() function inside Assets via delegate call from declaring ALL selectors in Assets', async () => {
        await assets.init().should.be.rejected;

        // add function (still not enough to call assets):
        selectors = deployUtils.extractSelectorsFromAbi(Assets.abi);
        contractId = 1;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        await assets.init().should.be.rejected;
        // activate function, now, enough to call assets:
        tx1 = await proxy.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);

        // test that deactivateContracts destroys all calls to assets functions
        tx1 = await proxy.deactivateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.rejected;

        // I can re-activate, and, because storage is preserved, I cannot init again, but nCountries is still OK
        contractId = 2;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        tx1 = await proxy.activateContracts(contractIds = [contractId]).should.be.fulfilled;
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);
        var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxy.getContractInfo(contractId).should.be.fulfilled;
        isActive.should.be.equal(true);

        // I can do the same thing in one atomic TX:
        contractId = 3;
        tx0 = await proxy.addContract(contractId, assetsAsLib.address, selectors, name = toBytes32("Assets")).should.be.fulfilled;
        tx1 = await proxy.deactivateAndActivateContracts(deactivate = [2], activate = [3]).should.be.fulfilled;

        now = Math.floor(Date.now()/1000);
        truffleAssert.eventEmitted(tx1, "ContractsActivated", (event) => { 
            console.log();
            return event.contractIds[0].toNumber().should.be.equal(3) && (Math.abs(event.time.toNumber()-now) < 30).should.be.equal(true)
        });
        truffleAssert.eventEmitted(tx1, "ContractsDeactivated", (event) => { 
            return event.contractIds[0].toNumber().should.be.equal(2) && (Math.abs(event.time.toNumber()-now) < 30).should.be.equal(true)
        });

        var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxy.getContractInfo(2).should.be.fulfilled;
        isActive.should.be.equal(false);
        await assets.init().should.be.rejected;
        result = await assets.countCountries(tz = 1).should.be.fulfilled;
        (result.toNumber() > 0).should.be.equal(true);
    });
    
    it('deploy and redeploy', async () => {
        // contact[0] is the NULL contract
        const nContractsToProxy = 4;
        assert.equal(await proxy.countContracts(), '1', "wrong init number of contracts in proxy");
        const {0: proxyV0, 1: assV0, 2: markV0, 3: updV0, 4: chllV0} = await deployUtils.deploy(versionNumber = 0, defaultSetup.owners, Proxy, proxyAddress = '0x0', Assets, Market, Updates, Challenges);
        assert.equal(await proxy.countContracts(), '1', "wrong init number of contracts in proxy");
        assert.equal(await proxyV0.countContracts(), '5', "wrong V0 number of contracts in proxy");

        expectedNamesV0 = ['Assets0', 'Market0', 'Updates0', 'Challenges0'];
        for (c = 1; c < 1+nContractsToProxy; c++) {
            var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxyV0.getContractInfo(c).should.be.fulfilled;
            isActive.should.be.equal(true);
            assert(fromBytes32(nom) == expectedNamesV0[c-1] , "wrong contract name");
        }    

        const {0: proxyV1, 1: assV1, 2: markV1, 3: updV1, 4: chllV1} = await deployUtils.deploy(versionNumber = 1, defaultSetup.owners, Proxy, proxyV0.address, Assets, Market, Updates, Challenges);
        assert.equal(await proxyV1.address, proxyV0.address);
        assert.equal(await proxyV0.countContracts(), '9', "wrong V1 number of contracts in proxyV0");
        assert.equal(await proxyV1.countContracts(), '9', "wrong V1 number of contracts in proxyV1");
        assert.equal(await assV1.address, assV0.address);
        assert.equal(await markV1.address, markV0.address);
        assert.equal(await updV1.address, updV0.address);
        assert.equal(await chllV1.address, chllV0.address);
        expectedNamesV1 = ['Assets1', 'Market1', 'Updates1', 'Challenges1'];
        for (c = 1; c < 1+nContractsToProxy; c++) {
            var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxyV1.getContractInfo(c).should.be.fulfilled;
            isActive.should.be.equal(false);
            assert(fromBytes32(nom) == expectedNamesV0[c-1] , "wrong contract name");
            var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxyV1.getContractInfo(c+nContractsToProxy).should.be.fulfilled;
            isActive.should.be.equal(true);
            assert(fromBytes32(nom) == expectedNamesV1[c-1] , "wrong contract name");
        }    
        

        
    });

    
    
});