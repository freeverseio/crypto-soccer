const { GAMELAYER_URL, HORIZON_URL } = process.env;

const horizonConfig = {
  url: HORIZON_URL,
};

const gamelayerConfig = {
  url: GAMELAYER_URL,
};

module.exports = {
  gamelayerConfig,
  horizonConfig,
};
