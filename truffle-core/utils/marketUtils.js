const truffleAssert = require('truffle-assertions');

function concatHash(types, vals) {
  assert(types.length == vals.length, "Length of inputs should be equal")
  return web3.utils.keccak256(
      web3.eth.abi.encodeParameters(types, vals)
  )
}

// this function does the crazy thing solidity does for hex...
function getMessageHash(msg) {
  assert(web3.utils.isHexStrict(msg), "We currently only support signing hashes, which are 0x stating hex numbers")
  message = web3.utils.hexToBytes(msg);
  var messageBuffer = Buffer.from(message);
  var preamble = "\x19Ethereum Signed Message:\n" + message.length;
  var preambleBuffer = Buffer.from(preamble);
  var ethMessage = Buffer.concat([preambleBuffer, messageBuffer]);
  return web3.utils.keccak256(ethMessage);
}

async function signPutAssetForSaleMTx(currencyId, price, rnd, validUntil, auctionDurationAfterOfferIsAccepted, asssetId, sellerAccount) {
  sellerTxMsg = getSellerTxMsg(currencyId, price, rnd, asssetId, validUntil, auctionDurationAfterOfferIsAccepted)
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

function getSellerTxMsg(currencyId, price, sellerRnd, assetId, validUntil, auctionDurationAfterOfferIsAccepted) {
  const sellerHiddenPrice = concatHash(
    ['uint8', 'uint256', 'uint256'],
    [currencyId, price, sellerRnd]
  );
  const sellerTxMsg = concatHash(
    ['bytes32', 'uint256', 'uint32', 'uint32'],
    [sellerHiddenPrice, assetId, validUntil, auctionDurationAfterOfferIsAccepted]
  )
  return sellerTxMsg;
}

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signAgreeToBuyPlayerMTx(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, auctionDurationAfterOfferIsAccepted, playerId, buyerTeamId, buyerAccount) {
  const buyerHiddenPrice = concatHash(
    ['uint256', 'uint256'],
    [extraPrice, buyerRnd]
  )
  sellerTxMsg = getSellerTxMsg(currencyId, price, sellerRnd, playerId, validUntil, auctionDurationAfterOfferIsAccepted);
  const sellerTxHash = getMessageHash(sellerTxMsg);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32', 'uint256'],
      [sellerTxHash, buyerHiddenPrice, buyerTeamId]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return [sigBuyer, sellerTxHash];
}

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signAgreeToBuyTeamMTx(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, auctionDurationAfterOfferIsAccepted, teamId, buyerAccount) {
  const buyerHiddenPrice = concatHash(
    ['uint256', 'uint256'],
    [extraPrice, buyerRnd]
  )
  sellerTxMsg = getSellerTxMsg(currencyId, price, sellerRnd, teamId, validUntil, auctionDurationAfterOfferIsAccepted)
  const sellerTxHash = getMessageHash(sellerTxMsg);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32'],
      [sellerTxHash, buyerHiddenPrice]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return [sigBuyer, sellerTxHash];
}


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

  tx = await market.freezePlayerViaPutForSale(
    sellerHiddenPrice,
    playerId,
    sigSellerRS = [NULL_BYTES32, NULL_BYTES32],
    sigSellerV = 0,
    validUntil,
    {from: contractOwner}
  ).should.be.fulfilled;
  return tx;
}

async function freezePlayer(contractOwner, currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, playerId, sellerAccount) {
  // Mobile app does this:
  sigSeller = await signPutAssetForSaleMTx(
    currencyId,
    price,
    sellerRnd,
    validUntil,
    auctionDurationAfterOfferIsAccepted,
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
  sellerTxMsgBC = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, playerId, validUntil, auctionDurationAfterOfferIsAccepted).should.be.fulfilled;
  sellerTxMsgBC.should.be.equal(sigSeller.message);

  // Then, the buyer builds a message to sign
  let isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(false);

  // and send the Freeze TX. 
  const sigSellerRS = [
    sigSeller.r,
    sigSeller.s
  ];
  var tx;
  
  isOffer = auctionDurationAfterOfferIsAccepted > 0;
  if (isOffer) {
    tx = await market.freezePlayerViaOffer(
      sellerHiddenPrice,
      playerId,
      sigSellerRS,
      sigSeller.v,
      validUntil,
      auctionDurationAfterOfferIsAccepted,
      {from: contractOwner}
    ).should.be.fulfilled;
  } else {
    tx = await market.freezePlayerViaPutForSale(
      sellerHiddenPrice,
      playerId,
      sigSellerRS,
      sigSeller.v,
      validUntil,
      {from: contractOwner}
    ).should.be.fulfilled;
  }
  return tx;
}

