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
const Directory = artifacts.require('Directory');

require('chai')
    .use(require('chai-as-promised'))
    .should();

const delegateUtils = require('../utils/delegateCallUtils.js');

module.exports = function (deployer, network, accounts) {
  deployer.then(async () => {
    
    const {0: assets, 1: market, 2: updates} = await delegateUtils.firstDeploy(versionNumber = 0, deployer, Proxy, proxyAddress = "0x0", Assets, Market, Updates);
    // const {0: assets, 1: market, 2: updates} = await delegateUtils.firstDeploy(versionNumber = 1, deployer, Proxy, proxyAddress, Assets, Market, Updates);

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
    const directory = await deployer.deploy(Directory).should.be.fulfilled;
    
    console.log("Setting up ...");
    await leagues.setEngineAdress(engine.address).should.be.fulfilled;
    await leagues.setAssetsAdress(assets.address).should.be.fulfilled;
    await trainingPoints.setAssetsAddress(assets.address).should.be.fulfilled;
    await trainingPoints.setMarketAddress(market.address).should.be.fulfilled;
    await engine.setPreCompAddr(enginePreComp.address).should.be.fulfilled;
    await engine.setApplyBoostersAddr(engineApplyBoosters.address).should.be.fulfilled;
    await playAndEvolve.setTrainingAddress(trainingPoints.address);
    await playAndEvolve.setEvolutionAddress(evolution.address).should.be.fulfilled;
    await playAndEvolve.setEngineAddress(engine.address).should.be.fulfilled;
    await playAndEvolve.setShopAddress(shop.address).should.be.fulfilled;

    namesAndAddresses = [
      ["ASSETS", assets.address],
      ["MARKET", market.address],
      ["ENGINE", engine.address],
      ["ENGINEPRECOMP", enginePreComp.address],
      ["ENGINEAPPLYBOOSTERS", engineApplyBoosters.address],
      ["LEAGUES", leagues.address],
      ["UPDATES", updates.address],
      ["TRAININGPOINTS", trainingPoints.address],
      ["EVOLUTION", evolution.address],
      ["TRAININGPOINTS", trainingPoints.address],
      ["FRIENDLIES", friendlies.address],
      ["SHOP_CONTRACT", shop.address],
      ["PRIVILEGED", privileged.address],
      ["UTILS", utils.address],
      ["PRIVILEGED", assets.address],
      ["PLAYANDEVOLVE", playAndEvolve.address],
      ["CONSTANTSGETTERS", constantsGetters.address]
    ]

    // Build arrays "names" and "addresses" and store in Directory contract
    names = [];
    namesBytes32 = [];
    addresses = [];
    for (c = 0; c < namesAndAddresses.length; c++) {
      names.push(namesAndAddresses[c][0]);
      namesBytes32.push(web3.utils.utf8ToHex(namesAndAddresses[c][0]));
      addresses.push(namesAndAddresses[c][1]);
    }
    await directory.deploy(namesBytes32, addresses).should.be.fulfilled;

    // Print Summary to Console
    namesAndAddresses.push(["DIRECTORY", directory.address]);
    console.log("");
    console.log("🚀  Deployed on:", deployer.network)
    console.log("-----------AddressesStart-----------");
    for (c = 0; c < names.length; c++) {
      console.log(names[c] + "_CONTRACT_ADDRESS=" + addresses[c]);
    }
    console.log("-----------AddressesEnd-----------");
  });
};

