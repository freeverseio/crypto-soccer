const HorizonService = require('../services/HorizonService.js');
const { selectOwnerMaxBidAllowed } = require('../repositories/index.js');
const utc = require('dayjs/plugin/utc');
const dayjs = require('dayjs');
dayjs.extend(utc);
class Validation {
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
}

module.exports = Validation;
