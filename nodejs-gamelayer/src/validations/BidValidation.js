const HorizonService = require('../services/HorizonService.js');
const { MINIMUM_DEFAULT_BID } = require('../config.js');
const { selectOwnerMaxBidAllowed } = require('../repositories/index.js');
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

  async isAllowedByUnpayments({ owner }) {
    const unpayments = await HorizonService.getUnpaymentsByOwner({ owner });
    if (unpayments.length) {
      switch (unpayments.length) {
        case 1:
        case 2:
          const timeOfUnpayment = dayjs(unpayments[0].timeOfUnpayment).utc();
          const daysSinceLastTimeOfUnpayment = dayjs.utc().diff(timeOfUnpayment, 'day');
          if (daysSinceLastTimeOfUnpayment < 7) {
            return false;
          }
          break;
        default:
          return false;
      }
    }
    return true;
  }

  async isAllowedToBidBySignerOwner() {
    const isSignerOwner = await this.isSignerOwner();
    if (!isSignerOwner) {
      return false;
    }
    return true;
  }

  async isAllowedToBidByMaxBidAllowedByOwner({ owner, totalPrice }) {
    const maxBidAllowedByOwnerRow = await selectOwnerMaxBidAllowed({ owner });
    if (maxBidAllowedByOwnerRow) {
      const maxBid = parseInt(maxBidAllowedByOwnerRow.max_bid_allowed);
      if (Number.isInteger(maxBid)) {
        if (totalPrice > maxBid) {
          return false;
        }
      }
    }

    return true;
  }

  async isAllowedToBidBySpentInWorldplayers({ owner }) {
    const teamsByOwner = await HorizonService.getTeamsByOwner({ owner });
    let hasSpentInWorldPlayers = false;
    for (const team of teamsByOwner) {
      const hasSpentWP = await HorizonService.hasSpentInWorldPlayers({ teamId: team.teamId });
      hasSpentInWorldPlayers = hasSpentInWorldPlayers || hasSpentWP;
    }

    if (hasSpentInWorldPlayers) {
      return true;
    }

    return false;
  }

  async isAllowedToBid() {
    const isAllowedToBidBySignerOwner = await this.isAllowedToBidBySignerOwner();
    if (!isAllowedToBidBySignerOwner) {
      return false;
    }

    const owner = await this.signerAddress();
    const auction = await HorizonService.getAuction({
      auctionId: this.auctionId,
    });
    const totalPrice = parseInt(auction.price) + parseInt(this.extraPrice);

    const isAllowedToBidByMaxBidAllowedByOwner = await this.isAllowedToBidByMaxBidAllowedByOwner({ owner, totalPrice });
    if (!isAllowedToBidByMaxBidAllowedByOwner) {
      return false;
    }

    const isAllowedToBidByUnpayments = await this.isAllowedByUnpayments({ owner });
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

    const isAllowedToBidByWorldplayers = await this.isAllowedToBidBySpentInWorldplayers({ owner });
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

module.exports = BidValidation;
