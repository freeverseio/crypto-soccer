const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const debug = require('../utils/debugUtils.js');
const timeTravel = require('../utils/TimeTravel.js');
const marketUtils = require('../utils/marketUtils.js');
const delegateUtils = require('../utils/delegateCallUtils.js');

const ConstantsGetters = artifacts.require('ConstantsGetters');
const Proxy = artifacts.require('Proxy');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const Updates = artifacts.require('Updates');
const MarketCrypto = artifacts.require('MarketCrypto');
const Privileged = artifacts.require('Privileged');

async function createPromoPlayer(targetTeamId, internalId = 144321433) {
  sk = [16383, 13, 4, 56, 456];
  traits = [potential = 5, forwardness = 3, leftishness = 4, aggressiveness = 1];
  secsInYear = 365*24*3600;
  playerId = await privileged.createPromoPlayer(
    sk,
    age = 24 * secsInYear,
    traits,
    internalId,
    targetTeamId
  ).should.be.fulfilled;
  return playerId;
}

async function createSpecialPlayerId(internalId = 144321433) {
  sk = [16383, 13, 4, 56, 456];
  traits = [potential = 5, forwardness = 3, leftishness = 4, aggressiveness = 1]
  secsInYear = 365*24*3600
  playerId = await privileged.createSpecialPlayer(
    sk,
    age = 24 * secsInYear,
    traits,
    internalId
  ).should.be.fulfilled;
  return playerId;
}



async function freezeTeam(currencyId, price, sellerRnd, validUntil, teamId, sellerAccount) {
  // Mobile app does this:
  sigSeller = await marketUtils.signPutAssetForSaleMTx(
    currencyId,
    price,
    sellerRnd,
    validUntil, 
    teamId.toNumber(),
    sellerAccount
  );

  // First of all, Freeverse and Buyer check the signature
  // In this case, using web3:
  recoveredSellerAddr = await web3.eth.accounts.recover(sigSeller);
  recoveredSellerAddr.should.be.equal(sellerAccount.address);

  // The correctness of the seller message can also be checked in the BC:
  const sellerHiddenPrice = marketUtils.concatHash(
    ["uint8", "uint256", "uint256"],
    [currencyId, price, sellerRnd]
  );
  sellerTxMsgBC = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, sellerTeamId.toNumber()).should.be.fulfilled;
  sellerTxMsgBC.should.be.equal(sigSeller.message);

  // Then, the buyer builds a message to sign
  let isTeamFrozen = await market.isTeamFrozen(teamId.toNumber()).should.be.fulfilled;
  isTeamFrozen.should.be.equal(false);

  // and send the Freeze TX. 
  const sigSellerRS = [
    sigSeller.r,
    sigSeller.s
  ];
  
  tx = await market.freezeTeam(
    sellerHiddenPrice,
    validUntil,
    teamId.toNumber(),
    sigSellerRS,
    sigSeller.v
  ).should.be.fulfilled;
  
  return tx;
};


async function freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount) {
    // Mobile app does this:
    sigSeller = await marketUtils.signPutAssetForSaleMTx(
      currencyId,
      price,
      sellerRnd,
      validUntil,
      playerId.toString(),
      sellerAccount
    );
    // First of all, Freeverse and Buyer check the signature
    // In this case, using web3:
    recoveredSellerAddr = await web3.eth.accounts.recover(sigSeller);
    recoveredSellerAddr.should.be.equal(sellerAccount.address);

    // The correctness of the seller message can also be checked in the BC:
    const sellerHiddenPrice = marketUtils.concatHash(
      ["uint8", "uint256", "uint256"],
      [currencyId, price, sellerRnd]
    );
    sellerTxMsgBC = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerId).should.be.fulfilled;
    sellerTxMsgBC.should.be.equal(sigSeller.message);

    // Then, the buyer builds a message to sign
    let isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(false);

    // and send the Freeze TX. 
    const sigSellerRS = [
      sigSeller.r,
      sigSeller.s
    ];
    tx = await market.freezePlayer(
      sellerHiddenPrice,
      validUntil,
      playerId,
      sigSellerRS,
      sigSeller.v
    ).should.be.fulfilled;
    return tx;
}

async function completeTeamAuction(
  currencyId, price, sellerRnd, validUntil, sellerTeamId, 
  extraPrice, buyerRnd, isOffer2StartAuctionSig, isOffer2StartAuctionBC, buyerAccount) 
{
  // Add some amount to the price where seller started, and a rnd to obfuscate it
  const buyerHiddenPrice = marketUtils.concatHash(
    ["uint256", "uint256"],
    [extraPrice, buyerRnd]
  );
  let sigBuyer = await marketUtils.signAgreeToBuyTeamMTx(
    currencyId,
    price,
    extraPrice,
    sellerRnd,
    buyerRnd,
    validUntil,
    sellerTeamId.toNumber(),
    isOffer2StartAuctionSig,
    buyerAccount
  ).should.be.fulfilled;

  // Freeverse checks the signature
  recoveredBuyerAddr = await web3.eth.accounts.recover(sigBuyer);
  recoveredBuyerAddr.should.be.equal(buyerAccount.address);

  // Freeverse waits until actual money has been transferred between users, and completes sale
  const sigBuyerRS = [
    sigBuyer.r,
    sigBuyer.s
  ];

  // The correctness of the seller message can also be checked in the BC:
  const sellerHiddenPrice = marketUtils.concatHash(
    ["uint8", "uint256", "uint256"],
    [currencyId, price, sellerRnd]
  );
  
  tx = await market.completeTeamAuction(
    sellerHiddenPrice,
    validUntil,
    sellerTeamId.toNumber(),
    buyerHiddenPrice,
    sigBuyerRS,
    sigBuyer.v,
    recoveredBuyerAddr,
    isOffer2StartAuctionBC
  ).should.be.fulfilled;
  return tx;
}

async function completePlayerAuction(
    currencyId, price, sellerRnd, validUntil, playerId, 
    extraPrice, buyerRnd, isOffer2StartAuctionSig, isOffer2StartAuctionBC, buyerTeamId, buyerAccount
  ) {
  // Add some amount to the price where seller started, and a rnd to obfuscate it
  const buyerHiddenPrice = marketUtils.concatHash(
    ["uint256", "uint256"],
    [extraPrice, buyerRnd]
  );

  let sigBuyer = await marketUtils.signAgreeToBuyPlayerMTx(
    currencyId,
    price,
    extraPrice,
    sellerRnd,
    buyerRnd,
    validUntil,
    playerId.toString(),
    isOffer2StartAuctionSig,
    buyerTeamId.toString(),
    buyerAccount
  ).should.be.fulfilled;

  // Freeverse checks the signature
  recoveredBuyerAddr = await web3.eth.accounts.recover(sigBuyer);
  recoveredBuyerAddr.should.be.equal(buyerAccount.address);

  // Freeverse waits until actual money has been transferred between users, and completes sale
  const sigBuyerRS = [
    sigBuyer.r,
    sigBuyer.s
  ];
  
  const sellerHiddenPrice = marketUtils.concatHash(
    ["uint8", "uint256", "uint256"],
    [currencyId, price, sellerRnd]
  );

  tx = await market.completePlayerAuction(
    sellerHiddenPrice,
    validUntil,
    playerId.toString(),
    buyerHiddenPrice,
    buyerTeamId.toString(),
    sigBuyerRS,
    sigBuyer.v,
    isOffer2StartAuctionBC
  ).should.be.fulfilled;

  return tx
}

async function provideFundsToAdresses(addresses, originAccounts, value = "1000000000000000000") {
  for (i = 0; i < addresses.length; i++) {
    to = addresses[i];
    await web3.eth.sendTransaction({ from: originAccounts[i], to, value }).should.be.fulfilled;
  }
}

