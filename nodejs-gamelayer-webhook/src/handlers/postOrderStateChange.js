const logger = require('../logger');
const HorizonService = require('../services/HorizonService');

const postOrderStateChange = async (req, res) => {
  try {
    const [body] = req.body;
    const {
      escrow: {
        name,
        status,
        trustee_shortlink: { hash: hashTrusteeShortLink },
        shortlink: { hash },
      },
    } = body;
    const regex = /[a-f0-9]{64}/g;
    const matchedAuctionId = name.match(regex);
    const auctionId = matchedAuctionId && matchedAuctionId[0] ? matchedAuctionId[0] : '';
    logger.debug(
      `Received:\nAuctionId: ${auctionId}\n--------\nTransaction Name: ${name}\n--------\nStatus: ${status}\n--------\nTrustee Shortlink hash: ${hashTrusteeShortLink}\n--------\nShortlink Hash: ${hash}\n--------\n`
    );

    if (auctionId) {
      await HorizonService.processAuction({ auctionId });
    }

    res.sendStatus(200);
  } catch (e) {
    logger.error(e);
    res.sendStatus(500);
  }
};

module.exports = postOrderStateChange;
