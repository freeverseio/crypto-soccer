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


async function signPutAssetForSaleMTx(currencyId, price, rnd, validUntil, asssetId, sellerAccount) {
  const hiddenPrice = concatHash(
      ['uint8', 'uint256', 'uint256'],
      [currencyId, price, rnd]
  )
  const sellerTxMsg = concatHash(
      ['bytes32', 'uint256', 'uint256'],
      [hiddenPrice, validUntil, asssetId]
  )
  
  const sigSeller = await sellerAccount.sign(sellerTxMsg);
  sigSeller.message.should.be.equal(sellerTxMsg);
  return sigSeller;
}

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signAgreeToBuyPlayerMTx(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, playerId, isOffer2StartAuction, buyerTeamId, buyerAccount) {
  const sellerHiddenPrice = concatHash(
    ['uint8', 'uint256', 'uint256'],
    [currencyId, price, sellerRnd]
  )
  const buyerHiddenPrice = concatHash(
    ['uint256', 'uint256'],
    [extraPrice, buyerRnd]
  )
  const sellerTxMsg = concatHash(
      ['bytes32', 'uint256', 'uint256'],
      [sellerHiddenPrice, validUntil, playerId]
  )
  const sellerTxHash = getMessageHash(sellerTxMsg);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32', 'uint256', 'bool'],
      [sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return sigBuyer;
}

// Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
async function signAgreeToBuyTeamMTx(currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, playerId, isOffer2StartAuction, buyerAccount) {
  const sellerHiddenPrice = concatHash(
    ['uint8', 'uint256', 'uint256'],
    [currencyId, price, sellerRnd]
  )
  const buyerHiddenPrice = concatHash(
    ['uint256', 'uint256'],
    [extraPrice, buyerRnd]
  )
  const sellerTxMsg = concatHash(
      ['bytes32', 'uint256', 'uint256'],
      [sellerHiddenPrice, validUntil, playerId]
  )
  const sellerTxHash = getMessageHash(sellerTxMsg);
  buyerTxMsg = concatHash(
      ['bytes32', 'bytes32', 'bool'],
      [sellerTxHash, buyerHiddenPrice, isOffer2StartAuction]
  )
  const sigBuyer = await buyerAccount.sign(buyerTxMsg);
  return sigBuyer;
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


module.exports = {
  concatHash,
  getMessageHash,
  signPutAssetForSaleMTx,
  signAgreeToBuyPlayerMTx,
  signAgreeToBuyTeamMTx,
  // signOfferToBuyPlayerMTx,
  // signOfferToBuyTeamMTx,  
  // buildOfferToBuyPlayerMsg,
  // buildOfferToBuyTeamMsg,
}