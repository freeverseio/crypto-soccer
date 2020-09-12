const truffleAssert = require('truffle-assertions');

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
  const hiddenPrice = hidePrice(currencyId, price, rnd);

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

function computePutAssetForSaleDigestNoPrefix(currencyId, price, sellerRnd, validUntil, offerValidUntil, assetId) {
  const sellerHiddenPrice = hidePrice(currencyId, price, sellerRnd);

  const sellerTxMsg = concatHash(
    ['bytes32', 'uint256', 'uint32', 'uint32'],
    [sellerHiddenPrice, assetId, validUntil, offerValidUntil]
  );
  return sellerTxMsg;
}

function computePutAssetForSaleDigest(currencyId, price, sellerRnd, validUntil, offerValidUntil, assetId) {
  return prefix(computePutAssetForSaleDigestNoPrefix(currencyId, price, sellerRnd, validUntil, offerValidUntil, assetId));
}

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signAgreeToBuyPlayerMTx(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, offerValidUntil, playerId, buyerTeamId, buyerAccount) {
  const buyerHiddenPrice = concatHash(
    ['uint256', 'uint256'],
    [extraPrice, buyerRnd]
  )
  const sellerDigest = computePutAssetForSaleDigest(currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32', 'uint256'],
      [sellerDigest, buyerHiddenPrice, buyerTeamId]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  auctionId = computeAuctionId(currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId);
  return [sigBuyer, auctionId];
}

function hidePrice(currencyId, price, rnd) {
  return concatHash(
    ['uint8', 'uint256', 'uint256'],
    [currencyId, price, rnd]
  );
}

function computeAuctionId(currencyId, price, sellerRnd, assetId, validUntil, offerValidUntil) {
  const sellerHiddenPrice = hidePrice(currencyId, price, sellerRnd);

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
  const sellerDigest = computePutAssetForSaleDigest(currencyId, price, sellerRnd, validUntil, offerValidUntil, teamId);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32'],
      [sellerDigest, buyerHiddenPrice]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return [sigBuyer, sellerDigest];
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

async function freezePlayer(contractOwner, currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId, sellerAccount) {
  // Mobile app does this:
  sigSeller = await signPutAssetForSaleMTx(
    currencyId,
    price,
    sellerRnd,
    validUntil,
    offerValidUntil,
    playerId.toString(),
    sellerAccount
  );
  // First of all, Freeverse and Buyer check the signature
  // In this case, using web3:
  recoveredSellerAddr = await web3.eth.accounts.recover(sigSeller);
  recoveredSellerAddr.should.be.equal(sellerAccount.address);

  // The correctness of the seller message can also be checked in the BC:
  const sellerHiddenPrice = concatHash(
    ["uint8", "uint256", "uint256"],
    [currencyId, price, sellerRnd]
  );
  sellerDigest = await market.computePutAssetForSaleDigest(sellerHiddenPrice, playerId, validUntil, offerValidUntil).should.be.fulfilled;
  sellerDigest.should.be.equal(prefix(sigSeller.message)); 

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
    playerId,
    sigSellerRS,
    sigSeller.v,
    validUntil,
    offerValidUntil,
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
  currencyId = 1;
  price = 41234;
  sellerRnd = 42321;
  extraPrice = 332;
  buyerRnd = 1243523;
  offerValidUntil = 0;
  
  now = await market.getBlockchainNowTime().should.be.fulfilled;
  AUCTION_TIME = 48 * 3600;
  validUntil = now.toNumber() + AUCTION_TIME;

  tx = await freezePlayer(contractOwner, currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId, sellerAccount).should.be.fulfilled;
  isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(true);
  truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
    return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
  });
  tx = await completePlayerAuction(
    contractOwner, 
    currencyId, price,  sellerRnd, validUntil, offerValidUntil, playerId, 
    extraPrice, buyerRnd,buyerTeamId, buyerAccount
  ).should.be.fulfilled;
  let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
  finalOwner.should.be.equal(buyerAccount.address);
}

async function completePlayerAuction(
  contractOwner, 
  currencyId, price, sellerRnd, validUntil, offerValidUntil, playerId,
  extraPrice, buyerRnd, buyerTeamId, buyerAccount
) {

  // Add some amount to the price where seller started, and a rnd to obfuscate it
  const buyerHiddenPrice = concatHash(
    ["uint256", "uint256"],
    [extraPrice, buyerRnd]
  );

  console.log(validUntil, offerValidUntil);
  var {0: sigBuyer, 1: auctionId} = await signAgreeToBuyPlayerMTx(
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
  console.log(auctionId);

  // Freeverse checks the signature
  recoveredBuyerAddr = await web3.eth.accounts.recover(sigBuyer);
  recoveredBuyerAddr.should.be.equal(buyerAccount.address);

  // Freeverse waits until actual money has been transferred between users, and completes sale
  const sigBuyerRS = [
    sigBuyer.r,
    sigBuyer.s
  ];

  tx = await market.completePlayerAuction(
    auctionId,
    playerId.toString(),
    buyerHiddenPrice,
    buyerTeamId.toString(),
    sigBuyerRS,
    sigBuyer.v,
    {from: contractOwner}
  ).should.be.fulfilled;

  return tx
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
  var {0: sigBuyer, 1: sellerDigest} = await signAgreeToBuyTeamMTx(
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
    sellerDigest,
    sellerTeamId.toNumber(),
    buyerHiddenPrice,
    sigBuyerRS,
    sigBuyer.v,
    recoveredBuyerAddr,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx;
}


module.exports = {
  concatHash,
  prefix,
  signPutAssetForSaleMTx,
  signAgreeToBuyPlayerMTx,
  signAgreeToBuyTeamMTx,
  transferPlayerViaAuction,
  freezePlayer,
  completePlayerAuction,
  freezeTeam,
  completeTeamAuction,
  signDismissPlayerMTx,
  transferTeamViaAuction,
  freezeAcademyPlayer,
  computePutAssetForSaleDigest,
  computePutAssetForSaleDigestNoPrefix,
  computeAuctionId,
  hidePrice
}