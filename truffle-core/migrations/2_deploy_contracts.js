const Storage = artifacts.require('Storage');
const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');
const EngineApplyBoosters = artifacts.require('EngineApplyBoosters');
const TrainingPoints = artifacts.require('TrainingPoints');
const Evolution = artifacts.require('Evolution');
const Assets = artifacts.require('Assets');
const Leagues = artifacts.require('Leagues');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const Friendlies = artifacts.require('Friendlies');
const Shop = artifacts.require('Shop');
const Privileged = artifacts.require('Privileged');
const Utils = artifacts.require('Utils');
const PlayAndEvolve = artifacts.require('PlayAndEvolve');


require('chai')
    .use(require('chai-as-promised'))
    .should();

module.exports = function (deployer) {
  deployer.then(async () => {
    const storage = await deployer.deploy(Storage).should.be.fulfilled;
    const engine = await deployer.deploy(Engine).should.be.fulfilled;
    const enginePreComp = await deployer.deploy(EnginePreComp).should.be.fulfilled;
    const engineApplyBoosters = await deployer.deploy(EngineApplyBoosters).should.be.fulfilled;
    const trainingPoints= await deployer.deploy(TrainingPoints).should.be.fulfilled;
    const evolution= await deployer.deploy(Evolution).should.be.fulfilled;
    const assets = await deployer.deploy(Assets).should.be.fulfilled;
    const leagues = await deployer.deploy(Leagues).should.be.fulfilled;
    const market = await deployer.deploy(Market).should.be.fulfilled;
    const updates = await deployer.deploy(Updates).should.be.fulfilled;
    const friendlies = await deployer.deploy(Friendlies).should.be.fulfilled;
    const shop = await deployer.deploy(Shop).should.be.fulfilled;
    const privileged = await deployer.deploy(Privileged).should.be.fulfilled;
    const utils = await deployer.deploy(Utils).should.be.fulfilled;
    const playAndEvolve = await deployer.deploy(PlayAndEvolve).should.be.fulfilled;
    
    console.log("Setting up ...");
    await storage.setAssetsOwner(assets.address).should.be.fulfilled;
    await storage.setUpdatesOwner(updates.address).should.be.fulfilled;
    await storage.setMarketOwner(market.address).should.be.fulfilled;
    await leagues.setAssetsAdress(assets.address).should.be.fulfilled;
    await market.setAssetsAddress(assets.address).should.be.fulfilled;
    await market.setStorageAddress(storage.address).should.be.fulfilled;
    await updates.initUpdates(storage.address).should.be.fulfilled;
    await trainingPoints.setAssetsAddress(assets.address).should.be.fulfilled;
    await engine.setPreCompAddr(enginePreComp.address).should.be.fulfilled;
    await engine.setApplyBoostersAddr(engineApplyBoosters.address).should.be.fulfilled;
    await market.setAcademyAddr("0xb8CE9ab6943e0eCED004cDe8e3bBed6568B2Fa01");
    await playAndEvolve.setTrainingAddress(trainingPoints.address);
    await playAndEvolve.setEvolutionAddress(evolution.address).should.be.fulfilled;
    await playAndEvolve.setEngineAddress(engine.address).should.be.fulfilled;
    await playAndEvolve.setShopAddress(shop.address).should.be.fulfilled;

    console.log("Setting up ... done");
    if (deployer.network === "production") {
      await assets.init(storage.address).should.be.fulfilled;
    } else {
      const timezone = 1;
      console.log("Initing only timezone " + timezone)
      await assets.initSingleTZ(timezone, storage.address).should.be.fulfilled; // TODO: bootstrap od all timezone using init()
    }
    console.log("Initing ... done");

    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("------------------------");
    console.log("STORAGE_CONTRACT_ADDRESS=" + storage.address);
    console.log("ENGINE_CONTRACT_ADDRESS=" + engine.address);
    console.log("ENGINEPRECOMP_CONTRACT_ADDRESS=" + enginePreComp.address);
    console.log("ENGINEAPPLYBOOSTERS_CONTRACT_ADDRESS=" + engineApplyBoosters.address);
    console.log("LEAGUES_CONTRACT_ADDRESS=" + leagues.address);
    console.log("MARKET_CONTRACT_ADDRESS=" + market.address);
    console.log("UPDATES_CONTRACT_ADDRESS=" + updates.address);
    console.log("ASSETS_CONTRACT_ADDRESS=" + assets.address);
    console.log("TRAININGPOINTS_CONTRACT_ADDRESS=" + trainingPoints.address);
    console.log("EVOLUTION_CONTRACT_ADDRESS=" + evolution.address);
    console.log("FRIENDLIES_CONTRACT_ADDRESS=" + friendlies.address);
    console.log("SHOP_CONTRACT_ADDRESS=" + shop.address);
    console.log("PRIVILEGED_CONTRACT_ADDRESS=" + privileged.address);
    console.log("UTILS_CONTRACT_ADDRESS=" + utils.address);
    console.log("PLAYANDEVOLVE_CONTRACT_ADDRESS=" + playAndEvolve.address);
  });
};

