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
const Merkle = artifacts.require('Merkle');
const Challenges = artifacts.require('Challenges');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Directory = artifacts.require('Directory');
const MarketCrypto = artifacts.require('MarketCrypto');
const Stakers = artifacts.require('Stakers');

require('chai')
    .use(require('chai-as-promised'))
    .should();
const assert = require('assert');
const deployUtils = require('../utils/deployUtils.js');



module.exports = function (deployer, network, accounts) {
  deployer.then(async () => {
    const { singleTimezone, owners, requiredStake } = deployUtils.getExplicitOrDefaultSetup(deployer.networks[network], accounts);

    const versionNumber = 0;
    const proxyAddress  = "0x0";
    const {0: proxy, 1: assets, 2: market, 3: updates, 4: challenges} = 
      await deployUtils.deploy(versionNumber, owners, Proxy, proxyAddress, Assets, Market, Updates, Challenges);
  
    const stakers  = await deployer.deploy(Stakers, requiredStake, {from: owners.superuser}).should.be.fulfilled;
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
    const merkle = await deployer.deploy(Merkle).should.be.fulfilled;
    const constantsGetters = await deployer.deploy(ConstantsGetters).should.be.fulfilled;
    const directory = await deployer.deploy(Directory).should.be.fulfilled;
    const marketCrypto = await deployer.deploy(MarketCrypto, {from: owners.superuser}).should.be.fulfilled;

    console.log("Setting up ...");

    if (versionNumber == 0) { 
      await deployUtils.setContractOwners(assets, updates, owners);
      // await assets.setMarket("0x7c34471e39c4A4De223c05DF452e28F0c4BD9BF0", {from: owners.superuser});
      await market.proposeNewMaxSumSkillsBuyNowPlayer(sumSkillsAllowed = 20000, newLapseTime = 5*24*3600, {from: owners.COO}).should.be.fulfilled;
      await market.updateNewMaxSumSkillsBuyNowPlayer({from: owners.COO}).should.be.fulfilled;
      await updates.initUpdates({from: owners.COO}).should.be.fulfilled; 
      await updates.setStakersAddress(stakers.address, {from: owners.superuser}).should.be.fulfilled;
      await stakers.setGameOwner(updates.address, {from: owners.superuser}).should.be.fulfilled;
      await deployUtils.addTrustedParties(stakers, owners.superuser, owners.trustedParties);
      await deployUtils.enroll(stakers, requiredStake, owners.trustedParties);
      if (singleTimezone != -1) {
        console.log("Init single timezone", singleTimezone);
        await assets.initSingleTZ(singleTimezone, {from: owners.COO}).should.be.fulfilled;
      } else {
        await assets.init({from: owners.COO}).should.be.fulfilled;
      }
    }

    await market.setCryptoMarketAddress(marketCrypto.address, {from: owners.COO}).should.be.fulfilled;
    await leagues.setEngineAdress(engine.address).should.be.fulfilled;
    await leagues.setAssetsAdress(assets.address).should.be.fulfilled;
    await trainingPoints.setAssetsAddress(assets.address).should.be.fulfilled;
    await trainingPoints.setMarketAddress(market.address).should.be.fulfilled;
    await engine.setPreCompAddr(enginePreComp.address).should.be.fulfilled;
    await engine.setApplyBoostersAddr(engineApplyBoosters.address).should.be.fulfilled;
    await playAndEvolve.setTrainingAddress(trainingPoints.address).should.be.fulfilled;
    await playAndEvolve.setEvolutionAddress(evolution.address).should.be.fulfilled;
    await playAndEvolve.setEngineAddress(engine.address).should.be.fulfilled;
    await playAndEvolve.setShopAddress(shop.address).should.be.fulfilled;
    await marketCrypto.setCOO(owners.COO, {from: owners.superuser}).should.be.fulfilled;
    await marketCrypto.setMarketFiatAddress(proxy.address, {from: owners.COO}).should.be.fulfilled;

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
      ["FRIENDLIES", friendlies.address],
      ["SHOP_CONTRACT", shop.address],
      ["PRIVILEGED", privileged.address],
      ["UTILS", utils.address],
      ["PLAYANDEVOLVE", playAndEvolve.address],
      ["MERKLE", merkle.address],
      ["CONSTANTSGETTERS", constantsGetters.address],
      ["PROXY", proxy.address],
      ["CHALLENGES", challenges.address],
      ["MARKETCRYPTO", marketCrypto.address],
      ["STAKERS", stakers.address]
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

