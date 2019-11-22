const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');
const Evolution = artifacts.require('Evolution');
const Assets = artifacts.require('Assets');
const Leagues = artifacts.require('Leagues');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Friendlies = artifacts.require('Friendlies');


require('chai')
    .use(require('chai-as-promised'))
    .should();

module.exports = function (deployer) {
  deployer.then(async () => {
    const engine = await deployer.deploy(Engine).should.be.fulfilled;
    const enginePreComp = await deployer.deploy(EnginePreComp).should.be.fulfilled;
    const evolution= await deployer.deploy(Evolution).should.be.fulfilled;
    const assets = await deployer.deploy(Assets).should.be.fulfilled;
    const leagues = await deployer.deploy(Leagues).should.be.fulfilled;
    const market = await deployer.deploy(Market).should.be.fulfilled;
    const updates = await deployer.deploy(Updates).should.be.fulfilled;
    const friendlies = await deployer.deploy(Friendlies).should.be.fulfilled;

    console.log("Setting up ...");
    await leagues.setEngineAdress(engine.address).should.be.fulfilled;
    await market.setAssetsAddress(assets.address).should.be.fulfilled;
    await updates.initUpdates(assets.address).should.be.fulfilled;
    await evolution.setAssetsAddress(assets.address).should.be.fulfilled;
    await evolution.setEngine(engine.address).should.be.fulfilled;
    await engine.setPreCompAddr(enginePreComp.address).should.be.fulfilled;
    console.log("Setting up ... done");

    console.log("Initing ... TODO : only one zone actually");
    await assets.initSingleTZ(1).should.be.fulfilled; // TODO: bootstrap od all timezone using init()
    console.log("Initing ... done");

    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("------------------------");
    console.log("ENGINE_CONTRACT_ADDRESS=" + engine.address);
    console.log("ENGINEPRECOMP_CONTRACT_ADDRESS=" + enginePreComp.address);
    console.log("LEAGUES_CONTRACT_ADDRESS=" + leagues.address);
    console.log("MARKET_CONTRACT_ADDRESS=" + market.address);
    console.log("UPDATES_CONTRACT_ADDRESS=" + updates.address);
    console.log("ASSETS_CONTRACT_ADDRESS=" + assets.address);
    console.log("EVOLUTION_CONTRACT_ADDRESS=" + evolution.address);
  });
};

