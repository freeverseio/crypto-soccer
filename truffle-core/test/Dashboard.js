const BN = require('bn.js');
var fs = require('fs');

require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const deployUtils = require('../utils/deployUtils.js');

const Dashboard = artifacts.require('Dashboard');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Challenges = artifacts.require('Challenges');


contract('Dashboard', (accounts) => {
    const it2 = async(text, f) => {};

    beforeEach(async () => {
        dashboard = await Dashboard.new().should.be.fulfilled;
        defaultSetup = deployUtils.getDefaultSetup(accounts);
        owners = defaultSetup.owners;
        depl = await deployUtils.deploy(versionNumber = 0, owners, Proxy, proxyAddress = '0x0', Assets, Market, Updates, Challenges);
        [proxy, assets, market, updates] = depl;
        await deployUtils.setProxyContractOwners(proxy, assets, owners, owners.company).should.be.fulfilled;
        proxyLib = await Proxy.at(dashboard.address).should.be.fulfilled;
    });

    it('setCOO', async () => {
        console.log("--a");
        await dashboard.setProxy(proxy.address).should.be.fulfilled;
        console.log("--b");
        await proxyLib.setSuperUser(market, {from: owners.company}).should.be.fulfilled;
        console.log("--c");
    });

});