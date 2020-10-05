const { GAMELAYER_URL, HORIZON_URL } = process.env;

const horizonConfig = {
  url: HORIZON_URL,
};

const gamelayerConfig = {
  url: GAMELAYER_URL,
};

const mailboxCronEntities = {
  auction: 'auction',
  offer: 'offer',
};

const offerStates = {
  started: 'started',
  accepted: 'accepted',
  rejected: 'rejected',
};

const auctionStates = {
  started: 'started',
  accepted: 'accepted',
  rejected: 'rejected',
  withdrawableBySeller: 'withdrawable_by_seller',
  paying: 'paying',
};

const bidStates = {
  paying: 'paying',
  accepted: 'accepted',
  failed: 'failed',
  paid: 'paid',
};

const mailboxTypes = {
  offer: 'offer',
  auction: 'auction',
  promo: 'promo',
  news: 'news',
  incident: 'incident',
  welcome: 'welcome',
};

const timeToNotifyPendinPayingBids = 12;

module.exports = {
  gamelayerConfig,
  horizonConfig,
  mailboxCronEntities,
  offerStates,
  auctionStates,
  mailboxTypes,
  bidStates,
  timeToNotifyPendinPayingBids,
};
