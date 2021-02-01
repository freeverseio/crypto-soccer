class CustomError extends Error {
  constructor(statusCode, ...params) {
    super(...params);

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, CustomError);
    }

    this.name = 'CustomError';
    this.statusCode = statusCode;
  }
}

module.exports = CustomError;
