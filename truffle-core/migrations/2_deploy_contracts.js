const Market = artifacts.require('Market');
const Assets = artifacts.require('Assets');
const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');
const EngineApplyBoosters = artifacts.require('EngineApplyBoosters');
const TrainingPoints = artifacts.require('TrainingPoints');
const Evolution = artifacts.require('Evolution');
const Leagues = artifacts.require('Leagues');
const Updates = artifacts.require('Updates');
const Friendlies = artifacts.require('Friendlies');
const Shop = artifacts.require('Shop');
const Privileged = artifacts.require('Privileged');
const Utils = artifacts.require('Utils');
const PlayAndEvolve = artifacts.require('PlayAndEvolve');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');

require('chai')
    .use(require('chai-as-promised'))
    .should();

const delegateUtils = require('../utils/delegateCallUtils.js');

module.exports = function (deployer, network, accounts) {
  deployer.then(async () => {
    const proxy = await deployer.deploy(Proxy).should.be.fulfilled;
    const {0: assets, 1: market, 2: updates} = await delegateUtils.deployDelegate(proxy, Assets, Market, Updates);

    const engine = await deployer.deploy(Engine).should.be.fulfilled;
    const enginePreComp = await deployer.deploy(EnginePreComp).should.be.fulfilled;
    const engineApplyBoosters = await deployer.deploy(EngineApplyBoosters).should.be.fulfilled;
    const trainingPoints= await deployer.deploy(TrainingPoints).should.be.fulfilled;
    const evolution= await deployer.deploy(Evolution).should.be.fulfilled;
    const leagues = await deployer.deploy(Leagues).should.be.fulfilled;
    const friendlies = await deployer.deploy(Friendlies).should.be.fulfilled;
    const shop = await deployer.deploy(Shop).should.be.fulfilled;
    const privileged = await deployer.deploy(Privileged).should.be.fulfilled;
    const utils = await deployer.deploy(Utils).should.be.fulfilled;
    const playAndEvolve = await deployer.deploy(PlayAndEvolve).should.be.fulfilled;
    const constantsGetters = await deployer.deploy(ConstantsGetters).should.be.fulfilled;
    
    console.log("Setting up ...");
    await leagues.setEngineAdress(engine.address).should.be.fulfilled;
    await leagues.setAssetsAdress(assets.address).should.be.fulfilled;
    await updates.initUpdates().should.be.fulfilled;
    await trainingPoints.setAssetsAddress(assets.address).should.be.fulfilled;
    await trainingPoints.setMarketAddress(market.address).should.be.fulfilled;
    await engine.setPreCompAddr(enginePreComp.address).should.be.fulfilled;
    await engine.setApplyBoostersAddr(engineApplyBoosters.address).should.be.fulfilled;
    await playAndEvolve.setTrainingAddress(trainingPoints.address);
    await playAndEvolve.setEvolutionAddress(evolution.address).should.be.fulfilled;
    await playAndEvolve.setEngineAddress(engine.address).should.be.fulfilled;
    await playAndEvolve.setShopAddress(shop.address).should.be.fulfilled;

    const value = "100000000000000000";
    const to = "0xeb3ce112d8610382a994646872c4361a96c82cf8";

    console.log("Setting up ... done");
    if (deployer.network === "xdai") {
      await assets.init().should.be.fulfilled;
    } else if (deployer.network === "local") {
      console.log("WARNING ... only timezone 1")
      await assets.initSingleTZ(1).should.be.fulfilled;
      console.log("Transfer " + value + " to " + to);
      web3.eth.sendTransaction({ from: accounts[0], to, value }).should.be.fulfilled;
    } else {
      console.log("WARNING ... only timezone 10")
      await assets.initSingleTZ(10).should.be.fulfilled;
      console.log("Transfer " + value + " to " + to);
      web3.eth.sendTransaction({ from: accounts[0], to, value }).should.be.fulfilled;
    }
    console.log("Initing ... done");

    await assets.setAcademyAddr("0x7c34471e39c4A4De223c05DF452e28F0c4BD9BF0");

    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("-----------AddressesStart-----------");
    console.log("ASSETS_CONTRACT_ADDRESS=" + assets.address);
    console.log("MARKET_CONTRACT_ADDRESS=" + market.address);
    console.log("ENGINE_CONTRACT_ADDRESS=" + engine.address);
    console.log("ENGINEPRECOMP_CONTRACT_ADDRESS=" + enginePreComp.address);
    console.log("ENGINEAPPLYBOOSTERS_CONTRACT_ADDRESS=" + engineApplyBoosters.address);
    console.log("LEAGUES_CONTRACT_ADDRESS=" + leagues.address);
    console.log("UPDATES_CONTRACT_ADDRESS=" + updates.address);
    console.log("TRAININGPOINTS_CONTRACT_ADDRESS=" + trainingPoints.address);
    console.log("EVOLUTION_CONTRACT_ADDRESS=" + evolution.address);
    console.log("FRIENDLIES_CONTRACT_ADDRESS=" + friendlies.address);
    console.log("SHOP_CONTRACT_ADDRESS=" + shop.address);
    console.log("PRIVILEGED_CONTRACT_ADDRESS=" + privileged.address);
    console.log("UTILS_CONTRACT_ADDRESS=" + utils.address);
    console.log("PLAYANDEVOLVE_CONTRACT_ADDRESS=" + playAndEvolve.address);
    console.log("CONSTANTSGETTERS_CONTRACT_ADDRESS=" + constantsGetters.address);
    console.log("-----------AddressesEnd-----------");
  });
};

