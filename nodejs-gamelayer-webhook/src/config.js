const { HORIZON_URL, LOG_LEVEL } = process.env;

const horizonConfig = {
  url: HORIZON_URL,
};

const loggerConfig = {
  level: LOG_LEVEL || 'debug',
};

module.exports = {
  horizonConfig,
  loggerConfig,
};
