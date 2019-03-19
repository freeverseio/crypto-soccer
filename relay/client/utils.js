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
    let prefix = buf(uint8(0x19)).toString('hex')
    //TODO: ask adria "\x19Ethereum Signed Message:\n"
    let hexType = buf(type).toString('hex')
    let hexValue = buf(value).toString('hex')
    let msg = "0x" + prefix +
      buf(uint8(0)).toString('hex') +
      buf(useraddr).toString('hex') +
      buf(uint256(nonce)).toString('hex') +
      hexType + hexValue

      let signedData = ethUtil.ecsign(buf(sha3(msg)),buf(privatekey));
      console.log("signedData:", signedData);
      let txData = {
              from: useraddr,
              type: hexType,
              value: hexValue,
              r: signedData.r.toString('hex'),
              s: signedData.s.toString('hex'),
              v: signedData.v
          };
      console.log("txData:", txData)
      return axios.post(userURL + '/action', txData)
  })
}

module.exports = {
  createUser, generateKeysMnemonic, submitAction
}
