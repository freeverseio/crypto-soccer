const { HORIZON_URL, LOG_LEVEL } = process.env;

const horizonConfig = {
  url: HORIZON_URL,
};

const loggerConfig = {
  level: LOG_LEVEL || 'info',
};

module.exports = {
  horizonConfig,
  loggerConfig,
};
