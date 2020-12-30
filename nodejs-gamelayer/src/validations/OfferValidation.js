const HorizonService = require('../services/HorizonService.js');
const { MINIMUM_DEFAULT_BID } = require('../config.js');
const utc = require('dayjs/plugin/utc');
const dayjs = require('dayjs');
const Validation = require('./Validation.js');
dayjs.extend(utc);
class OfferValidation {
  constructor({ currencyId, rnd, price, signature, web3, playerId, validUntil, buyerTeamId }) {
    this.currencyId = currencyId;
    this.playerId = playerId;
    this.price = price;
    this.rnd = rnd;
    this.signature = signature;
    this.web3 = web3;
    this.validUntil = validUntil;
    this.buyerTeamId = buyerTeamId;
    this.validation = new Validation();
  }
  computeAuctionId() {
    const paramsSellerHiddenPrice = this.web3.eth.abi.encodeParameters(
      ['uint8', 'uint256', 'uint256'],
      [this.currencyId || 0, this.price || 0, this.rnd || 0]
    );
    const sellerHiddenPrice = this.web3.utils.soliditySha3(paramsSellerHiddenPrice);
    const params = this.web3.eth.abi.encodeParameters(
      ['bytes32', 'uint256', 'uint32'],
      [sellerHiddenPrice || '', this.playerId || 0, this.validUntil || 0]
    );
    return this.web3.utils.soliditySha3(params).slice(2);
  }

  hash() {
    const dummyExtraPrice = 0;
    const dummyRnd = 0;
    const auctionId = this.computeAuctionId();
    const paramsBidHiddenPrice = this.web3.eth.abi.encodeParameters(
      ['uint256', 'uint256'],
      [dummyExtraPrice, dummyRnd]
    );
    const bidHiddenPrice = this.web3.utils.soliditySha3(paramsBidHiddenPrice);
    const params = this.web3.eth.abi.encodeParameters(
      ['bytes32', 'bytes32', 'uint256'],
      ['0x' + auctionId || '', bidHiddenPrice || '', this.buyerTeamId || 0]
    );
    return this.web3.utils.soliditySha3(params);
  }

  prefixedHash() {
    const prefixedHash = this.web3.utils.soliditySha3('\x19Ethereum Signed Message:\n32', this.hash());

    return prefixedHash;
  }

  async signerAddress() {
    const hash = this.prefixedHash();
    const signatureObject = {
      messageHash: hash,
      r: '0x' + this.signature.split('').slice(0, 66).join(''),
      s: '0x' + this.signature.split('').slice(66, 130).join(''),
      v: '0x' + this.signature.split('').slice(130, 132).join(''),
    };

    const pubKeyRecovered = await this.web3.eth.accounts.recover(signatureObject);

    return pubKeyRecovered;
  }

  async isSignerOwner() {
    const teamOwner = await HorizonService.getTeamOwner({
      teamId: this.buyerTeamId,
    });
    const signerAddress = await this.signerAddress();

    return teamOwner === signerAddress;
  }

  async isAllowedToBidBySignerOwner() {
    const isSignerOwner = await this.isSignerOwner();
    if (!isSignerOwner) {
      return false;
    }
    return true;
  }

  async isAllowedToOffer() {
    const isAllowedToBidBySignerOwner = await this.isAllowedToBidBySignerOwner();
    if (!isAllowedToBidBySignerOwner) {
      return false;
    }

    const owner = await this.signerAddress();
    const totalPrice = parseInt(this.price);

    const isAllowedToBidByMaxBidAllowedByOwner = await this.validation.isAllowedToBidByMaxBidAllowedByOwner({
      owner,
      totalPrice,
    });
    if (!isAllowedToBidByMaxBidAllowedByOwner) {
      return false;
    }

    const isAllowedToBidByUnpayments = await this.validation.isAllowedByUnpayments({ owner });
    if (!isAllowedToBidByUnpayments) {
      return false;
    }

    const isTotalPriceLessThanMinimum = parseInt(totalPrice) < parseInt(MINIMUM_DEFAULT_BID);
    if (isTotalPriceLessThanMinimum) {
      return true;
    }

    const hasAuctionPass = await HorizonService.hasAuctionPass({ owner });
    if (hasAuctionPass) {
      return true;
    }

    const isAllowedToBidByWorldplayers = await this.validation.isAllowedToBidBySpentInWorldplayers({ owner });
    if (isAllowedToBidByWorldplayers) {
      return true;
    }

    const bidsPayed = await HorizonService.getBidsPayedByOwner({ owner });
    const hasSpentInBids = bidsPayed.length > 0;
    if (hasSpentInBids) {
      return true;
    }

    return false;
  }
}

module.exports = OfferValidation;
