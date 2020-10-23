const winston = require('winston');
const { serverConfig } = require('./config.js');

const logFormatter = winston.format.printf((info) => {
  let { timestamp, level, code, stack, message } = info;
  code = code ? ` ${code}` : '';
  message = stack || message;

  return `${timestamp} ${level}${code}: ${message}`;
});

const logger = winston.createLogger({
  format: winston.format.errors({ stack: true }),
  transports: [
    new winston.transports.Console({
      level: serverConfig.level,
      format: winston.format.combine(
        winston.format.colorize(),
        winston.format.timestamp(),
        logFormatter
      ),
    }),
  ],
});

module.exports = logger;
