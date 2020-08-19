const { PG_CONNECTION_STRING, HORIZON_URL } = process.env;

const postgreSQLConfig = {
  connectionString: PG_CONNECTION_STRING
};

const horizonConfig = {
  url: HORIZON_URL,
}

module.exports = {
  postgreSQLConfig,
  horizonConfig,
};
