const truffleAssert = require('truffle-assertions');


// Signing  functions

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signPlayerBid(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, offerValidUntil, playerId, buyerTeamId, buyerAccount) {
  const buyerHiddenPrice = hideBuyerPrice(extraPrice, buyerRnd);
  const auctionId = computeAuctionId(currencyId, price, sellerRnd, playerId, validUntil, offerValidUntil);
  const buyerTxMsg = concatHash(
      ['bytes32', 'bytes32', 'uint256'],
      [auctionId, buyerHiddenPrice, buyerTeamId]
  );
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return [sigBuyer, auctionId, buyerHiddenPrice];
}



// Main functions that write to the BC


async function acceptOffer(contractOwner, currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId, sellerAccount) {
  const sigSeller = await signPutAssetForSaleMTx(
    currencyId,
    price,
    sellerRnd,
    validUntil,
    offerValidUntil,
    playerId.toString(),
    sellerAccount
  );

  const isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(false);

  const tx = await market.freezePlayer(
    hideSellerPrice(currencyId, price, sellerRnd),
    playerId,
    [sigSeller.r, sigSeller.s],
    sigSeller.v,
    validUntil,
    offerValidUntil,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx;
}


async function putPlayerForSale(contractOwner, currencyId, price, sellerRnd, validUntil, playerId, sellerAccount) {
  const offerValidUntil = 0;
  return acceptOffer(contractOwner, currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId, sellerAccount);
}

async function completePlayerAuction(
  contractOwner, 
  auctionId, 
  buyerHiddenPrice,
  playerId,
  buyerTeamId,
  sigBuyer
) {
  const tx = await market.completePlayerAuction(
    auctionId,
    playerId,
    buyerHiddenPrice,
    buyerTeamId,
    [sigBuyer.r, sigBuyer.s],
    sigBuyer.v,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx
}


// Internal functions

function concatHash(types, vals) {
  assert(types.length == vals.length, "Length of inputs should be equal")
  return web3.utils.keccak256(
      web3.eth.abi.encodeParameters(types, vals)
  )
}

// this function does the crazy thing solidity does for hex...
function prefix(msg) {
  assert(web3.utils.isHexStrict(msg), "We currently only support signing hashes, which are 0x stating hex numbers")
  message = web3.utils.hexToBytes(msg);
  var messageBuffer = Buffer.from(message);
  var preamble = "\x19Ethereum Signed Message:\n" + message.length;
  var preambleBuffer = Buffer.from(preamble);
  var ethMessage = Buffer.concat([preambleBuffer, messageBuffer]);
  return web3.utils.keccak256(ethMessage);
}

async function signPutAssetForSaleMTx(currencyId, price, rnd, validUntil, offerValidUntil, asssetId, sellerAccount) {
  const hiddenPrice = hideSellerPrice(currencyId, price, rnd);

  const sellerTxMsg = concatHash(
      ['bytes32', 'uint256', 'uint32', 'uint32'],
      [hiddenPrice, asssetId, validUntil, offerValidUntil]
  )
  
  const sigSeller = await sellerAccount.sign(sellerTxMsg);
  sigSeller.message.should.be.equal(sellerTxMsg);
  return sigSeller;
}

async function signDismissPlayerMTx(validUntil, asssetId, sellerAccount) {
  const sellerTxMsg = concatHash(
      ['uint256', 'uint256'],
      [validUntil, asssetId]
  )
  const sigSeller = await sellerAccount.sign(sellerTxMsg);
  sigSeller.message.should.be.equal(sellerTxMsg);
  return sigSeller;
}

function computePutAssetForSaleDigestNoPrefixFromHiddenPrice(sellerHiddenPrice, assetId, validUntil, offerValidUntil) {
  return concatHash(
    ['bytes32', 'uint256', 'uint32', 'uint32'],
    [sellerHiddenPrice, assetId, validUntil, offerValidUntil]
  );
}

function computePutAssetForSaleDigestNoPrefix(currencyId, price, sellerRnd, validUntil, offerValidUntil, assetId) {
  return computePutAssetForSaleDigestNoPrefixFromHiddenPrice(hideSellerPrice(currencyId, price, sellerRnd), assetId, validUntil, offerValidUntil);
}

function computePutAssetForSaleDigest(currencyId, price, sellerRnd, validUntil, offerValidUntil, assetId) {
  return prefix(computePutAssetForSaleDigestNoPrefix(currencyId, price, sellerRnd, validUntil, offerValidUntil, assetId));
}



function computeAuctionId(currencyId, price, sellerRnd, assetId, validUntil, offerValidUntil) {
  const sellerHiddenPrice = hideSellerPrice(currencyId, price, sellerRnd);

  return (offerValidUntil == 0) ?
    concatHash(['bytes32', 'uint256', 'uint32'], [sellerHiddenPrice, assetId, validUntil]) :
    concatHash(['bytes32', 'uint256', 'uint32'], [sellerHiddenPrice, assetId, offerValidUntil]);
}

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signAgreeToBuyTeamMTx(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, offerValidUntil, teamId, buyerAccount) {
  const buyerHiddenPrice = concatHash(
    ['uint256', 'uint256'],
    [extraPrice, buyerRnd]
  )
  const auctionId = computeAuctionId(currencyId, price, sellerRnd, teamId, validUntil, offerValidUntil);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32'],
      [auctionId, buyerHiddenPrice]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return [sigBuyer, auctionId];
}


// function buildOfferToBuyPlayerMsg(currencyId, price, rnd, validUntil, playerId, buyerTeamId) {
//   const privHash = concatHash(
//       ['uint8', 'uint256', 'uint256'],
//       [currencyId, price, rnd]
//   )
//   const buyerTxMsg = concatHash(
//       ['bytes32', 'uint256', 'uint256', 'uint256'],
//       [privHash, validUntil, playerId, buyerTeamId]
//   )
//   return buyerTxMsg;
// }

// function buildOfferToBuyTeamMsg(currencyId, price, rnd, validUntil, playerId) {
//   const privHash = concatHash(
//       ['uint8', 'uint256', 'uint256'],
//       [currencyId, price, rnd]
//   )
//   const buyerTxMsg = concatHash(
//       ['bytes32', 'uint256', 'uint256'],
//       [privHash, validUntil, playerId]
//   )
//   return buyerTxMsg;
// }

// async function signOfferToBuyPlayerMTx(currencyId, price, rnd, validUntil, playerId, buyerTeamId, buyerAccount) {
//   const buyerTxMsg = buildOfferToBuyPlayerMsg(currencyId, price, rnd, validUntil, playerId, buyerTeamId);
//   const sigBuyer = await buyerAccount.sign(buyerTxMsg);
//   sigBuyer.message.should.be.equal(buyerTxMsg);
//   return sigBuyer;
// }

// async function signOfferToBuyTeamMTx(currencyId, price, rnd, validUntil, playerId, buyerAccount) {
//   const buyerTxMsg = buildOfferToBuyPlayerMsg(currencyId, price, rnd, validUntil, playerId);
//   const sigBuyer = await buyerAccount.sign(buyerTxMsg);
//   sigBuyer.message.should.be.equal(buyerTxMsg);
//   return sigBuyer;
// }

async function freezeAcademyPlayer(contractOwner, currencyId, price, sellerRnd, validUntil, playerId) {
  // The correctness of the seller message can also be checked in the BC:
  const NULL_BYTES32 = web3.eth.abi.encodeParameter('bytes32','0x0');

  const sellerHiddenPrice = concatHash(
    ["uint8", "uint256", "uint256"],
    [currencyId, price, sellerRnd]
  );
  
  // Then, the buyer builds a message to sign
  let isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(false);

  tx = await market.freezePlayer(
    sellerHiddenPrice,
    playerId,
    sigSellerRS = [NULL_BYTES32, NULL_BYTES32],
    sigSellerV = 0,
    validUntil,
    offValidUntil = 0,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx;
}


async function transferTeamViaAuction(contractOwner, market, sellerTeamId, sellerAccount, buyerAccount) {
  currencyId = 1;
  price = 41234;
  sellerRnd = 42321;
  extraPrice = 332;
  buyerRnd = 1243523;
  offerValidUntil = 0;

  now = await market.getBlockchainNowTime().should.be.fulfilled;
  AUCTION_TIME = 48 * 3600;
  validUntil = now.toNumber() + AUCTION_TIME;
    
  tx = await freezeTeam(contractOwner, currencyId, price, sellerRnd, validUntil, offerValidUntil, sellerTeamId, sellerAccount).should.be.fulfilled;
  isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
  isTeamFrozen.should.be.equal(true);
  truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
    return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(true);
  });
  
  tx = await completeTeamAuction(
    contractOwner, 
    currencyId, price, sellerRnd, validUntil, offerValidUntil, sellerTeamId, 
    extraPrice, buyerRnd, buyerAccount
  ).should.be.fulfilled;
  
  truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
    return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(false);
  });

  let finalOwner = await market.getOwnerTeam(sellerTeamId.toNumber()).should.be.fulfilled;
  finalOwner.should.be.equal(buyerAccount.address);
  return tx;
}



