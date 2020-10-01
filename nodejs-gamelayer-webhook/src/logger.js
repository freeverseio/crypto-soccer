const winston = require('winston');
const { loggerConfig } = require('./config.js');

const logger = winston.createLogger({
  transports: [
    new winston.transports.Console({
      level: loggerConfig.level,
      format: winston.format.combine(winston.format.colorize(), winston.format.simple()),
    }),
  ],
});

module.exports = logger;
