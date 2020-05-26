/*
 Tests for all functions in Directoy.sol
*/
const BN = require('bn.js');
var fs = require('fs');

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const deployUtils = require('../utils/deployUtils.js');

const Directory = artifacts.require('Directory');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Challenges = artifacts.require('Challenges');

contract('Directory', (accounts) => {
    const FREEVERSE = accounts[0];
    const ALICE = accounts[1];
    const BOB = accounts[2];
    const CAROL = accounts[3];
    const it2 = async(text, f) => {};

    function toBytes32(name) { return web3.utils.utf8ToHex(name); }
    function fromBytes32(name) { return web3.utils.hexToUtf8(name); }

    beforeEach(async () => {
        defaultSetup = deployUtils.getDefaultSetup(accounts);
        owners = defaultSetup.owners;
        depl = await deployUtils.deploy(owners, Proxy, Assets, Market, Updates, Challenges);
        [proxy, assets, market, updates] = depl;
        await deployUtils.setProxyContractOwners(proxy, assets, owners, owners.company).should.be.fulfilled;
        directory = await Directory.new(proxy.address).should.be.fulfilled;
    });
    
    it('deploy tests', async () => {
        names = ["Baby1", "Baby2_Weird"];
        names32 = []
        for (n = 0; n < names.length; n++) names32.push(toBytes32(names[n]));        
        addresses = [ALICE, BOB];
        console.log(directory.address, proxy.address);
        await directory.deploy(names32, addresses, {from: owners.relay}).should.be.rejected;
        await directory.deploy(names32, addresses, {from: owners.COO}).should.be.fulfilled;
        // result = await directory.getAddress(toBytes32("Baby2_Weird")).should.be.fulfilled;
        // result.should.be.equal(BOB);
        var {0: noms, 1: addr} = await directory.getDirectory().should.be.fulfilled;
        assert.equal(noms.length, 0, "there should be no contract until activated");
        assert.equal(addr.length, 0, "there should be no contract until activated");

        await directory.activateNewDeploy({from: owners.COO}).should.be.fulfilled;
        console.log("++");        
        var {0: noms, 1: addr} = await directory.getDirectory().should.be.fulfilled;

        console.log(noms);        
        debug.compareArrays(addr, addresses, toNum = false);
        for (n = 0; n < noms.length; n++) noms[n] = fromBytes32(noms[n]);        
        debug.compareArrays(noms, names, toNum = false);
        
        // check events
    });

    
});