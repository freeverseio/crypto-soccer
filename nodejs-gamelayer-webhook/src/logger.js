const winston = require('winston');
const { serverConfig } = require('./config.js');

const logger = winston.createLogger({
  transports: [
    new winston.transports.Console({
      level: serverConfig.level,
      format: winston.format.combine(winston.format.colorize(), winston.format.simple()),
    }),
  ],
});

module.exports = logger;