async function transferTeamViaAuction(contractOwner, market, sellerTeamId, sellerAccount, buyerAccount) {
  currencyId = 1;
  price = 41234;
  sellerRnd = 42321;
  extraPrice = 332;
  buyerRnd = 1243523;
  isOffer = false;
  auctionDurationAfterOfferIsAccepted = isOffer ? 3600*24 : 0;

  now = await market.getBlockchainNowTime().should.be.fulfilled;
  AUCTION_TIME = 48 * 3600;
  validUntil = now.toNumber() + AUCTION_TIME;
    
  tx = await freezeTeam(contractOwner, currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, sellerTeamId, sellerAccount).should.be.fulfilled;
  isTeamFrozen = await market.isTeamFrozen(sellerTeamId.toNumber()).should.be.fulfilled;
  isTeamFrozen.should.be.equal(true);
  truffleAssert.eventEmitted(tx, "TeamFreeze", (event) => {
    return event.teamId.should.be.bignumber.equal(sellerTeamId) && event.frozen.should.be.equal(true);
  });
  tx = await completeTeamAuction(
    contractOwner, 
    currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, sellerTeamId, 
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
  isOffer = false;
  auctionDurationAfterOfferIsAccepted = isOffer ? 3600*24 : 0;

  now = await market.getBlockchainNowTime().should.be.fulfilled;
  AUCTION_TIME = 48 * 3600;
  validUntil = now.toNumber() + AUCTION_TIME;
  tx = await freezePlayer(contractOwner, currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, playerId, sellerAccount).should.be.fulfilled;
  isPlayerFrozen = await market.isPlayerFrozenFiat(playerId).should.be.fulfilled;
  isPlayerFrozen.should.be.equal(true);
  truffleAssert.eventEmitted(tx, "PlayerFreeze", (event) => {
    return event.playerId.should.be.bignumber.equal(playerId) && event.frozen.should.be.equal(true);
  });
  
  tx = await completePlayerAuction(
    contractOwner, 
    currencyId, price,  sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, playerId, 
    extraPrice, buyerRnd, buyerTeamId, buyerAccount
  ).should.be.fulfilled;
  let finalOwner = await market.getOwnerPlayer(playerId).should.be.fulfilled;
  finalOwner.should.be.equal(buyerAccount.address);
}

async function completePlayerAuction(
  contractOwner, 
  currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, playerId, 
  extraPrice, buyerRnd, buyerTeamId, buyerAccount
) {
  // Add some amount to the price where seller started, and a rnd to obfuscate it
  const buyerHiddenPrice = concatHash(
    ["uint256", "uint256"],
    [extraPrice, buyerRnd]
  );

  var {0: sigBuyer, 1: sellerTxHash} = await signAgreeToBuyPlayerMTx(
    currencyId,
    price,
    extraPrice, // 0
    sellerRnd, // offererRnd
    buyerRnd,  // 0
    validUntil, // offerValidUntil
    auctionDurationAfterOfferIsAccepted,
    playerId.toString(),
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

  tx = await market.completePlayerAuction(
    sellerTxHash,
    playerId.toString(),
    buyerHiddenPrice,
    buyerTeamId.toString(),
    sigBuyerRS,
    sigBuyer.v,
    {from: contractOwner}
  ).should.be.fulfilled;

  return tx
}

async function freezeTeam(contractOwner, currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, teamId, sellerAccount) {
  // Mobile app does this:
  sigSeller = await signPutAssetForSaleMTx(
    currencyId,
    price,
    sellerRnd,
    validUntil, 
    auctionDurationAfterOfferIsAccepted,
    teamId.toString(),
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
  sellerTxMsgBC = await market.buildPutAssetForSaleTxMsg(sellerHiddenPrice, teamId.toString(), validUntil, auctionDurationAfterOfferIsAccepted).should.be.fulfilled;
  sellerTxMsgBC.should.be.equal(sigSeller.message);

  // Then, the buyer builds a message to sign
  let isTeamFrozen = await market.isTeamFrozen(teamId.toNumber()).should.be.fulfilled;
  isTeamFrozen.should.be.equal(false);

  // and send the Freeze TX. 
  const sigSellerRS = [
    sigSeller.r,
    sigSeller.s
  ];
  isOffer = auctionDurationAfterOfferIsAccepted > 0;
  if (isOffer) {
    tx = await market.freezeTeamViaOffer(
      sellerHiddenPrice,
      teamId.toNumber(),
      sigSellerRS,
      sigSeller.v,
      validUntil,
      auctionDurationAfterOfferIsAccepted,
      {from: contractOwner}
    ).should.be.fulfilled;
  } else {
    tx = await market.freezeTeamViaPutForSale(
      sellerHiddenPrice,
      teamId.toNumber(),
      sigSellerRS,
      sigSeller.v,
      validUntil,
      {from: contractOwner}
    ).should.be.fulfilled;
  }
  return tx;
};

async function completeTeamAuction(
  contractOwner,
  currencyId, price, sellerRnd, validUntil, auctionDurationAfterOfferIsAccepted, sellerTeamId, 
  extraPrice, buyerRnd, buyerAccount) 
{
  // Add some amount to the price where seller started, and a rnd to obfuscate it
  const buyerHiddenPrice = concatHash(
    ["uint256", "uint256"],
    [extraPrice, buyerRnd]
  );
  var {0: sigBuyer, 1: sellerTxHash} = await signAgreeToBuyTeamMTx(
    currencyId,
    price,
    extraPrice,
    sellerRnd,
    buyerRnd,
    validUntil,
    auctionDurationAfterOfferIsAccepted,
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
    sellerTxHash,
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
  getMessageHash,
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
  freezeAcademyPlayer
}