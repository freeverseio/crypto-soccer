const HorizonService = require('../services/HorizonService.js');
const {
  selectOwnerMaximumBid,
  updateOwnerMaximumBid,
} = require('../repositories');
const { MINIMUM_DEFAULT_BID } = require('../config.js');

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
      ['bytes32', 'bytes32', 'uint256', 'bool'],
      [
        '0x' + this.auctionId || '',
        bidHiddenPrice || '',
        this.teamId || 0,
        false,
      ]
    );
    return this.web3.utils.soliditySha3(params);
  }

  prefixedHash() {
    const prefixedHash = this.web3.utils.soliditySha3(
      '\x19Ethereum Signed Message:\n32',
      this.hash()
    );

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

    const pubKeyRecovered = await this.web3.eth.accounts.recover(
      signatureObject
    );

    return pubKeyRecovered;
  }

  async isSignerOwner() {
    const teamOwner = await HorizonService.getTeamOwner({
      teamId: this.teamId,
    });
    const signerAddress = await this.signerAddress();

    return teamOwner === signerAddress;
  }

  async isAllowedToBid() {
    const isSignerOwner = await this.isSignerOwner();
    if (!isSignerOwner) {
      return false;
    }

    const owner = await this.signerAddress();
    let maximumBid = await selectOwnerMaximumBid({ owner });

    if (parseInt(maximumBid) === 0) {
      return false;
    }
    maximumBid = parseInt(maximumBid) || MINIMUM_DEFAULT_BID;

    const auction = await HorizonService.getAuction({
      auctionId: this.auctionId,
    });
    const totalPrice = parseInt(auction.price) + parseInt(this.extraPrice);

    if (maximumBid >= totalPrice) {
      return true;
    }

    const bidsPayed = await HorizonService.getBidsPayedByOwner({ owner });
    const totalAmountSpent = bidsPayed.reduce(
      (acc, curr) =>
        (acc +=
          parseInt(curr.extraPrice) + parseInt(curr.auctionByAuctionId.price)),
      0
    );
    const newMaximumBid = parseInt(totalAmountSpent) * 1.5;
    await updateOwnerMaximumBid({ owner, maximumBid: newMaximumBid });

    if (newMaximumBid >= totalPrice) {
      return true;
    }

    return false;
  }
}

module.exports = BidValidation;
