const HorizonService = require("../services/HorizonService.js");
class PlayerValidation {
    constructor({ playerId, name, signature, web3 }) {
      this.playerId = playerId
      this.name = name
      this.signature = signature
      this.web3 = web3
    }
    
    hash() {
      const params = this.web3.eth.abi.encodeParameters( ["uint256", "string"], [this.playerId || 0, this.name || ''] )
      return this.web3.utils.soliditySha3(params);
    }
    
    prefixedHash() {
      const params = this.web3.eth.abi.encodeParameters( ["uint256", "string"], [this.playerId || 0, this.name || ''] )
      const hash = this.web3.utils.soliditySha3(params);
      const prefixedHash = this.web3.utils.soliditySha3("\x19Ethereum Signed Message:\n32", hash)
    
      return prefixedHash
    }
  
    async signerAddress() {
      const hash = this.prefixedHash()
      const signatureObject = {
        messageHash: hash,
        r: '0x' + this.signature.split('').slice(0, 66).join(''),
        s: '0x' + this.signature.split('').slice(66, 130).join(''),
        v: '0x' + this.signature.split('').slice(130, 132).join(''),
      }
  
      const pubKeyRecovered = await this.web3.eth.accounts.recover(signatureObject)
  
      return pubKeyRecovered
    }

    async isSignerOwner() {
      const playerOwner = await HorizonService.getPlayerOwner({ playerId: this.playerId })
      const signerAddress = await this.signerAddress()

      return playerOwner === signerAddress
    }
  
  }
  
  module.exports = PlayerValidation;
   