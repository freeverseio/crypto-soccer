const HorizonService = require("../services/HorizonService.js");
const { selectTeamMaximumBid, updateTeamMaximumBid } = require("../repositories");

class BidValidation {
  constructor({ teamId, rnd, auctionId, extraPrice, signature, web3 }) {
    this.teamId = teamId
    this.auctionId = auctionId
    this.extraPrice = extraPrice
    this.rnd = rnd
    this.signature = signature
    this.web3 = web3  		
  }

  async isSignerAllowedToBid() {
    const minimumDefaultBid = 10
    let maximumBid = await selectTeamMaximumBid({ teamId: this.teamId })

    if(parseInt(maximumBid) === 0) {
      return false
    }
    maximumBid = maximumBid || minimumDefaultBid


    const auction = await HorizonService.getAuction({ auctionId: this.auctionId })
    const totalPrice = parseInt(auction.price) + parseInt(this.extraPrice)

    if(parseInt(maximumBid) > totalPrice) {
      return true
    }

    const bidsPayed = await HorizonService.getBidsPayed({ teamId: this.teamId })
    const totalAmountSpent = bidsPayed.reduce((acc, curr) => acc += parseInt(curr.extraPrice) + parseInt(curr.auctionByAuctionId.price), 0)
    const newMaximumBid = totalAmountSpent * 1.5
    await updateTeamMaximumBid({ teamId: this.teamId, teamMaximumBid: newMaximumBid })

    if(parseInt(newMaximumBid) > totalPrice) {
      return true
    }
    
    return false
  }

}

module.exports = BidValidation;
 