async function transferPlayerViaAuction(market, playerId, buyerTeamId, sellerAccount, buyerAccount) {
  currencyId = 1;
  price = 41234;
  rnd = 42321;  
  sellerRnd = 42321;
  extraPrice = 332;
  buyerRnd = 1243523;
  
  now = await market.getBlockchainNowTime().should.be.fulfilled;
  AUCTION_TIME = await constants.get_AUCTION_TIME().should.be.fulfilled;
  validUntil = now.toNumber() + AUCTION_TIME.toNumber();
  tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
  isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(true);
  truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
    return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
  });
  tx = await completePlayerAuction(
    currencyId, price,  sellerRnd, validUntil, playerId, 
    extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
  ).should.be.fulfilled;
  let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
  finalOwner.should.be.equal(buyerAccount.address);
}

contract("Market", accounts => {
  let ok;
  const NULL_TEAMID = 0;
  const NULL_PLAYERID = 0;
  const ACADEMY_TEAM_ID = 1;
  const NULL_ADDR = '0x0000000000000000000000000000000000000000';
  
  const it2 = async(text, f) => {};
  
  beforeEach(async () => {
    depl =  await delegateUtils.deploy(versionNumber = 0, Proxy, '0x0', Assets, Market, Updates);
    proxy  = depl[0];
    assets = depl[1];
    market = depl[2];
    // done with delegate calls
    
    constants = await ConstantsGetters.new().should.be.fulfilled;
    marketCrypto = await MarketCrypto.new().should.be.fulfilled;

    freeverseAccount = await web3.eth.accounts.create("iamFreeverse");
    await assets.init().should.be.fulfilled;
    privileged = await Privileged.new().should.be.fulfilled;
    sellerAccount = await web3.eth.accounts.create("iamaseller");
    buyerAccount = await web3.eth.accounts.create("iamabuyer");
    playerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry = 4);
    sellerTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry1 = 0);
    buyerTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry2 = 1);
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, sellerAccount.address).should.be.fulfilled;
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, buyerAccount.address).should.be.fulfilled;
    now = await market.getBlockchainNowTime().should.be.fulfilled;

    AUCTION_TIME = await constants.get_AUCTION_TIME().should.be.fulfilled;
    AUCTION_TIME = AUCTION_TIME.toNumber();
    
    POST_AUCTION_TIME = await constants.get_POST_AUCTION_TIME().should.be.fulfilled;
    POST_AUCTION_TIME = POST_AUCTION_TIME.toNumber();
    
    validUntil = now.toNumber() + AUCTION_TIME;
    currencyId = 1;
    price = 41234;
    sellerRnd = 42321;
    extraPrice = 332;
    buyerRnd = 1243523;

  });
  
  it2("normal players, go above 25, and get rid of player", async () => {
    playerIds = [];
    nPlayersToBuy = 9;
    for (i = 0; i < nPlayersToBuy; i++) {
      playerIds.push(playerId.add(web3.utils.toBN(i))); 
      tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerIds[i], sellerAccount).should.be.fulfilled;
    }
    for (i = 0; i < nPlayersToBuy; i++) {
      tx = await completePlayerAuction(
        currencyId, price,  sellerRnd, validUntil, playerIds[i], 
        extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
      ).should.be.fulfilled;
    }

    nTransit = await market.getNPlayersInTransitInTeam(buyerTeamId).should.be.fulfilled;
    nTransit.toNumber().should.be.equal(2);

    // transfer fails because team is still full
    await market.completePlayerTransit(playerIds[nPlayersToBuy-1]).should.be.rejected;
    await market.completePlayerTransit(playerIds[nPlayersToBuy-2]).should.be.rejected;
    
  });

  it2("crypto flow with player" , async () => {
    // set up teams: team 2 - ALICE, team 3 - BOB, team 4 - CAROL
    ALICE = accounts[0];
    BOB = accounts[1];
    CAROL = accounts[2];
    await marketCrypto.setMarketAddress(proxy.address).should.be.fulfilled;
    startingPrice = web3.utils.toWei('1');
    teamIdxInCountry0 = 2; 
    playerId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry0 = teamIdxInCountry0*18+3);
    sellerTeamId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0);
    tx = await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, ALICE).should.be.fulfilled;

    // ALICE puts for sale
    tx = await marketCrypto.putPlayerForSale(playerId0, startingPrice, {from: ALICE}).should.be.fulfilled;
    truffleAssert.eventEmitted(tx, "PlayerPutForSaleCrypto", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId0) && event.startingPrice.should.be.bignumber.equal(web3.utils.toBN(startingPrice));
    });

    auctionId = await marketCrypto.getCurrentAuctionForPlayer(playerId0).should.be.fulfilled;
    auctionId.toNumber().should.be.equal(1); 
    
    // TEST all auction getters at this stage:
    var {0: validUn, 1: teamIdHi, 2: hiBid, 3: hiBidder, 4: sell, 5: assetWent} = await marketCrypto.getAuctionData(auctionId).should.be.fulfilled;
    teamIdHi.toNumber().should.be.equal(0);
    hiBid.toNumber().should.be.equal(0);
    hiBidder.should.be.equal('0x0000000000000000000000000000000000000000');
    sell.should.be.equal(ALICE);
    assetWent.should.be.equal(false);
    (Math.abs(validUn - now.toNumber()) < 24*3600 + 10).should.be.equal(true);
    (Math.abs(validUn - now.toNumber()) > 24*3600 - 10).should.be.equal(true);

    
    // BOB does first bid
    tx = await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, BOB).should.be.fulfilled;
    buyerTeamId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0 + 1);

    tx = await marketCrypto.bidForPlayer(playerId0, buyerTeamId0, {from: BOB, value: startingPrice}).should.be.fulfilled;

    past = await market.getPastEvents( 'PlayerFreezeCrypto', { fromBlock: 0, toBlock: 'latest' } ).should.be.fulfilled;
    past[0].args.playerId.should.be.bignumber.equal(playerId0);
    past[0].args.frozen.should.be.equal(true);

    truffleAssert.eventEmitted(tx, "BidForPlayerCrypto", (event) => {
      return  event.playerId.should.be.bignumber.equal(playerId0) && 
              event.bidderTeamId.should.be.bignumber.equal(buyerTeamId0) && 
              event.totalAmount.should.be.bignumber.equal(web3.utils.toBN(startingPrice));
    });

    // TEST all auction getters at this stage:
    var {0: validUn, 1: teamIdHi, 2: hiBid, 3: hiBidder, 4: sell, 5: assetWent} = await marketCrypto.getAuctionData(auctionId).should.be.fulfilled;
    teamIdHi.should.be.bignumber.equal(buyerTeamId0);
    hiBid.should.be.bignumber.equal(web3.utils.toBN(startingPrice));
    hiBidder.should.be.equal(BOB);
    sell.should.be.equal(ALICE);
    assetWent.should.be.equal(false);
    (Math.abs(validUn - now.toNumber()) < 24*3600 + 10).should.be.equal(true);
    (Math.abs(validUn - now.toNumber()) > 24*3600 - 10).should.be.equal(true);
    
 
    tx = await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, CAROL).should.be.fulfilled;
    buyerTeamId1 = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0 + 2);

    await marketCrypto.bidForPlayer(playerId0, buyerTeamId1, {from: CAROL, value: startingPrice}).should.be.rejected;

    newBid = web3.utils.toWei('1.1');
    await marketCrypto.bidForPlayer(playerId0, buyerTeamId1, {from: CAROL, value: (newBid)}).should.be.rejected;

    newBid = web3.utils.toWei('1.5');
    // ALICE cannot bid for her own player:
    tx = await marketCrypto.bidForPlayer(playerId0, buyerTeamId1, {from: ALICE, value: (newBid)}).should.be.rejected;
    // but CAROL can:
    tx = await marketCrypto.bidForPlayer(playerId0, buyerTeamId1, {from: CAROL, value: (newBid)}).should.be.fulfilled;

    truffleAssert.eventEmitted(tx, "BidForPlayerCrypto", (event) => {
      return  event.playerId.should.be.bignumber.equal(playerId0) && 
              event.bidderTeamId.should.be.bignumber.equal(buyerTeamId1) && 
              event.totalAmount.should.be.bignumber.equal(web3.utils.toBN(newBid));
    });

    // TEST all auction getters at this stage:
    var {0: validUn, 1: teamIdHi, 2: hiBid, 3: hiBidder, 4: sell, 5: assetWent} = await marketCrypto.getAuctionData(auctionId).should.be.fulfilled;
    teamIdHi.should.be.bignumber.equal(buyerTeamId1);
    hiBid.should.be.bignumber.equal(web3.utils.toBN(newBid));
    hiBidder.should.be.equal(CAROL);
    sell.should.be.equal(ALICE);
    assetWent.should.be.equal(false);
    (Math.abs(validUn - now.toNumber()) < 24*3600 + 10).should.be.equal(true);
    (Math.abs(validUn - now.toNumber()) > 24*3600 - 10).should.be.equal(true);
    
     
    await marketCrypto.withdraw(auctionId, {from: CAROL}).should.be.rejected;

    
    balanceBefore = await web3.eth.getBalance(BOB);
    await marketCrypto.withdraw(auctionId, {from: BOB}).should.be.fulfilled;
    balanceAfter = await web3.eth.getBalance(BOB);
    // checks that BOB has as much as he had at the beginning up to gas costs
    (startingPrice - (balanceAfter - balanceBefore) < web3.utils.toWei('0.001')).should.be.equal(true);

    await marketCrypto.withdraw(auctionId, {from: CAROL}).should.be.rejected;
    await marketCrypto.withdraw(auctionId, {from: ALICE}).should.be.rejected;
    await marketCrypto.executePlayerTransfer(playerId0).should.be.rejected;

    await timeTravel.advanceTime(24*3600-100);
    await timeTravel.advanceBlock().should.be.fulfilled;
    await marketCrypto.withdraw(auctionId, {from: ALICE}).should.be.rejected;
    await marketCrypto.executePlayerTransfer(playerId0).should.be.rejected;

    await timeTravel.advanceTime(0.1*3600);
    await timeTravel.advanceBlock().should.be.fulfilled;
    await marketCrypto.withdraw(auctionId, {from: ALICE}).should.be.fulfilled;
    tx = await marketCrypto.executePlayerTransfer(playerId0).should.be.fulfilled;
    truffleAssert.eventEmitted(tx, "AssetWentToNewOwner", (event) => {
      return  event.playerId.should.be.bignumber.equal(playerId0) && 
              event.auctionId.should.be.bignumber.equal(web3.utils.toBN(1)) 
    });

    past = await market.getPastEvents( 'PlayerFreezeCrypto', { fromBlock: 0, toBlock: 'latest' } ).should.be.fulfilled;
    past[1].args.playerId.should.be.bignumber.equal(playerId0);
    past[1].args.frozen.should.be.equal(false);

    finalOwner = await market.getOwnerPlayer(playerId0).should.be.fulfilled;
    finalOwner.should.be.equal(CAROL);
    
    // TEST all auction getters at this stage:
    var {0: validUn, 1: teamIdHi, 2: hiBid, 3: hiBidder, 4: sell, 5: assetWent} = await marketCrypto.getAuctionData(auctionId).should.be.fulfilled;
    teamIdHi.should.be.bignumber.equal(buyerTeamId1);
    hiBid.should.be.bignumber.equal(web3.utils.toBN(0));
    hiBidder.should.be.equal(CAROL);
    sell.should.be.equal(ALICE);
    assetWent.should.be.equal(true);
    (Math.abs(validUn - now.toNumber()) < 24*3600 + 10).should.be.equal(true);
    (Math.abs(validUn - now.toNumber()) > 24*3600 - 10).should.be.equal(true);
    
  });

  it2("crypto mkt shows that we can get past 25 players" , async () => {
    // set up teams: team 2 - ALICE, team 3 - BOB
    ALICE = accounts[0];
    BOB = accounts[1];
    await marketCrypto.setMarketAddress(proxy.address).should.be.fulfilled;
    startingPrice = web3.utils.toWei('1');
    teamIdxInCountry0 = 2; 

    // ALICE will be selling
    sellerTeamId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0);
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, ALICE).should.be.fulfilled;
    // BOB will be buying
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, BOB).should.be.fulfilled;
    buyerTeamId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0 + 1);

    nTransfers = 10;
    playerIds = [];
    for (n = 0; n < nTransfers; n++) { 
      thisPlayerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry0 = teamIdxInCountry0*18+2+n);
      playerIds.push(thisPlayerId);
      // ALICE puts for sale and BOB bids
      await marketCrypto.putPlayerForSale(thisPlayerId, startingPrice, {from: ALICE}).should.be.fulfilled;

      result = await market.isPlayerFrozenInAnyMarket(thisPlayerId).should.be.fulfilled;
      result.should.be.equal(false);

      await marketCrypto.bidForPlayer(thisPlayerId, buyerTeamId0, {from: BOB, value: startingPrice}).should.be.fulfilled;
      result = await market.isPlayerFrozenInAnyMarket(thisPlayerId).should.be.fulfilled;
      result.should.be.equal(true);
      }

    await timeTravel.advanceTime(24*3600+100);
    await timeTravel.advanceBlock().should.be.fulfilled;

    // show that the first 7 transfers works normally, for the others, players go to limbo
    IN_TRANSIT_TEAM = 2;
    for (n = 0; n < nTransfers; n++) { 
      await marketCrypto.executePlayerTransfer(playerIds[n]).should.be.fulfilled;
      if (n >= 7) {
        ownTeam = await market.getCurrentTeamIdFromPlayerId(playerIds[n]).should.be.fulfilled;
        ownTeam.toNumber().should.be.equal(IN_TRANSIT_TEAM);
      }
    }
    // note that the players are not frozen anymore. However, it'll be impossible to freeze them since
    // they currently belong to IN_TRANSIT_TEAM
    result = await market.isPlayerFrozenInAnyMarket(playerIds[8]).should.be.fulfilled;
    result.should.be.equal(false);
  
    // buyer can not continue bidding because he already has players in limbo:
    thisPlayerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry0 = teamIdxInCountry0*18+2+nTransfers+1);
    await marketCrypto.putPlayerForSale(thisPlayerId, startingPrice, {from: ALICE}).should.be.fulfilled;
    await marketCrypto.bidForPlayer(thisPlayerId, buyerTeamId0, {from: BOB, value: startingPrice}).should.be.rejected;

    // so player 7, 8, 9 are in transit. We can still not complete transits because team is still full
    await market.completePlayerTransit(playerIds[7]).should.be.rejected;

    // so let's sell again, return back to ALICE :-)
    for (n = 0; n < 3; n++) { 
      await marketCrypto.putPlayerForSale(playerIds[n], startingPrice, {from: BOB}).should.be.fulfilled;
      await marketCrypto.bidForPlayer(playerIds[n], sellerTeamId0, {from: ALICE, value: startingPrice}).should.be.fulfilled;
    }
    await timeTravel.advanceTime(24*3600+100);
    await timeTravel.advanceBlock().should.be.fulfilled;
    for (n = 0; n < 3; n++) { 
      await marketCrypto.executePlayerTransfer(playerIds[n]).should.be.fulfilled;
    }
    await market.completePlayerTransit(playerIds[7]).should.be.fulfilled;
    await market.completePlayerTransit(playerIds[8]).should.be.fulfilled;
    await market.completePlayerTransit(playerIds[9]).should.be.fulfilled;
    ownTeam = await market.getCurrentTeamIdFromPlayerId(playerIds[8]).should.be.fulfilled;
    ownTeam.should.be.bignumber.equal(buyerTeamId0);

    result = await market.isPlayerFrozenInAnyMarket(playerIds[8]).should.be.fulfilled;
    result.should.be.equal(false);
  });

  it2("crypto mkt cannot bid if buyerTeamId = sellerTeamId" , async () => {
    // set up teams: ALICE, BOB, ALICE
    ALICE = accounts[0];
    BOB = accounts[1];

    await marketCrypto.setMarketAddress(proxy.address).should.be.fulfilled;
    startingPrice = web3.utils.toWei('1');
    teamIdxInCountry0 = 2; 

    // ALICE will be 
    sellerTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0);
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, ALICE).should.be.fulfilled;
    // BOB will be buying
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, BOB).should.be.fulfilled;
    buyerTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0 + 1);
    // ALICE will be buying too
    await assets.transferFirstBotToAddr(tz = 1, countryIdxInTZ = 0, BOB).should.be.fulfilled;
    allice2ndTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry0 + 2);

    
    // ALICE puts for sale and BOB bids, all OK
    thisPlayerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry0 = teamIdxInCountry0 * 18 + 1);
    await marketCrypto.putPlayerForSale(thisPlayerId, startingPrice, {from: ALICE}).should.be.fulfilled;
    // await web3.eth.sendTransaction({ from: originAccounts[i], to, value }).should.be.fulfilled;

    await marketCrypto.bidForPlayer(thisPlayerId, buyerTeamId, {from: BOB, value: startingPrice}).should.be.fulfilled;
    result = await market.isPlayerFrozenInAnyMarket(thisPlayerId).should.be.fulfilled;
    result.should.be.equal(true);

    // ALICE puts for sale and ALICE bids, fails because of same teamId and same owner
    thisPlayerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry0 = teamIdxInCountry0 * 18 + 2);
    await marketCrypto.putPlayerForSale(thisPlayerId, startingPrice, {from: ALICE}).should.be.fulfilled;
    await marketCrypto.bidForPlayer(thisPlayerId, sellerTeamId, {from: ALICE, value: startingPrice}).should.be.rejected;

    // ALICE puts for sale and ALICE bids with her 2nd team, fails because of same owner
    thisPlayerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry0 = teamIdxInCountry0 * 18 + 3);
    await marketCrypto.putPlayerForSale(thisPlayerId, startingPrice, {from: ALICE}).should.be.fulfilled;
    await marketCrypto.bidForPlayer(thisPlayerId, allice2ndTeamId, {from: ALICE, value: startingPrice}).should.be.rejected;

    // TODO: test when BOB aquires ALICE's team, and then tries to bid. He'll be a different ownerAddr, but same teamID. Should fail
    // problem is I didn't figure out how to have an account that both i) signs as needed for our MarketAPI, ii) signs TXs in Truffle
  });

  
  it2('setAcquisitionConstraint of constraints in friendliess', async () => {
    maxNumConstraints = 7;
    remainingAcqs = 0;
    for (c = 0; c < maxNumConstraints; c++) {
      acq = c;
      isIt = await market.isAcquisitionFree(remainingAcqs, acq).should.be.fulfilled;
      isIt.should.be.equal(true);
      remainingAcqs = await market.setAcquisitionConstraint(remainingAcqs, valUnt = now.toNumber() + c * 4400, numRemain = c, acq).should.be.fulfilled;
      isIt = await market.isAcquisitionFree(remainingAcqs, acq).should.be.fulfilled;
      isIt.should.be.equal(false);
      valid = await market.getAcquisitionConstraintValidUntil(remainingAcqs, acq = c).should.be.fulfilled;
      num = await market.getAcquisitionConstraintNRemain(remainingAcqs, acq = c).should.be.fulfilled;
      valid.toNumber().should.be.equal(valUnt);
      num.toNumber().should.be.equal(numRemain);
    }
  });
  
  it2('addAcquisitionConstraint of constraints in friendlies', async () => {
    maxNumConstraints = 6;
    teamId = buyerTeamId;
    for (c = 0; c < maxNumConstraints; c++) {
      await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + (c + 1) * 4400, numRemain = c + 1).should.be.fulfilled;
    }
    // the team is full already
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + (c + 1) * 4400, numRemain = c + 1).should.be.rejected;
    // as just enough time passes it can do one more submission again:
    await timeTravel.advanceTime(4400 + 1000);
    await timeTravel.advanceBlock().should.be.fulfilled;
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + (c + 1) * 4400, numRemain = c + 1).should.be.fulfilled;
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + (c + 1) * 4400, numRemain = c + 1).should.be.rejected;
  });

  it2('encoding of constraints pass with time', async () => {
    teamId = buyerTeamId;
    remainingAcqs = 0;
    acq = 5;
    isIt = await market.isAcquisitionFree(remainingAcqs, acq).should.be.fulfilled;
    isIt.should.be.equal(true);
    remainingAcqs = await market.setAcquisitionConstraint(remainingAcqs, valUnt = now.toNumber() - 10, numRemain = c, acq).should.be.fulfilled;
    isIt = await market.isAcquisitionFree(remainingAcqs, acq).should.be.fulfilled;
    isIt.should.be.equal(true);
  });

  it2('getMaxAllowedAcquisitions and decreaseMaxAllowedAcquisitions', async () => {
    teamId = buyerTeamId;
    // initially, isContrained = false
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(false);
    nRemain.toNumber().should.be.equal(0);
    // we add 1 contraint
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + 4400, numRemain = 8).should.be.fulfilled;
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(true);
    nRemain.toNumber().should.be.equal(numRemain);
    // we another constraint, but in the past, so nothing changes
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() - 4400, n = 7).should.be.fulfilled;
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(true);
    nRemain.toNumber().should.be.equal(numRemain);
    // we another constraint, it takes into account the lowest constaint (in this case, the newest)
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + 6666, n = 7).should.be.fulfilled;
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(true);
    nRemain.toNumber().should.be.equal(n);
    // we another constraint, it takes into account the lowest constaint (in this case, the previous one)
    await market.addAcquisitionConstraint(teamId, valUnt = now.toNumber() + 6666, n2 = 15).should.be.fulfilled;
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(true);
    nRemain.toNumber().should.be.equal(n);
    // decreaseMaxAllowedAcquisitions twice
    await market.decreaseMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    await market.decreaseMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(true);
    nRemain.toNumber().should.be.equal(n-2);
    // after a long time, it's ready again
    await timeTravel.advanceTime(6666 + 1000);
    await timeTravel.advanceBlock().should.be.fulfilled;
    result = await  market.getMaxAllowedAcquisitions(teamId).should.be.fulfilled;
    var {0: isConstrained, 1: nRemain} = result;
    isConstrained.should.be.equal(false);
    nRemain.toNumber().should.be.equal(0);
  });
  

  // *************************************************************************
  // *********************************   TEST  *******************************
  // *************************************************************************
   
  it2('players: deterministic sign (values used in market.notary test)', async () => {
    sellerTeamId.should.be.bignumber.equal('274877906944');
    buyerTeamId.should.be.bignumber.equal('274877906945');
    sellerTeamPlayerIds = await market.getPlayerIdsInTeam(sellerTeamId).should.be.fulfilled;
    const playerIdToSell = sellerTeamPlayerIds[0];
    playerIdToSell.should.be.bignumber.equal('274877906944');

    const sellerAccount = web3.eth.accounts.privateKeyToAccount('0x3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54');
    const buyerAccount = await web3.eth.accounts.privateKeyToAccount('0x3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882');

    // Define params of the seller, and sign
    validUntil = 2000000000;
    buyerHiddenPrice = marketUtils.concatHash(
      ["uint256", "uint256"],
      [extraPrice, buyerRnd]
    );

    buyerHiddenPrice.should.be.equal('0xd46585a1479af8dcc6f2ce8495174282f8c874f1583182bbe2c9df7ae2358edc');

    const sellerHiddenPrice = await market.hashPrivateMsg(currencyId, price, sellerRnd).should.be.fulfilled;
    sellerHiddenPrice.should.be.equal('0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa');
    const message = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerIdToSell).should.be.fulfilled;

    message.should.be.equal('0x909e2fbc45b398649f58c7ea4b632ff1b457ee5f60a43a70abfe00d50e7c917d');
    const sigSeller = sellerAccount.sign(message);
    sigSeller.messageHash.should.be.equal('0x55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a');
    sigSeller.signature.should.be.equal('0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c');

    const prefixed = await market.prefixed(message).should.be.fulfilled;
    prefixed.should.be.equal('0x55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a');
    const isOffer2StartAuction = true;
    const buyerMsg = await market.buildAgreeToBuyPlayerTxMsg(prefixed, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction).should.be.fulfilled;
    buyerMsg.should.be.equal('0xc049e2032b873dd67cc7cc43fb2488f7c770d1654716fc75024cda693c74488c');
    const sigBuyer = buyerAccount.sign(buyerMsg);
    sigBuyer.messageHash.should.be.equal('0xe04d23ec0424b22adec87879118715ce75997a4fd47897c398f3a8cad79b3041');
    sigBuyer.signature.should.be.equal('0xdbe104e7b51c9b1e38cdda4e31c2036e531f7d3338d392bee2f526c4c892437f5e50ddd44224af8b3bd92916b93e4b0d7af2974175010323da7dedea19f30d621c');
  });

  it2('teams: deterministic sign (values used in market.notary test)', async () => {
    sellerTeamId.should.be.bignumber.equal('274877906944');

    const sellerAccount = web3.eth.accounts.privateKeyToAccount('0x3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54');
    const buyerAccount = await web3.eth.accounts.privateKeyToAccount('0x3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882');

    // Define params of the seller, and sign
    validUntil = 2000000000;
    buyerHiddenPrice = marketUtils.concatHash(
      ["uint256", "uint256"],
      [extraPrice, buyerRnd]
    );

    const sellerHiddenPrice = await market.hashPrivateMsg(currencyId, price, sellerRnd).should.be.fulfilled;
    sellerHiddenPrice.should.be.equal('0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa');
    const message = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, sellerTeamId).should.be.fulfilled;

    message.should.be.equal('0x909e2fbc45b398649f58c7ea4b632ff1b457ee5f60a43a70abfe00d50e7c917d');
    const sigSeller = sellerAccount.sign(message);
    sigSeller.messageHash.should.be.equal('0x55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a');
    sigSeller.signature.should.be.equal('0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c');

    const prefixed = await market.prefixed(message).should.be.fulfilled;
    const isOffer2StartAuction = true;
    const buyerMsg = await market.buildAgreeToBuyTeamTxMsg(prefixed, buyerHiddenPrice, isOffer2StartAuction).should.be.fulfilled;
    buyerMsg.should.be.equal('0xdd3d39b424073a7a74a333d3b35bc2b0adea64c4a51c47c4669d190111e7b5e5');
    const sigBuyer = buyerAccount.sign(buyerMsg);
    sigBuyer.messageHash.should.be.equal('0xeb0feff7cbf76cd8f6a6bb07b2d92305e1978c66a157b7738e249689682942f7');
    sigBuyer.signature.should.be.equal('0x7c6b08dfff430bd5dd1785463846f3961f3844b9b4d1cccc941ad2d5441b4496556ffc4518f9be660e2c34ba3d74ef67665af727c25eae6758695354b36462f71b');
  });

  
  // ------------------------------------------------------------------------------------ 
  // ------------------------------------------------------------------------------------ 
  // ------------------------------------------------------------------------------------ 
  // ----------------------------------------------------------------- TEAMS 
  // ------------------------------------------------------------------------------------
  // ------------------------------------------------------------------------------------
  // ------------------------------------------------------------------------------------
  
  it2("teams: completes a MAKE_AN_OFFER via MTXs", async () => {
    // now, sellerRnd is fixed by offerer
    offererRnd = 23987435;
    offerValidUntil = now.toNumber() + 3600; // valid for an hour
    const validUntil = now.toNumber() + 3000 + AUCTION_TIME; // this is, at most, offerValidUntil + AUCTION_TIME
    
    tx = await freezeTeam(currencyId, price, offererRnd, validUntil, sellerTeamId, sellerAccount).should.be.fulfilled;
    isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
      return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completeTeamAuction(
      currencyId, price, offererRnd, offerValidUntil, sellerTeamId, 
      extraPrice = 0, buyerRnd = 0, isOffer2StartAuctionSig = true, isOffer2StartAuctionBC = true, buyerAccount
    ).should.be.fulfilled;
  });
  
  it2("teams: fails a MAKE_AN_OFFER via MTXs because offerValidUntil had expired", async () => {
    // now, sellerRnd is fixed by offerer
    offererRnd = 23987435;
    offerValidUntil = now.toNumber() + 3600; // valid for an hour
    const validUntil = now.toNumber() + 3601 + AUCTION_TIME; // this is, at most, offerValidUntil + AUCTION_TIME

    tx = await freezeTeam(currencyId, price, offererRnd, validUntil, sellerTeamId, sellerAccount).should.be.fulfilled;
    isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
      return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completeTeamAuction(
      currencyId, price, offererRnd, offerValidUntil, sellerTeamId, 
      extraPrice = 0, buyerRnd = 0, isOffer2StartAuctionSig = true, isOffer2StartAuctionBC = true, buyerAccount
    ).should.be.rejected;
  });

  it2("teams: fails a MAKE_AN_OFFER via MTXs because validUntil is too large", async () => {
    validUntil = now.toNumber() + 3600*24*2; // two days

    sigSeller = await marketUtils.signPutAssetForSaleMTx(
      currencyId,
      price,
      offererRnd, // he reuses the rnd provided
      validUntil, 
      sellerTeamId.toNumber(),
      sellerAccount
    );


    // First of all, Freeverse and Buyer check the signature
    // In this case, using web3:
    recoveredSellerAddr = await web3.eth.accounts.recover(sigSeller);
    recoveredSellerAddr.should.be.equal(sellerAccount.address);

    // The correctness of the seller message can also be checked in the BC:
    const sellerHiddenPrice = marketUtils.concatHash(
      ["uint8", "uint256", "uint256"],
      [currencyId, price, offererRnd]
    );
    sellerTxMsgBC = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, sellerTeamId.toNumber()).should.be.fulfilled;
    sellerTxMsgBC.should.be.equal(sigSeller.message);

    // Then, the buyer builds a message to sign
    let isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(false);

    // and send the Freeze TX. 
    const sigSellerRS = [
      sigSeller.r,
      sigSeller.s
    ];

    // we can double-check that it would work
    ok = await market.areFreezeTeamRequirementsOK(
      sellerHiddenPrice,
      validUntil,
      sellerTeamId.toNumber(),
      sigSellerRS,
      sigSeller.v
    ).should.be.fulfilled;
    ok.should.be.equal(false);
    
    // and finally do the freeze 
    tx = await market.freezeTeam(
      sellerHiddenPrice,
      validUntil,
      sellerTeamId.toNumber(),
      sigSellerRS,
      sigSeller.v
    ).should.be.rejected;
  });
  
 
  it2("teams: completes a PUT_FOR_SALE and AGREE_TO_BUY via MTXs", async () => {
    // 1. buyer's mobile app sends to Freeverse: sigBuyer AND params (currencyId, price, ....)
    // 2. Freeverse checks signature and returns to buyer: OK, failed
    // 3. Freeverse advertises to owner that there is an offer to buy his asset at price
    // 4. seller's mobile app sends to Freeverse: sigSeller and params
    // 5. Freeverse checks signature and returns to seller: OK, failed
    // 6. Freeverse FREEZES the player by sending a TX to the BLOCKCHAIN
    // 7. If freeze went OK:
    //          urges buyer to complete payment
    //    If freeze not OK (he probably sold the player in a different market)
    //          tells the buyer to forget about this player
    // 8. Freeverse receives confirmation from Paypal, Apple, GooglePay... of payment buyer -> seller
    // 9. Freeverse COMPLETES TRANSFER OF PLAYER USING BLOCKCHAIN
    tx = await freezeTeam(currencyId, price, sellerRnd, validUntil, sellerTeamId, sellerAccount).should.be.fulfilled;
    isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
      return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completeTeamAuction(
      currencyId, price, sellerRnd, validUntil, sellerTeamId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerAccount
    ).should.be.fulfilled;
    
    truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
      return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(false);
    });

    let finalOwner = await market.getOwnerTeam(sellerTeamId.toNumber()).should.be.fulfilled;
    finalOwner.should.be.equal(buyerAccount.address);
  });

  it2("teams: fails a PUT_FOR_SALE and AGREE_TO_BUY via MTXs because isOffer2StartAuction is not correctly set ", async () => {
    tx, sellerHiddenPrice = await freezeTeam(currencyId, price, sellerRnd, validUntil, sellerTeamId, sellerAccount).should.be.fulfilled;
    isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
      return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completeTeamAuction(
      currencyId, price, sellerRnd, validUntil, sellerTeamId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = true, buyerAccount
    ).should.be.rejected;    
  });

  it2("teams: fails a PUT_FOR_SALE and AGREE_TO_BUY via MTXs because one of its players already frozen", async () => {

    // make sure we'll put for sale a player who belongs to the team that we will also put for sale.
    teamId = await market.getCurrentTeamIdFromPlayerId(playerId).should.be.fulfilled;
    teamId.should.be.bignumber.equal(sellerTeamId);
    
    // put player:
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    
    // fail to put team:
    tx = await freezeTeam(currencyId, price, sellerRnd, validUntil, sellerTeamId, sellerAccount).should.be.rejected;
    isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(false);
    
  });


  // ------------------------------------------------------------------------------------ 
  // ------------------------------------------------------------------------------------ 
  // ------------------------------------------------------------------------------------ 
  // ----------------------------------------------------------------- PLAYERS 
  // ------------------------------------------------------------------------------------
  // ------------------------------------------------------------------------------------
  // ------------------------------------------------------------------------------------

  it2("players: fails a PUT_FOR_SALE and AGREE_TO_BUY via MTXs because his team is already frozen", async () => {

    // make sure we'll put for sale a player who belongs to the team that we will also put for sale.
    teamId = await market.getCurrentTeamIdFromPlayerId(playerId).should.be.fulfilled;
    teamId.should.be.bignumber.equal(sellerTeamId);

    // put team:
    tx = await freezeTeam(currencyId, price, sellerRnd, validUntil, sellerTeamId, sellerAccount).should.be.fulfilled;
    isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
    isTeamFrozen.should.be.equal(true);
    
    // fail to put player:
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.rejected;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(false);
  });
  
  
  it2("players: completes a MAKE_AN_OFFER via MTXs", async () => {
    // now, sellerRnd is fixed by offerer
    offererRnd = 23987435;
    offerValidUntil = now.toNumber() + 3600; // valid for an hour
    const validUntil = now.toNumber() + 3000 + AUCTION_TIME; // this is, at most, offerValidUntil + AUCTION_TIME

    tx = await freezePlayer(currencyId, price, offererRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });

    // the MTX was actually created before the seller put the asset for sale, but it is used now to complete the auction  
    tx = await completePlayerAuction(
      currencyId, price,  offererRnd, offerValidUntil, playerId, 
      extraPrice = 0, buyerRnd = 0, isOffer2StartAuctionSig = true, isOffer2StartAuctionBC = true, buyerTeamId, buyerAccount
    ).should.be.fulfilled;

    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(false);
    });

    let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    finalOwner.should.be.equal(buyerAccount.address);
  });
  
  it2("players: fails a MAKE_AN_OFFER via MTXs because offerValidUntil had expired", async () => {
    // now, sellerRnd is fixed by offerer
    offererRnd = 23987435;
    offerValidUntil = now.toNumber() + 3600; // valid for an hour
    const validUntil = now.toNumber() + 3601 + AUCTION_TIME; // this is, at most, offerValidUntil + AUCTION_TIME

    tx = await freezePlayer(currencyId, price, offererRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });

    // the MTX was actually created before the seller put the asset for sale, but it is used now to complete the auction  
    tx = await completePlayerAuction(
      currencyId, price,  offererRnd, offerValidUntil, playerId, 
      extraPrice = 0, buyerRnd = 0, isOffer2StartAuctionSig = true, isOffer2StartAuctionBC = true, buyerTeamId, buyerAccount
    ).should.be.rejected;
    
  });
  
  it2("players: fails a MAKE_AN_OFFER via MTXs because validUntil is too large", async () => {
    tx, sellerHiddenPrice = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    validUntil = now.toNumber() + 3600*24*2; // two days
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.rejected;
  });
  
  it2("players: fails a PUT_FOR_SALE and AGREE_TO_BUY via MTXs because targetTeam = originTeam", async () => {
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, sellerTeamId, sellerAccount
    ).should.be.rejected;

  });
  
  it2("players: completes a PUT_FOR_SALE and AGREE_TO_BUY via MTXs", async () => {
    // 1. buyer's mobile app sends to Freeverse: sigBuyer AND params (currencyId, price, ....)
    // 2. Freeverse checks signature and returns to buyer: OK, failed
    // 3. Freeverse advertises to owner that there is an offer to buy his asset at price
    // 4. seller's mobile app sends to Freeverse: sigSeller and params
    // 5. Freeverse checks signature and returns to seller: OK, failed
    // 6. Freeverse FREEZES the player by sending a TX to the BLOCKCHAIN
    // 7. If freeze went OK:
    //          urges buyer to complete payment
    //    If freeze not OK (he probably sold the player in a different market)
    //          tells the buyer to forget about this player
    // 8. Freeverse receives confirmation from Paypal, Apple, GooglePay... of payment buyer -> seller
    // 9. Freeverse COMPLETES TRANSFER OF PLAYER USING BLOCKCHAIN

    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
    ).should.be.fulfilled;

    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(false);
    });

    let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    finalOwner.should.be.equal(buyerAccount.address);
  });
  
  it2("players: completes a PUT_FOR_SALE and AGREE_TO_BUY via MTXs", async () => {
    // 1. buyer's mobile app sends to Freeverse: sigBuyer AND params (currencyId, price, ....)
    // 2. Freeverse checks signature and returns to buyer: OK, failed
    // 3. Freeverse advertises to owner that there is an offer to buy his asset at price
    // 4. seller's mobile app sends to Freeverse: sigSeller and params
    // 5. Freeverse checks signature and returns to seller: OK, failed
    // 6. Freeverse FREEZES the player by sending a TX to the BLOCKCHAIN
    // 7. If freeze went OK:
    //          urges buyer to complete payment
    //    If freeze not OK (he probably sold the player in a different market)
    //          tells the buyer to forget about this player
    // 8. Freeverse receives confirmation from Paypal, Apple, GooglePay... of payment buyer -> seller
    // 9. Freeverse COMPLETES TRANSFER OF PLAYER USING BLOCKCHAIN

    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
    ).should.be.fulfilled;

    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(false);
    });

    let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    finalOwner.should.be.equal(buyerAccount.address);
  });
  
  it("players: completes a PUT_FOR_SALE and AGREE_TO_BUY via MTXs - via function call", async () => {
    await transferPlayerViaAuction(market, playerId, buyerTeamId, sellerAccount, buyerAccount).should.be.fulfilled;
  });
  
  
  it2("players: tests constraints on players", async () => {
    await market.addAcquisitionConstraint(buyerTeamId, valUnt = now.toNumber() + 1000, n = 1).should.be.fulfilled;
    // first acquisition works:
    playerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry = 4);
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
    ).should.be.fulfilled;
    // second acquisition should fail:
    playerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry = 5);
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
    ).should.be.rejected;
  });  

  it2("behaviour of getCurrentTeamIdFromPlayerId", async () => {
    playerId = await createSpecialPlayerId();

    teamId = await market.getCurrentTeamIdFromPlayerId(playerId).should.be.fulfilled;
    teamId.toNumber().should.be.equal(ACADEMY_TEAM_ID);

    teamId = await market.getCurrentTeamIdFromPlayerId(id = 0).should.be.fulfilled;
    teamId.toNumber().should.be.equal(NULL_TEAMID);
    
    teamId = await market.getCurrentTeamIdFromPlayerId(id = 43214234).should.be.fulfilled;
    teamId.toNumber().should.be.equal(NULL_TEAMID);
  });
  
  
  it2("ownership functions of Academy Players", async () => {
    owner = await market.getOwnerTeam(ACADEMY_TEAM_ID).should.be.fulfilled;
    owner.should.be.equal(NULL_ADDR);

    tx = await assets.setAcademyAddr(freeverseAccount.address).should.be.fulfilled;

    owner = await market.getOwnerTeam(ACADEMY_TEAM_ID).should.be.fulfilled;
    owner.should.be.equal(freeverseAccount.address);

    playerId = await createSpecialPlayerId();

    teamId = await market.getCurrentTeamIdFromPlayerId(playerId).should.be.fulfilled;
    teamId.toNumber().should.be.equal(ACADEMY_TEAM_ID);

    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(freeverseAccount.address);

    is = await market.getIsSpecial(playerId).should.be.fulfilled;
    is.should.be.equal(true);

    was = await market.wasPlayerCreatedVirtually(playerId).should.be.fulfilled;
    was.should.be.equal(false);
  
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(false)
    
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, freeverseAccount).should.be.fulfilled;

    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);

    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(freeverseAccount.address);

    was = await market.wasPlayerCreatedVirtually(playerId).should.be.fulfilled;
    was.should.be.equal(false);
  });

  
  it2("special players: completes a PUT_FOR_SALE and AGREE_TO_BUY via MTXs", async () => {
    playerId = await createSpecialPlayerId();

    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, freeverseAccount).should.be.rejected;
    tx = await assets.setAcademyAddr(freeverseAccount.address).should.be.fulfilled;
    truffleAssert.eventEmitted(tx, "TeamTransfer", (event) => {
      return event.teamId.toNumber() == 1 && event.to == freeverseAccount.address;
    });

    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, freeverseAccount).should.be.fulfilled;

    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);

    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });
    
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, buyerTeamId, buyerAccount
    ).should.be.fulfilled;

    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(false);
    });
    let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    finalOwner.should.be.equal(buyerAccount.address);

    // test that Freeverse cannot put the same player again in the market
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, freeverseAccount).should.be.rejected;
    
    // test that the new owner can sell freely as always
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, buyerAccount).should.be.fulfilled;
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = false, sellerTeamId, sellerAccount
    ).should.be.fulfilled;
    

  });

  it2("buy now player", async () => {
    // TODO: add test that it fails if not sent from Academy address.
    // CURRENTLY: it works regardless of: await assets.setAcademyAddr(freeverseAccount.address).should.be.fulfilled;
    
    playerId = await createSpecialPlayerId(id = 4312432432);

    // it currently has no owner:
    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(NULL_ADDR);
    
    // set target to a bot team => should fail
    targetTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry2 = 2);
    tx = await market.transferBuyNowPlayer(playerId.toString(), targetTeamId).should.be.rejected;

    // set target to a Academy team => should fail
    tx = await market.transferBuyNowPlayer(playerId.toString(), ACADEMY_TEAM_ID).should.be.rejected;

    // set target to a team that does not exist => should fail
    targetTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 2, teamIdxInCountry2 = 2);
    was = await market.wasTeamCreatedVirtually(targetTeamId).should.be.fulfilled;
    was.should.be.equal(false);
    tx = await market.transferBuyNowPlayer(playerId.toString(), targetTeamId).should.be.rejected;
    
    // set target to a human team => should work
    targetTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry2 = 1);
    tx = await market.transferBuyNowPlayer(playerId.toString(), targetTeamId).should.be.fulfilled;
    
    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(buyerAccount.address);

    // a non academy player would not work
    playerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry = 4);
    tx = await market.transferBuyNowPlayer(playerId.toString(), targetTeamId).should.be.rejected;
  });

  it2("promo players: cannot be used to transfer players that already belong to humans", async () => {
    // We try to transfer an existing player as if it was a promo player
    // Note that teamIdx = 0 => seller, teamIdx = 1 => buyer
    // ATTEMPT 1: Human to null team
    FREEVERSE = accounts[0];
    playerId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, playerIdxInCountry = 4);

    result = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    result.should.be.equal(sellerAccount.address);

    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    var sigBuyerRS  = [sigBuyer.r, sigBuyer.s];

    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    // this one fails because targetTeamId = 0 for a standard playerId created as part of a league.
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
    
    // ATTEMPT 2: Human to bot team: player in teamIdx = 0, trying to transfer to a bot teamIdx = 2
    targetTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry2 = 2);
    result = await market.getOwnerTeam(targetTeamId).should.be.fulfilled;
    result.should.be.equal('0x0000000000000000000000000000000000000000')

    playerId = await market.setTargetTeamId(playerId, targetTeamId).should.be.fulfilled;
    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));

    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];

    // fails because targetTeam is a bot
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;

    // ATTEMPT 3: Human to human team: player in teamIdx = 0, trying to transfer to a bot teamIdx = 1
    targetTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry2 = 1);

    // make sure player belongs to "seller"
    result = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    result.should.be.equal(sellerAccount.address);

    result = await market.getOwnerTeam(targetTeamId).should.be.fulfilled;
    result.should.be.equal(buyerAccount.address)

    playerId = await market.setTargetTeamId(playerId, targetTeamId).should.be.fulfilled;
    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    var sigBuyerRS  = [sigBuyer.r, sigBuyer.s];

    // should be rejected because it is simply not an Academy player
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
  });
     
  it2("promo players: completes an offering and accepting", async () => {
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId).should.be.fulfilled;
    FREEVERSE = accounts[0];

    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    const sigBuyerRS  = [sigBuyer.r, sigBuyer.s];

    // it currently does not exist:
    finalPlayerId = await assets.setTargetTeamId(playerId, 0).should.be.fulfilled;

    // it currently has no owner since AcamedyAddr is not set:
    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(NULL_ADDR);
      
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
    // let's fix it:
    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(FREEVERSE);

    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.fulfilled;
    // change of academy address immediately reflects in change of who owns the academy players
    owner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
    owner.should.be.equal(FREEVERSE);
    
    // when transferred, the "targetTeamId" is erased (set to zero)
    finalPlayerId = await assets.setTargetTeamId(playerId, 0).should.be.fulfilled;
    owner = await market.getOwnerPlayer(finalPlayerId).should.be.fulfilled;
    owner.should.be.equal(buyerAccount.address);
  });
  
  it2("promo players: failures when sending to bot teams", async () => {
    FREEVERSE = accounts[0];
    buyerTeamId = await assets.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, teamIdxInCountry2 = 5);
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId).should.be.fulfilled;

    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    const sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    // fails because owner does not own the team (it is a bot, after all)
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
  });

  it2("promo players: failures when sending to Academy team", async () => {
    FREEVERSE = accounts[0];
    playerId = await createPromoPlayer(targetTeamId = ACADEMY_TEAM_ID).should.be.fulfilled;

    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    const sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    // fails because owner does not own the team (it is a bot, after all)
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
  });


  it2("promo players: effect of constraints", async () => {
    FREEVERSE = accounts[0];
    await market.addAcquisitionConstraint(buyerTeamId, valUnt = now.toNumber() + 1000, n = 1).should.be.fulfilled;
    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    // first acquisition works:
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId).should.be.fulfilled;
    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.fulfilled;
    // first acquisition fails:
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId, 432153).should.be.fulfilled;
    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
  });

  it2("promo players: a promo player cannot be acquired by any team other than targetTeam", async () => {
    FREEVERSE = accounts[0];
    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId).should.be.fulfilled;
    sigBuyer = await sellerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil])); // note the signer is not the targetTeam owner
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;
  });
  
  it2("promo players: cannot offer a promo player that already exists", async () => {
    FREEVERSE = accounts[0];
    await assets.setAcademyAddr(FREEVERSE).should.be.fulfilled;
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId).should.be.fulfilled;
    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.fulfilled;
    finalPlayerId = await assets.setTargetTeamId(playerId, 0).should.be.fulfilled;
    owner = await market.getOwnerPlayer(finalPlayerId).should.be.fulfilled;
    owner.should.be.equal(buyerAccount.address);

    // try to offer it again to the same buyer (literal copy-paste of previous paragraph)
    playerId = await createPromoPlayer(targetTeamId = buyerTeamId).should.be.fulfilled;
    sigBuyer = await buyerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil]));
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;

    // try to offer the same promo player to another user (e.g. seller)
    playerId = await createPromoPlayer(targetTeamId = sellerTeamId).should.be.fulfilled; // note the different target team
    sigSeller = await freeverseAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil])); 
    sigBuyer = await sellerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil])); // note the different signer
    sigSellerRS = [sigSeller.r, sigSeller.s];
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.rejected;

    // do double check: any other playerId would've worked
    playerId = await createPromoPlayer(targetTeamId = sellerTeamId, 54235342).should.be.fulfilled; // note the different target team
    sigBuyer = await sellerAccount.sign(marketUtils.concatHash(["uint256", "uint256"], [playerId.toString(), validUntil])); // note the different signer
    sigBuyerRS  = [sigBuyer.r, sigBuyer.s];
    tx = await market.transferPromoPlayer(playerId.toString(), validUntil, sigBuyerRS, sigBuyer.v, {from: FREEVERSE}).should.be.fulfilled;
  });
  
  it2("players: fails a PUT_FOR_SALE and AGREE_TO_BUY via MTXs because isOffer2StartAuction is not correctly set ", async () => {
    tx = await freezePlayer(currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;
    isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
    isPlayerFrozen.should.be.equal(true);
    truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
      return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
    });
    tx = await completePlayerAuction(
      currencyId, price,  sellerRnd, validUntil, playerId, 
      extraPrice, buyerRnd, isOffer2StartAuctionSig = false, isOffer2StartAuctionBC = true, buyerTeamId, buyerAccount
    ).should.be.rejected;
  });
  
  
  // OTHER TESTS
  
  it2("test accounts from truffle and web3", async () => {
    accountsWeb3 = await web3.eth.getAccounts().should.be.fulfilled;
    accountsWeb3[0].should.be.equal(accounts[0]);
  });
  
  it2('players: put for sale msg', async () => {
    const validUntil = 2000000000;
    const playerId = 10;
    const currencyId = 1;
    const price = 41234;
    const rnd = 42321;
    const sellerAccount = web3.eth.accounts.privateKeyToAccount('0x3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54');

    const sellerHiddenPrice = await market.hashPrivateMsg(currencyId, price, rnd).should.be.fulfilled;
    sellerHiddenPrice.should.be.equal('0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa');
    const message = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerId).should.be.fulfilled;
    message.should.be.equal('0x07d43490a59d38783f03854081c1ecd738a6cb320c1767befdbc147e6b496eed');
    const sigSeller = sellerAccount.sign(message);
    sigSeller.messageHash.should.be.equal('0xc50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796');
    sigSeller.signature.should.be.equal('0x075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c');
  });
  
   
  it2('players: deterministic sign (values used in market.notary test)', async () => {
    sellerTeamId.should.be.bignumber.equal('274877906944');
    buyerTeamId.should.be.bignumber.equal('274877906945');
    sellerTeamPlayerIds = await market.getPlayerIdsInTeam(sellerTeamId).should.be.fulfilled;
    const playerIdToSell = sellerTeamPlayerIds[0];
    playerIdToSell.should.be.bignumber.equal('274877906944');

    const sellerAccount = web3.eth.accounts.privateKeyToAccount('0x3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54');
    const buyerAccount = await web3.eth.accounts.privateKeyToAccount('0x3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882');

    // Define params of the seller, and sign
    validUntil = 2000000000;
    buyerHiddenPrice = marketUtils.concatHash(
      ["uint256", "uint256"],
      [extraPrice, buyerRnd]
    );

    const sellerHiddenPrice = await market.hashPrivateMsg(currencyId, price, sellerRnd).should.be.fulfilled;
    sellerHiddenPrice.should.be.equal('0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa');
    const message = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, playerIdToSell).should.be.fulfilled;

    message.should.be.equal('0x909e2fbc45b398649f58c7ea4b632ff1b457ee5f60a43a70abfe00d50e7c917d');
    const sigSeller = sellerAccount.sign(message);
    sigSeller.messageHash.should.be.equal('0x55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a');
    sigSeller.signature.should.be.equal('0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c');

    const prefixed = await market.prefixed(message).should.be.fulfilled;
    const isOffer2StartAuction = true;
    const buyerMsg = await market.buildAgreeToBuyPlayerTxMsg(prefixed, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction).should.be.fulfilled;
    buyerMsg.should.be.equal('0xc049e2032b873dd67cc7cc43fb2488f7c770d1654716fc75024cda693c74488c');
    const sigBuyer = buyerAccount.sign(buyerMsg);
    sigBuyer.messageHash.should.be.equal('0xe04d23ec0424b22adec87879118715ce75997a4fd47897c398f3a8cad79b3041');
    sigBuyer.signature.should.be.equal('0xdbe104e7b51c9b1e38cdda4e31c2036e531f7d3338d392bee2f526c4c892437f5e50ddd44224af8b3bd92916b93e4b0d7af2974175010323da7dedea19f30d621c');
  });

  it2('teams: deterministic sign (values used in market.notary test)', async () => {
    sellerTeamId.should.be.bignumber.equal('274877906944');

    const sellerAccount = web3.eth.accounts.privateKeyToAccount('0x3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54');
    const buyerAccount = await web3.eth.accounts.privateKeyToAccount('0x3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882');

    // Define params of the seller, and sign
    validUntil = 2000000000;
    buyerHiddenPrice = marketUtils.concatHash(
      ["uint256", "uint256"],
      [extraPrice, buyerRnd]
    );

    const sellerHiddenPrice = await market.hashPrivateMsg(currencyId, price, sellerRnd).should.be.fulfilled;
    sellerHiddenPrice.should.be.equal('0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa');
    const message = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, validUntil, sellerTeamId).should.be.fulfilled;

    message.should.be.equal('0x909e2fbc45b398649f58c7ea4b632ff1b457ee5f60a43a70abfe00d50e7c917d');
    const sigSeller = sellerAccount.sign(message);
    sigSeller.messageHash.should.be.equal('0x55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a');
    sigSeller.signature.should.be.equal('0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c');

    const prefixed = await market.prefixed(message).should.be.fulfilled;
    const isOffer2StartAuction = true;
    const buyerMsg = await market.buildAgreeToBuyTeamTxMsg(prefixed, buyerHiddenPrice, isOffer2StartAuction).should.be.fulfilled;
    buyerMsg.should.be.equal('0xdd3d39b424073a7a74a333d3b35bc2b0adea64c4a51c47c4669d190111e7b5e5');
    const sigBuyer = buyerAccount.sign(buyerMsg);
    sigBuyer.messageHash.should.be.equal('0xeb0feff7cbf76cd8f6a6bb07b2d92305e1978c66a157b7738e249689682942f7');
    sigBuyer.signature.should.be.equal('0x7c6b08dfff430bd5dd1785463846f3961f3844b9b4d1cccc941ad2d5441b4496556ffc4518f9be660e2c34ba3d74ef67665af727c25eae6758695354b36462f71b');
  });

});
