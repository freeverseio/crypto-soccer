const EncodingMatchLog = artifacts.require('EncodingMatchLog');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(EncodingMatchLog);
    })
    .catch(console.error);
};

