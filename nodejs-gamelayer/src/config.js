const { PG_CONNECTION_STRING, HORIZON_URL } = process.env;

const postgreSQLConfig = {
  connectionString: PG_CONNECTION_STRING,
};

const horizonConfig = {
  url: HORIZON_URL,
};

const errorCodes = {
  BID_NOT_ALLOWED: 100,
  OFFER_NOT_ALLOWED: 101,
  BID_NOT_ALLOWED_BY_BAN: 102,
  OFFER_NOT_ALLOWED_BY_BAN: 103,
};

module.exports = {
  postgreSQLConfig,
  horizonConfig,
  MINIMUM_DEFAULT_BID: 1000,
  errorCodes,
};
