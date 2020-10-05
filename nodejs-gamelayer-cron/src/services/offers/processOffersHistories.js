const processAcceptedOffers = require('./processAcceptedOffers');
const processRejectedOffers = require('./processRejectedOffers');
const processStartedOffers = require('./processStartedOffers');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxCronTypes, offerStates } = require('../config');

const processOffersHistories = async () => {
  const offerLastChecked = await selectLastChecked({
    entity: mailboxCronTypes.offer,
  });

  const lastOffersHistories = await HorizonService.getLastOfferHistories({
    lastChecked: offerLastChecked,
  });

  const newLastChecked =
    lastOffersHistories[lastOffersHistories.length].insertedAt;

  await updateLastChecked({
    entity: mailboxCronTypes.offer,
    lastChecked: newLastChecked,
  });

  for (const offerHistory of lastOffersHistories) {
    switch (offerHistory.state) {
      case offerStates.started:
        await processStartedOffers({ offerHistory });
        break;
      case offerStates.accepted:
        await processAcceptedOffers({ offerHistory });
        break;
      case offerStates.rejected:
        await processRejectedOffers({ offerHistory });
        break;
    }
  }
};

module.exports = processOffersHistories;
