const { HORIZON_URL, LOG_LEVEL } = process.env;

const horizonConfig = {
  url: HORIZON_URL,
};

const serverConfig = {
  level: LOG_LEVEL || 'info',
  port: 5000,
};

module.exports = {
  horizonConfig,
  serverConfig,
};
