const {
  GAMELAYER_URL,
  HORIZON_URL,
  MAILBOX_CRON,
  PG_CONNECTION_STRING,
  LOG_LEVEL,
} = process.env;

const serverConfig = {
  level: LOG_LEVEL || 'info',
};

const postgreSQLConfig = {
  connectionString: PG_CONNECTION_STRING,
};
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
  cancelled: 'cancelled',
};

const auctionStates = {
  started: 'started',
  accepted: 'accepted',
  rejected: 'rejected',
  withdrawableBySeller: 'withadrable_by_seller',
  paying: 'paying',
  assetFrozen: 'asset_frozen',
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
  postgreSQLConfig,
  MAILBOX_CRON,
  serverConfig,
};
