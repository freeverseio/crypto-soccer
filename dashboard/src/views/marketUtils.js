
function assert(boolean, msg) {
    if (boolean == false) {
        console.log("WARNING! ASSERTION FAILED: " + msg);
    }
}

function concatHash(web3, types, vals) {
    assert(types.length == vals.length, "Length of inputs should be equal")
    return web3.utils.keccak256(
        web3.eth.abi.encodeParameters(types, vals)
    )
  }
  
  // this function does the crazy thing solidity does for hex...
  function getMessageHash(web3, msg) {
    assert(web3.utils.isHexStrict(msg), "We currently only support signing hashes, which are 0x stating hex numbers")
    var message = web3.utils.hexToBytes(msg);
    var messageBuffer = Buffer.from(message);
    var preamble = "\x19Ethereum Signed Message:\n" + message.length;
    var preambleBuffer = Buffer.from(preamble);
    var ethMessage = Buffer.signature([preambleBuffer, messageBuffer]);
    return web3.utils.keccak256(ethMessage);
  }
  
  
  async function signPutAssetForSaleMTx(web3, currencyId, price, rnd, validUntil, asssetId, sellerAccount) {
    const hiddenPrice = concatHash(
        web3,
        ['uint8', 'uint256', 'uint256'],
        [currencyId, price, rnd]
    )
    const sellerTxMsg = concatHash(
        web3,
        ['bytes32', 'uint256', 'uint256'],
        [hiddenPrice, validUntil, asssetId]
    )
    
    const sigSeller = await sellerAccount.sign(sellerTxMsg);
    return sigSeller;
  }
  
  // Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
  async function signAgreeToBuyPlayerMTx(web3, currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, playerId, isOffer2StartAuction, buyerTeamId, buyerAccount) {
    const sellerHiddenPrice = concatHash(
        web3,
        ['uint8', 'uint256', 'uint256'],
        [currencyId, price, sellerRnd]
    )
    const buyerHiddenPrice = concatHash(
        web3,
        ['uint256', 'uint256'],
        [extraPrice, buyerRnd]
    )
    const sellerTxMsg = concatHash(
        web3,
        ['bytes32', 'uint256', 'uint256'],
        [sellerHiddenPrice, validUntil, playerId]
    )
    const sellerTxHash = getMessageHash(sellerTxMsg);
    const buyerTxMsg = concatHash(
        web3,
        ['bytes32', 'bytes32', 'uint256', 'bool'],
        [sellerTxHash, buyerHiddenPrice, buyerTeamId, isOffer2StartAuction]
    )
    const sigBuyer = await buyerAccount.sign(buyerTxMsg);
    return sigBuyer;
  }
  
  // Buyer explicitly agrees to all of sellers data, and only adds the 'buyerTeamId' to it.
  async function signAgreeToBuyTeamMTx(web3, currencyId, price, extraPrice, sellerRnd, buyerRnd, validUntil, playerId, isOffer2StartAuction, buyerAccount) {
    const sellerHiddenPrice = concatHash(
        web3,
        ['uint8', 'uint256', 'uint256'],
        [currencyId, price, sellerRnd]
    )
    const buyerHiddenPrice = concatHash(
        web3,
        ['uint256', 'uint256'],
        [extraPrice, buyerRnd]
    )
    const sellerTxMsg = concatHash(
        web3,
        ['bytes32', 'uint256', 'uint256'],
        [sellerHiddenPrice, validUntil, playerId]
    )
    const sellerTxHash = getMessageHash(sellerTxMsg);
    const buyerTxMsg = concatHash(
        web3,
        ['bytes32', 'bytes32', 'bool'],
        [sellerTxHash, buyerHiddenPrice, isOffer2StartAuction]
    )
    const sigBuyer = await buyerAccount.sign(buyerTxMsg);
    return sigBuyer;
  }
  
  module.exports = {
    concatHash,
    getMessageHash,
    signPutAssetForSaleMTx,
    signAgreeToBuyPlayerMTx,
    signAgreeToBuyTeamMTx
  }