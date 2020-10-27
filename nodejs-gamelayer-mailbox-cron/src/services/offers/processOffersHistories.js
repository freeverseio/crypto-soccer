const dayjs = require('dayjs');
const processAcceptedOffers = require('./processAcceptedOffers');
const processRejectedOffers = require('./processRejectedOffers');
const processStartedOffers = require('./processStartedOffers');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxTypes, offerStates } = require('../../config');
const logger = require('../../logger');

const processOffersHistories = async () => {
  const offerLastChecked = await selectLastChecked({
    entity: mailboxTypes.offer,
  });

  const lastOffersHistories = await HorizonService.getLastOfferHistories({
    lastChecked: dayjs(offerLastChecked).add(1, 'second').format(),
  });

  const newLastChecked =
    lastOffersHistories[lastOffersHistories.length - 1] &&
    lastOffersHistories[lastOffersHistories.length - 1].insertedAt
      ? lastOffersHistories[lastOffersHistories.length - 1].insertedAt
      : offerLastChecked;

  await updateLastChecked({
    entity: mailboxTypes.offer,
    lastChecked: newLastChecked,
  });
  const areNewHistories =
    dayjs(offerLastChecked).format() != dayjs(newLastChecked).format();

  logger.info(
    areNewHistories
      ? `Processing Offer Histories - LastCheckedOfferTime: ${offerLastChecked} - NewLastCheckedOfferTime: ${newLastChecked}`
      : `Processing Offer Histories - No new offer histories since lastCheckedOfferTime: ${offerLastChecked}`
  );
  if (areNewHistories) {
    for (const offerHistory of lastOffersHistories) {
      try {
        switch (offerHistory.state.toLowerCase()) {
          case offerStates.started:
            await processStartedOffers({ offerHistory });
            break;
          case offerStates.accepted:
            await processAcceptedOffers({ offerHistory });
            break;
          case offerStates.cancelled:
            await processRejectedOffers({ offerHistory });
            break;
        }
      } catch (e) {
        logger.info(
          `Error processing offerHistory: ${JSON.stringify(offerHistory)}`
        );
        logger.error(e);
      }
    }
  }
};

module.exports = processOffersHistories;