async function transferPlayerViaAuction(contractOwner, market, playerId, buyerTeamId, sellerAccount, buyerAccount) {
  const currencyId = 1;
  const price = 41234;
  const sellerRnd = 42321;
  const extraPrice = 332;
  const buyerRnd = 1243523;
  const offerValidUntil = 0;
  
  const now = await market.getBlockchainNowTime().should.be.fulfilled;
  const AUCTION_TIME = 48 * 3600;
  const validUntil = now.toNumber() + AUCTION_TIME;

  const tx = await putPlayerForSale(contractOwner, currencyId, price, sellerRnd, validUntil, playerId, sellerAccount).should.be.fulfilled;

  const isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(true);
  truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
    return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
  });

  const {0: sigBuyer, 1: auctionId, 2: buyerHiddenPrice} = await signPlayerBid(
    currencyId,
    price,
    extraPrice,
    sellerRnd,
    buyerRnd,
    validUntil,
    offerValidUntil,
    playerId.toString(),
    buyerTeamId.toString(),
    buyerAccount
  ).should.be.fulfilled;

  await completePlayerAuction(
    contractOwner, 
    auctionId, 
    buyerHiddenPrice,
    playerId,
    buyerTeamId,
    sigBuyer
  ).should.be.fulfilled;
  const finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
  finalOwner.should.be.equal(buyerAccount.address);
}



