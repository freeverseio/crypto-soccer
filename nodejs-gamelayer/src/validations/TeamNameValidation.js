
const crypto = require('crypto')
const elliptic = require('elliptic')
let sha3 = require('js-sha3');
const { ethers } = require('ethers')
class TeamNameValidation {
  constructor({ teamId, name, signature, web3}) {
    this.teamId = teamId
    this.name = name
    this.signature = signature
    this.web3 = web3
    this.ec = new elliptic.ec('secp256k1')
  }

  hash() {
    
    const params = this.web3.eth.abi.encodeParameters( ["uint256", "string"], [this.teamId || 0, this.name || ''] )
    return this.web3.utils.soliditySha3(params);
  }

  verifySignature() {

    // let keyPair = this.ec.keyFromPrivate("97ddae0f3a25b92268175400149d65d6887b9cefaf28ea2c078e05cdc15a3c0a");
    // let privKey = keyPair.getPrivate("hex");
    // let pubKey = keyPair.getPublic();
    // console.log(`Private key: ${privKey}`);
    // console.log("Public key :", pubKey.encode("hex").substr(2));
    // console.log("Public key (compressed):",
    //     pubKey.encodeCompressed("hex"));
    
    // let msg = 'Message for signing';
    // let msgHash = sha3.keccak256(msg);
    // let signature = this.ec.sign(msgHash, privKey, "hex", {canonical: true});
    // console.log(`Msg: ${msg}`);
    // console.log(`Msg hash: ${msgHash}`);
    // console.log("Signature:", signature);

    const hash = this.prefixedHash()
    console.log("verifySignature -> prefixedhash", hash)
    console.log("verifySignature -> hash", this.hash())
    // let hexToDecimal = (x) => this.ec.keyFromPrivate(x, "hex").getPrivate().toString(10);
    // const signatureObject = {
    //   messageHash: hash,
    //   // r: toHexString(this.signature.split('').map(c => c.charCodeAt(c) || 0).slice(0, 31)),
    //   r: this.signature.split('').slice(0, 63).join(''),
    //   // s: toHexString(this.signature.split('').map(c => c.charCodeAt(c) || 0).slice(31, 63)),
    //   s: this.signature.split('').slice(64, 127).join(''),
    //   recoveryParam: 0,
    //   // signature: this.signature
    // }
    // const pubKeyRecovered = this.ec.recoverPubKey(hexToDecimal(hash), signatureObject, 0, "hex")
    // console.log("verifySignature -> pubKeyRecovered", pubKeyRecovered.encodeCompressed("hex"))

    // console.log("Signature in bytes", this.signature.split('').map(c => c.charCodeAt(c) || 0))
    // console.log("Signature in bytes length", this.signature.split('').map(c => c.charCodeAt(c) || 0).length)




    const signatureObject = {
      messageHash: hash,
      // r: toHexString(this.signature.split('').map(c => c.charCodeAt(c) || 0).slice(0, 31)),
      r: '0x' + this.signature.split('').slice(0, 63).join(''),
      // s: toHexString(this.signature.split('').map(c => c.charCodeAt(c) || 0).slice(31, 63)),
      s: '0x' + this.signature.split('').slice(64, 127).join(''),
      v: "0x1c",
      // signature: this.signature
    }
  //   const signatureObject2 = {
  //     messageHash: "0x1da44b586eb0729ff70a73c326926f6ed5a25f5b056e7f47fbc6e58d86871655",
  //     v: "0x1c",
  //     r: "0xb91467e570a6466aa9e9876cbcd013baba02900b8979d43fe208a4a4f339f5fd",
  //     s: "0x6007e74cd82e037b800186422fc2da167c747ef045e5d18a5f5d4300f8e1a029"
  // }
  //   console.log("verifySignature -> signatureObject", signatureObject)
  //   const pubKeyRecovered = this.web3.eth.accounts.recover(signatureObject)
  //   console.log("verifySignature -> pubKeyRecovered", pubKeyRecovered)



  // const pubKeyRecovered = ethers.utils.recoverAddress(ethers.utils.arrayify(hash), signatureObject)



  //   return pubKeyRecovered
  }

  prefixedHash() {
    const params = this.web3.eth.abi.encodeParameters( ["uint256", "string"], [this.teamId || 0, this.name || ''] )
    const hash = this.web3.utils.soliditySha3(params);
    const prefixedHash = this.web3.utils.soliditySha3("\x19Ethereum Signed Message:\n32", hash)

    return prefixedHash
  }

}

module.exports = TeamNameValidation;
 