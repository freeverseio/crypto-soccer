const express = require('express');
const bodyParser = require('body-parser');
const routes = require('./routes');
const logger = require('./logger');
const { serverConfig } = require('./config.js');
const { port } = serverConfig;

const app = express();

app.use(bodyParser.json());
app.use('/', routes);

const start = async () => {
  try {
    app.listen(port, () => {
      logger.info(`Running on ${port}`);
    });
  } catch (e) {
    logger.error(e);
    process.exit(1);
  }
};

start();