async function freezeTeam(contractOwner, currencyId, price, sellerRnd, validUntil, offerValidUntil, teamId, sellerAccount) {
  // Mobile app does this:
  sigSeller = await signPutAssetForSaleMTx(
    currencyId,
    price,
    sellerRnd,
    validUntil, 
    offerValidUntil,
    teamId.toString(),
    sellerAccount
  );

  // First of all, Freeverse and Buyer check the signature
  // In this case, using web3:
  recoveredSellerAddr = await web3.eth.accounts.recover(sigSeller);
  recoveredSellerAddr.should.be.equal(sellerAccount.address);

  const sellerHiddenPrice = concatHash(
    ["uint8", "uint256", "uint256"],
    [currencyId, price, sellerRnd]
  );
  sellerDigest = await market.computePutAssetForSaleDigest(sellerHiddenPrice, teamId.toString(), validUntil, offerValidUntil).should.be.fulfilled;
  sellerDigest.should.be.equal(prefix(sigSeller.message));

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
    teamId.toString(),
    sigSellerRS,
    sigSeller.v,
    validUntil,
    offerValidUntil,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx;
};

async function completeTeamAuction(
  contractOwner,
  currencyId, price, sellerRnd, validUntil, offerValidUntil, sellerTeamId, 
  extraPrice, buyerRnd, buyerAccount) 
{
  // Add some amount to the price where seller started, and a rnd to obfuscate it
  const buyerHiddenPrice = concatHash(
    ["uint256", "uint256"],
    [extraPrice, buyerRnd]
  );
  var {0: sigBuyer, 1: auctionId} = await signAgreeToBuyTeamMTx(
    currencyId,
    price,
    extraPrice,
    sellerRnd,
    buyerRnd,
    validUntil,
    offerValidUntil,
    sellerTeamId.toNumber(),
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

  tx = await market.completeTeamAuction(
    auctionId,
    sellerTeamId.toNumber(),
    buyerHiddenPrice,
    sigBuyerRS,
    sigBuyer.v,
    recoveredBuyerAddr,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx;
}

function hideSellerPrice(currencyId, price, rnd) {
  return concatHash(['uint8', 'uint256', 'uint256'], [currencyId, price, rnd]);
}

function hideBuyerPrice(price, rnd) {
  return concatHash(['uint256', 'uint256'], [extraPrice, buyerRnd]);
}

module.exports = {
  concatHash,
  prefix,
  signPutAssetForSaleMTx,
  signPlayerBid,
  signAgreeToBuyTeamMTx,
  transferPlayerViaAuction,
  completePlayerAuction,
  freezeTeam,
  completeTeamAuction,
  signDismissPlayerMTx,
  transferTeamViaAuction,
  freezeAcademyPlayer,
  computePutAssetForSaleDigest,
  computePutAssetForSaleDigestNoPrefix,
  computePutAssetForSaleDigestNoPrefixFromHiddenPrice,
  computeAuctionId,
  hideSellerPrice
}