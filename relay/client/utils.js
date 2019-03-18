let ethUtil = require('ethereumjs-util');
let Web3 = require('web3');
let HDKey = require('hdkey')
let bip39 = require('bip39')
let toastr = require('toastr')
let axios = require('axios')


var web3 = new Web3(new Web3.providers.HttpProvider(''));
const buf = b => ethUtil.toBuffer(b)
const sha3 = b => web3.utils.soliditySha3(b)
const uint256 = n => "0x"+n.toString(16).padStart(64,'0')
const uint8 = n => "0x"+n.toString(16)

const RELAYURL = 'http://localhost:8080';

function bytesToHex(buff) {
  return `0x${buff.toString('hex')}`;
}

function hexToBytes(hex) {
  if (hex.substr(0, 2) === '0x') {
    return Buffer.from(hex.substr(2), 'hex');
  }

  return Buffer.from(hex, 'hex');
}

function createUser(addr) {
  return axios.post(RELAYURL + '/relay/createuser', {useraddr:addr})
}

function generateKeysMnemonic(mnemonic) {
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
  return {address: addressHex, privatekey: privKHex, mnemonic: mnemonic};
}

function submitAction(useraddr, privatekey, type, value) {
  if(useraddr==undefined) {
      toastr.error("Undefined address");
      return;
  }
  if(useraddr=="") { // TODO check also if it's a valid eth address
      toastr.error("Empty address");
      return;
  }

  let nonce = -1;
  let userURL = RELAYURL + '/relay/v1/' + useraddr
  return axios.get(userURL + '/nonce')
  .then(function(res) {
    nonce = res.data.nonce;
    console.log("response from server:", res.data);

    // after getting nonce, generate & sign & send transaction
    let msg = "0x" +
      buf(uint8(0x19)).toString('hex') +
      buf(uint8(0)).toString('hex') +
      buf(useraddr).toString('hex') +
      buf(uint256(nonce)).toString('hex') +
      buf(type).toString('hex') +
      buf(value).toString('hex')

      let sig = ethUtil.ecsign(buf(sha3(msg)),buf(privatekey));
      console.log(sig);
      stringSig = sig.r.toString('hex') + sig.s.toString('hex') + buf(uint8(sig.v)).toString('hex')
      return axios.post(userURL + '/action', {type:type, value:stringSig})
      // TODO: send sig and not stringSig
      //return axios.post(userURL + '/action', {type:type, value:sig})
  })
}

module.exports = {
  createUser, generateKeysMnemonic, submitAction
}
