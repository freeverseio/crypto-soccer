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


module.exports = {
  concatHash,
  getMessageHash,
  signPutAssetForSaleMTx
}