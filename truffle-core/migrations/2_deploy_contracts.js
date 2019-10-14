const Engine = artifacts.require('Engine');
const Leagues = artifacts.require('Leagues');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
require('chai')
    .use(require('chai-as-promised'))
    .should();

module.exports = function (deployer) {
  deployer.then(async () => {
    const engine = await deployer.deploy(Engine).should.be.fulfilled;
    const leagues = await deployer.deploy(Leagues).should.be.fulfilled;
    const market = await deployer.deploy(Market).should.be.fulfilled;
    const updates = await deployer.deploy(Updates).should.be.fulfilled;

    console.log("Setting up ...");
    await leagues.setEngineAdress(engine.address).should.be.fulfilled;
    await market.setAssetsAddress(leagues.address).should.be.fulfilled;
    console.log("Setting up ... done");

    console.log("Initing ...");
    await leagues.initSingleTZ(1).should.be.fulfilled; // TODO: bootstrap od all timezone using init()
    await updates.initUpdates(leagues.address).should.be.fulfilled;
    console.log("Initing ... done");

    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("------------------------");
    console.log("ENGINE_CONTRACT_ADDRESS=" + engine.address);
    console.log("LEAGUES_CONTRACT_ADDRESS=" + leagues.address);
    console.log("MARKET_CONTRACT_ADDRESS=" + market.address);
    console.log("UPDATES_CONTRACT_ADDRESS=" + updates.address);
  });
};

