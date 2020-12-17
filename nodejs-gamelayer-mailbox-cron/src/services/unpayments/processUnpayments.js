const HorizonService = require('../HorizonService');
const logger = require('../../logger');

const processUnpayments = async () => {
  const lastUnpayments = await HorizonService.getLastUnpayments();

  const areNewUnpayments = lastUnpayments.length > 0;

  logger.info(
    areNewUnpayments
      ? `Processing New Unpayments`
      : `Processing Unpayments - No new unpayments`
  );
  if (areNewUnpayments) {
    for (const unpayment of lastUnpayments) {
      try {
        //get teams of owner
        const teamsOfOwner = await HorizonService.getTeamIdsFromOwner({
          owner: unpayment.owner,
        });
        const unpaymentsByOwner = await HorizonService.getUnpaymentsByOwner({
          owner,
        });
        //send mailbox for each team
        message = {
          destinatary: team,
          category: 'ban',
          auctionId: '',
          title: '',
          text: 'ban',
          customImageUrl: '',
          metadata: '',
        };
        if (unpaymentsByOwner.length > 2) {
          message = {
            destinatary: team,
            category: 'permaban',
            auctionId: '',
            title: '',
            text: 'permaban',
            customImageUrl: '',
            metadata: '',
          };
        }

        for (const team of teamsOfOwner) {
          await GamelayerService.setMessage(message);
        }
        //set unpayment notified
        await HorizonService.setUnpaymentNotified({ unpayment });
      } catch (e) {
        logger.info(`Error processing unpayment: ${JSON.stringify(unpayment)}`);
        logger.error(e);
      }
    }
  }
};

module.exports = processUnpayments;
