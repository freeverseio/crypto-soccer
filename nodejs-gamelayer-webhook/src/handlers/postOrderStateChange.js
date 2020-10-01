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
    const regex = /^[a-f0-9]{64}$/gm;
    const auctionId = name.match(regex);
    logger.debug(
      `Received:\nAuctionId: ${auctionId}\n--------\nTransactionName${name}\n--------\nStatus: ${status}\n--------\nTrustee Shortlink hash: ${hashTrusteeShortLink}\n--------\nShortlink Hash: ${hash}\n--------\n`
    );
    // user horizon service to call notary new mutation
    await HorizonService.processAuction({ auctionId });
    res.sendStatus(200);
  } catch (e) {
    logger.error(e);
    res.sendStatus(500);
  }
};

module.exports = postOrderStateChange;
