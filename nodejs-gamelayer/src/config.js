const { PG_CONNECTION_STRING } = process.env;

const postgreSQLConfig = {
  connectionString: PG_CONNECTION_STRING
};

module.exports = postgreSQLConfig;
