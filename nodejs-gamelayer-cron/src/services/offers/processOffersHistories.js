const processAcceptedOffers = require('./processAcceptedOffers');
const processRejectedOffers = require('./processRejectedOffers');
const processStartedOffers = require('./processStartedOffers');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');
const { mailboxTypes, offerStates } = require('../../config');

const processOffersHistories = async () => {
  const offerLastChecked = await selectLastChecked({
    entity: mailboxTypes.offer,
  });

  const lastOffersHistories = await HorizonService.getLastOfferHistories({
    lastChecked: offerLastChecked,
  });

  const newLastChecked =
    lastOffersHistories[lastOffersHistories.length - 1].insertedAt;

  await updateLastChecked({
    entity: mailboxTypes.offer,
    lastChecked: newLastChecked,
  });

  for (const offerHistory of lastOffersHistories) {
    switch (offerHistory.state.toLowerCase()) {
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
