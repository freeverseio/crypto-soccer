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
    const account0Owners = deployUtils.getAccount0Owner(accounts[0]);
    const versionNumber = 0;
    const proxyAddress  = "0x0";
    const {0: proxy, 1: assets, 2: market, 3: updates, 4: challenges} = 
      await deployUtils.deploy(versionNumber, account0Owners, Proxy, proxyAddress, Assets, Market, Updates, Challenges).should.be.fulfilled;

    const stakers  = await deployer.deploy(Stakers, requiredStake).should.be.fulfilled;
    const enginePreComp = await deployer.deploy(EnginePreComp).should.be.fulfilled;
    const engineApplyBoosters = await deployer.deploy(EngineApplyBoosters).should.be.fulfilled;
    const engine = await deployer.deploy(Engine, enginePreComp.address, engineApplyBoosters.address).should.be.fulfilled;
    const trainingPoints= await deployer.deploy(TrainingPoints, assets.address).should.be.fulfilled;
    const evolution= await deployer.deploy(Evolution).should.be.fulfilled;
    const leagues = await deployer.deploy(Leagues, assets.address).should.be.fulfilled;
    const friendlies = await deployer.deploy(Friendlies).should.be.fulfilled;
    const shop = await deployer.deploy(Shop).should.be.fulfilled;
    const privileged = await deployer.deploy(Privileged).should.be.fulfilled;
    const utils = await deployer.deploy(Utils).should.be.fulfilled;
    const playAndEvolve = await deployer.deploy(PlayAndEvolve, trainingPoints.address, evolution.address, engine.address, shop.address).should.be.fulfilled;
    const merkle = await deployer.deploy(Merkle).should.be.fulfilled;
    const constantsGetters = await deployer.deploy(ConstantsGetters).should.be.fulfilled;
    const directory = await deployer.deploy(Directory).should.be.fulfilled;
    const marketCrypto = await deployer.deploy(MarketCrypto).should.be.fulfilled;

    console.log("Setting up ...");

    if (versionNumber == 0) { 
      // first set all owners to accounts[0] so that we can do some operations
      await assets.setCOO(accounts[0]).should.be.fulfilled;
      await assets.setMarket(accounts[0]).should.be.fulfilled;
      await updates.setRelay(accounts[0]).should.be.fulfilled;
      await stakers.setCOO(accounts[0]).should.be.fulfilled;

      // do these operations:
      await market.setCryptoMarketAddress(marketCrypto.address).should.be.fulfilled;
      await market.proposeNewMaxSumSkillsBuyNowPlayer(sumSkillsAllowed = 20000, newLapseTime = 5*24*3600).should.be.fulfilled;
      await market.updateNewMaxSumSkillsBuyNowPlayer().should.be.fulfilled;
      await updates.initUpdates().should.be.fulfilled; 
      await updates.setStakersAddress(stakers.address).should.be.fulfilled;
      await stakers.setGameOwner(updates.address).should.be.fulfilled;

      if (singleTimezone != -1) {
        console.log("Init single timezone", singleTimezone);
        await assets.initSingleTZ(singleTimezone).should.be.fulfilled;
      } else {
        await assets.init().should.be.fulfilled;
      }

      // Prepare the final ownerships
      await marketCrypto.setCOO(owners.COO).should.be.fulfilled;
      await stakers.setCOO(owners.COO).should.be.fulfilled;
      await assets.setCOO(owners.COO).should.be.fulfilled;
      await assets.setMarket(owners.market).should.be.fulfilled;
      await updates.setRelay(owners.relay).should.be.fulfilled;
      await proxy.setSuperUser(owners.superuser).should.be.fulfilled;

      await marketCrypto.proposeOwner(owners.superuser).should.be.fulfilled;
      await stakers.proposeOwner(owners.superuser).should.be.fulfilled;
      await proxy.proposeCompany(owners.company).should.be.fulfilled;

      // Execute the final ownerships (WARNING: needs privKeys)
      await marketCrypto.acceptOwner({from: owners.superuser}).should.be.fulfilled;
      await stakers.acceptOwner({from: owners.superuser}).should.be.fulfilled;
      await proxy.acceptCompany({from: owners.company}).should.be.fulfilled;
 

      // If we want stakers signed up during deploy, uncomment this:
      // await deployUtils.addTrustedParties(stakers, owners.COO, owners.trustedParties);
      // await deployUtils.enroll(stakers, requiredStake, owners.trustedParties);
    }


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

