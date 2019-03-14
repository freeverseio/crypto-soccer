var ethUtil = require('ethereumjs-util');
var HDKey = require('hdkey')
var bip39 = require('bip39')

let bytesToHex = function bytesToHex(buff) {
  return `0x${buff.toString('hex')}`;
}

let hexToBytes = function hexToBytes(hex) {
  if (hex.substr(0, 2) === '0x') {
    return Buffer.from(hex.substr(2), 'hex');
  }

  return Buffer.from(hex, 'hex');
}

let generateKeysMnemonic = function generateKeysMnemonic(mnemonic) {
  if (mnemonic == undefined) {
    mnemonic = bip39.generateMnemonic();
  }

  let len = mnemonic.trim().split(/\s+/g).length
  if (len < 12) {
    throw "Mnemonic phrase needs 12 word at least. Only " + len + ' given.'
  }

  const hdkey = HDKey.fromMasterSeed(mnemonic);
  const masterPrivateKey = hdkey.privateKey;
  const masterPubKey = hdkey.publicKey;
  var path = "m/44'/60'/0'/0/0";
  const addrNode = hdkey.derive(path);
  let privK = addrNode._privateKey;
  const pubKey = ethUtil.privateToPublic(addrNode._privateKey);
  let address = ethUtil.privateToAddress(addrNode._privateKey);
  let addressHex = bytesToHex(address);
  let privKHex = bytesToHex(privK);
  //localStorage.setItem(addressHex, privKHex);
  return {address: addressHex, mnemonic: mnemonic};
}

module.exports = {
  hexToBytes, bytesToHex, generateKeysMnemonic
}
