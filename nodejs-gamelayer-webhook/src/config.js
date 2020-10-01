const { HORIZON_URL, LOG_LEVEL, PORT } = process.env;

const horizonConfig = {
  url: HORIZON_URL,
};

const serverConfig = {
  level: LOG_LEVEL || 'info',
  port: PORT,
};

module.exports = {
  horizonConfig,
  serverConfig,
};
