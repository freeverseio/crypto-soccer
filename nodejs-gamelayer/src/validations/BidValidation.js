const HorizonService = require('../services/HorizonService.js');
const Validation = require('./Validation');
const { MINIMUM_DEFAULT_BID, errorCodes } = require('../config.js');
const utc = require('dayjs/plugin/utc');
const dayjs = require('dayjs');
dayjs.extend(utc);
class BidValidation {
  constructor({ teamId, rnd, auctionId, extraPrice, signature, web3 }) {
    this.teamId = teamId;
    this.auctionId = auctionId;
    this.extraPrice = extraPrice;
    this.rnd = rnd;
    this.signature = signature;
    this.web3 = web3;
    this.validation = new Validation();
  }

  hash() {
    const paramsBidHiddenPrice = this.web3.eth.abi.encodeParameters(
      ['uint256', 'uint256'],
      [this.extraPrice || 0, this.rnd || 0]
    );
    const bidHiddenPrice = this.web3.utils.soliditySha3(paramsBidHiddenPrice);
    const params = this.web3.eth.abi.encodeParameters(
      ['bytes32', 'bytes32', 'uint256'],
      ['0x' + this.auctionId || '', bidHiddenPrice || '', this.teamId || 0]
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
      teamId: this.teamId,
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

  async isAllowedToBid() {
    const isAllowedToBidBySignerOwner = await this.isAllowedToBidBySignerOwner();
    if (!isAllowedToBidBySignerOwner) {
      return { allowed: false, code: errorCodes.BID_NOT_ALLOWED };
    }

    const owner = await this.signerAddress({ web3: this.web3, signature: this.signature });
    const auction = await HorizonService.getAuction({
      auctionId: this.auctionId,
    });
    const totalPrice = parseInt(auction.price) + parseInt(this.extraPrice);

    const isAllowedToBidByMaxBidAllowedByOwner = await this.validation.isAllowedToBidByMaxBidAllowedByOwner({
      owner,
      totalPrice,
    });
    if (!isAllowedToBidByMaxBidAllowedByOwner) {
      return { allowed: false, code: errorCodes.BID_NOT_ALLOWED_BY_BAN };
    }

    const isAllowedToBidByUnpayments = await this.validation.isAllowedByUnpayments({ owner });
    if (!isAllowedToBidByUnpayments) {
      return { allowed: false, code: errorCodes.BID_NOT_ALLOWED_BY_BAN };
    }

    const isTotalPriceLessThanMinimum = parseInt(totalPrice) < parseInt(MINIMUM_DEFAULT_BID);
    if (isTotalPriceLessThanMinimum) {
      return { allowed: true };
    }

    const hasAuctionPass = await HorizonService.hasAuctionPass({ owner });
    if (hasAuctionPass) {
      return { allowed: true };
    }

    const isAllowedToBidByWorldplayers = await this.validation.isAllowedToBidBySpentInWorldplayers({ owner });
    if (isAllowedToBidByWorldplayers) {
      return { allowed: true };
    }

    const bidsPayed = await HorizonService.getBidsPayedByOwner({ owner });
    const hasSpentInBids = bidsPayed.length > 0;
    if (hasSpentInBids) {
      return { allowed: true };
    }

    return { allowed: false, code: errorCodes.BID_NOT_ALLOWED };
  }
}

module.exports = BidValidation;
