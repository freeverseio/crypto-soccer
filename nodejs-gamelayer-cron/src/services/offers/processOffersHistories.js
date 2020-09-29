const processAcceptedOffers = require('./processAcceptedOffers');
const processRejectedOffers = require('./processRejectedOffers');
const processStartedOffers = require('./processStartedOffers');
const selectLastChecked = require('../../repositories/selectLastChecked');
const updateLastChecked = require('../../repositories/updateLastChecked');
const HorizonService = require('../HorizonService');

const processOffersHistories = async () => {
  const offerLastChecked = await selectLastChecked({ entity: 'offer' });

  const lastOffersHistories = await HorizonService.getLastOfferHistories({
    lastChecked: offerLastChecked,
  });

  const newLastChecked =
    lastOffersHistories[lastOffersHistories.length].insertedAt;

  await updateLastChecked({ entity: 'offer', lastChecked: newLastChecked });

  for (const offerHistory of lastOffersHistories) {
    switch (offerHistory.state) {
      case 'started':
        await processStartedOffers({ offerHistory });
        break;
      case 'accepted':
        await processAcceptedOffers({ offerHistory });
        break;
      case 'rejected':
        await processRejectedOffers({ offerHistory });
        break;
    }
  }
};

module.exports = processOffersHistories;
